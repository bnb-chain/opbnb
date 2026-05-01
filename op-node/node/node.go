package node

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/ethereum-optimism/optimism/op-node/node/safedb"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	plasma "github.com/ethereum-optimism/optimism/op-plasma"
	"github.com/ethereum-optimism/optimism/op-service/httputil"

	"github.com/hashicorp/go-multierror"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-node/heartbeat"
	"github.com/ethereum-optimism/optimism/op-node/metrics"
	"github.com/ethereum-optimism/optimism/op-node/p2p"
	"github.com/ethereum-optimism/optimism/op-node/rollup/conductor"
	"github.com/ethereum-optimism/optimism/op-node/rollup/driver"
	"github.com/ethereum-optimism/optimism/op-node/rollup/sync"
	"github.com/ethereum-optimism/optimism/op-node/version"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/oppprof"
	"github.com/ethereum-optimism/optimism/op-service/retry"
	"github.com/ethereum-optimism/optimism/op-service/sources"
)

var ErrAlreadyClosed = errors.New("node is already closed")

type closableSafeDB interface {
	derive.SafeHeadListener
	SafeDBReader
	io.Closer
}

type OpNode struct {
	log        log.Logger
	appVersion string
	metrics    *metrics.Metrics

	l1HeadsSub     ethereum.Subscription // Subscription to get L1 heads (automatically re-subscribes on error)
	l1SafeSub      ethereum.Subscription // Subscription to get L1 safe blocks, a.k.a. justified data (polling)
	l1FinalizedSub ethereum.Subscription // Subscription to get L1 safe blocks, a.k.a. justified data (polling)

	l1Source  *sources.L1Client     // L1 Client to fetch data from
	l2Driver  *driver.Driver        // L2 Engine to Sync
	l2Source  *sources.EngineClient // L2 Execution Engine RPC bindings
	server    *rpcServer            // RPC server hosting the rollup-node API
	p2pNode   *p2p.NodeP2P          // P2P node functionality
	p2pSigner p2p.Signer            // p2p gossip application messages will be signed with this signer
	tracer    Tracer                // tracer to get events for testing/debugging
	runCfg    *RuntimeConfig        // runtime configurables

	safeDB closableSafeDB

	rollupHalt string // when to halt the rollup, disabled if empty

	pprofService *oppprof.Service
	metricsSrv   *httputil.HTTPServer

	l1Blob *sources.BSCBlobClient // L1 Blob Client to fetch blobs

	// some resources cannot be stopped directly, like the p2p gossipsub router (not our design),
	// and depend on this ctx to be closed.
	resourcesCtx   context.Context
	resourcesClose context.CancelFunc

	// Indicates when it's safe to close data sources used by the runtimeConfig bg loader
	runtimeConfigReloaderDone chan struct{}

	closed atomic.Bool

	// catchUpEnabled mirrors cfg.CatchUp; controls whether Start() runs the pre-gossip catch-up phase.
	catchUpEnabled bool

	// gossipReady gates incoming gossip payloads during the startup catch-up phase.
	// While false, gossip payloads received via OnUnsafeL2Payload are silently dropped
	// to prevent the clSync queue from accumulating orphan payloads (parent != op-geth.UnsafeL2Head)
	// while op-geth's unsafe head is still being advanced via L1 derivation.
	// Set to true once catch-up completes (or is disabled / times out).
	gossipReady atomic.Bool

	// firstPayloadAllowed lets exactly one gossip payload pass through OnUnsafeL2Payload
	// while gossipReady is still false. This is required when running in ELSync mode:
	// the engineController initial state is syncStatusWillStartEL, which causes
	// IsEngineSyncing() to return true and prevents the driver eventLoop from running
	// derivation (the stepReqCh handler short-circuits with `continue`). Until at least
	// one payload reaches Driver.OnUnsafeL2Payload -> InsertUnsafePayload, the engine's
	// "Skipping EL sync" finalized-block check never fires and syncStatus is stuck.
	// Allowing exactly one payload through unblocks that transition and lets derivation
	// drive op-geth forward during the catch-up phase. After this single payload,
	// subsequent payloads continue to be dropped until gossipReady is flipped to true.
	firstPayloadAllowed atomic.Bool

	// cancels execution prematurely, e.g. to halt. This may be nil.
	cancel context.CancelCauseFunc
	halted atomic.Bool
}

// Startup catch-up parameters. Hardcoded; tweak here if needed.
const (
	// catchUpLagThreshold is how close op-geth's unsafe head timestamp must be
	// to the current wall-clock time before gossip is enabled.
	// On opBNB (~500ms blocks) 30s ≈ 60 blocks of remaining gap, well below the
	// threshold that triggers the activity loop in tested scenarios.
	catchUpLagThreshold = 30 * time.Second

	// catchUpMaxWait is the absolute maximum time we are willing to defer gossip.
	// If catch-up does not complete within this window (e.g. L1 derivation is unhealthy),
	// gossip is enabled regardless and the system degrades to the pre-fix behavior
	// rather than blocking forever.
	catchUpMaxWait = 10 * time.Minute

	// catchUpPollInterval is how often we re-check op-geth's unsafe head during catch-up.
	catchUpPollInterval = 5 * time.Second
)

// The OpNode handles incoming gossip
var _ p2p.GossipIn = (*OpNode)(nil)

// New creates a new OpNode instance.
// The provided ctx argument is for the span of initialization only;
// the node will immediately Stop(ctx) before finishing initialization if the context is canceled during initialization.
func New(ctx context.Context, cfg *Config, log log.Logger, snapshotLog log.Logger, appVersion string, m *metrics.Metrics) (*OpNode, error) {
	if err := cfg.Check(); err != nil {
		return nil, err
	}

	n := &OpNode{
		log:            log,
		appVersion:     appVersion,
		metrics:        m,
		rollupHalt:     cfg.RollupHalt,
		cancel:         cfg.Cancel,
		catchUpEnabled: cfg.CatchUp,
	}
	// If catch-up is disabled, gossip should be processed immediately as before.
	// Set the gate to "ready" up front so OnUnsafeL2Payload behaves identically to the pre-fix code path.
	if !n.catchUpEnabled {
		n.gossipReady.Store(true)
	}
	// not a context leak, gossipsub is closed with a context.
	n.resourcesCtx, n.resourcesClose = context.WithCancel(context.Background())

	err := n.init(ctx, cfg, snapshotLog)
	if err != nil {
		log.Error("Error initializing the rollup node", "err", err)
		// ensure we always close the node resources if we fail to initialize the node.
		if closeErr := n.Stop(ctx); closeErr != nil {
			return nil, multierror.Append(err, closeErr)
		}
		return nil, err
	}
	return n, nil
}

// waitForOpGethCatchUp blocks until op-geth's unsafe head timestamp is within
// catchUpLagThreshold of the current time, or until catchUpMaxWait elapses.
//
// Background:
// On RPC node restart, op-geth's unsafe head is frozen at the pre-restart height for
// the duration of the pod outage. When op-node comes back up and immediately subscribes
// to gossip, incoming gossip payloads have a parent that does not match op-geth's
// stale unsafe head; the clSync queue accumulates orphan payloads. The driver's
// checkForGapInUnsafeQueue then triggers alt-sync via an unbuffered rangeRequests
// channel, while alt-sync's mainLoop -- when promoting results back via receivePayload
// -- itself blocks on the driver's unsafeL2Payloads channel (buf=10). The two goroutines
// form a livelock that only releases through ctx timeouts, leaving the unsafe head
// stalled for some time.
//
// This function defers gossip subscription (via the gossipReady gate) until the L1
// derivation pipeline has advanced op-geth's unsafe head close enough to the live tip
// that no significant gap exists when gossip is finally enabled, eliminating the
// activity loop's preconditions at the source.
//
// Returns nil on successful catch-up; returns an error on context cancellation or timeout.
// In case of timeout, the caller should still enable gossip and degrade gracefully.
func (n *OpNode) waitForOpGethCatchUp(ctx context.Context) error {
	n.log.Info("starting op-geth catch-up phase before enabling gossip",
		"lag_threshold", catchUpLagThreshold,
		"max_wait", catchUpMaxWait,
	)

	deadline := time.Now().Add(catchUpMaxWait)
	ticker := time.NewTicker(catchUpPollInterval)
	defer ticker.Stop()

	for {
		// Query op-geth's current unsafe head.
		queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		unsafeHead, err := n.l2Source.L2BlockRefByLabel(queryCtx, eth.Unsafe)
		cancel()

		if err != nil {
			n.log.Warn("Failed to query op-geth unsafe head during catch-up, will retry", "error", err)
		} else {
			headTime := time.Unix(int64(unsafeHead.Time), 0)
			lag := time.Since(headTime)

			// Treat negative lag (clock skew or future-timestamp head) as caught up.
			if lag < catchUpLagThreshold {
				n.log.Info("op-geth caught up; enabling gossip", "unsafe_head", unsafeHead.Number, "lag", lag)
				return nil
			}

			n.log.Info("op-geth still catching up via L1 derivation", "unsafe_head", unsafeHead.Number,
				"lag", lag, "deadline_in", time.Until(deadline))
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("startup catch-up timeout after %v", catchUpMaxWait)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			continue
		}
	}
}

func (n *OpNode) init(ctx context.Context, cfg *Config, snapshotLog log.Logger) error {
	n.log.Info("Initializing rollup node", "version", n.appVersion)
	if err := n.initTracer(ctx, cfg); err != nil {
		return fmt.Errorf("failed to init the trace: %w", err)
	}
	if err := n.initL1(ctx, cfg); err != nil {
		return fmt.Errorf("failed to init L1: %w", err)
	}
	if err := n.initL1Blob(ctx, cfg); err != nil {
		return fmt.Errorf("failed to init L1 blob: %w", err)
	}
	if err := n.initL2(ctx, cfg, snapshotLog); err != nil {
		return fmt.Errorf("failed to init L2: %w", err)
	}
	if err := n.initRuntimeConfig(ctx, cfg); err != nil { // depends on L2, to signal initial runtime values to
		return fmt.Errorf("failed to init the runtime config: %w", err)
	}
	if err := n.initP2PSigner(ctx, cfg); err != nil {
		return fmt.Errorf("failed to init the P2P signer: %w", err)
	}
	if err := n.initP2P(ctx, cfg); err != nil {
		return fmt.Errorf("failed to init the P2P stack: %w", err)
	}
	// Only expose the server at the end, ensuring all RPC backend components are initialized.
	if err := n.initRPCServer(cfg); err != nil {
		return fmt.Errorf("failed to init the RPC server: %w", err)
	}
	if err := n.initMetricsServer(cfg); err != nil {
		return fmt.Errorf("failed to init the metrics server: %w", err)
	}
	n.metrics.RecordInfo(n.appVersion)
	n.metrics.RecordUp()
	n.initHeartbeat(cfg)
	if err := n.initPProf(cfg); err != nil {
		return fmt.Errorf("failed to init profiling: %w", err)
	}
	return nil
}

func (n *OpNode) initTracer(ctx context.Context, cfg *Config) error {
	if cfg.Tracer != nil {
		n.tracer = cfg.Tracer
	} else {
		n.tracer = new(noOpTracer)
	}
	return nil
}

func (n *OpNode) initL1(ctx context.Context, cfg *Config) error {
	l1Node, rpcCfg, err := cfg.L1.Setup(ctx, n.log, &cfg.Rollup)
	if err != nil {
		return fmt.Errorf("failed to get L1 RPC client: %w", err)
	}

	// Set the RethDB path in the EthClientConfig, if there is one configured.
	rpcCfg.EthClientConfig.RethDBPath = cfg.RethDBPath

	n.l1Source, err = sources.NewL1Client(
		client.NewInstrumentedRPC(l1Node, &n.metrics.RPCMetrics.RPCClientMetrics), n.log, n.metrics.L1SourceCache, rpcCfg)
	if err != nil {
		return fmt.Errorf("failed to create L1 source: %w", err)
	}

	if err := cfg.Rollup.ValidateL1Config(ctx, n.l1Source); err != nil {
		return fmt.Errorf("failed to validate the L1 config: %w", err)
	}

	// Keep subscribed to the L1 heads, which keeps the L1 maintainer pointing to the best headers to sync
	n.l1HeadsSub = event.ResubscribeErr(time.Second*10, func(ctx context.Context, err error) (event.Subscription, error) {
		if err != nil {
			n.log.Warn("resubscribing after failed L1 subscription", "err", err)
		}
		return eth.WatchHeadChanges(ctx, n.l1Source, n.OnNewL1Head)
	})

	if fallbackClient, ok := l1Node.(*sources.FallbackClient); ok {
		fallbackClient.RegisterSubscribeFunc(func() (event.Subscription, error) {
			return eth.WatchHeadChanges(n.resourcesCtx, n.l1Source, n.OnNewL1Head)
		}, &n.l1HeadsSub)
		fallbackClient.RegisterMetrics(n.metrics)
	}
	go func() {
		err, ok := <-n.l1HeadsSub.Err()
		if !ok {
			return
		}
		n.log.Error("l1 heads subscription error", "err", err)
	}()

	// Poll for the safe L1 block and finalized block,
	// which only change once per epoch at most and may be delayed.
	n.l1SafeSub = eth.PollBlockChanges(n.log, n.l1Source, n.OnNewL1Safe, eth.Safe,
		cfg.L1EpochPollInterval, time.Second*10)
	n.l1FinalizedSub = eth.PollBlockChanges(n.log, n.l1Source, n.OnNewL1Finalized, eth.Finalized,
		cfg.L1EpochPollInterval, time.Second*10)
	return nil
}

func (n *OpNode) initRuntimeConfig(ctx context.Context, cfg *Config) error {
	// attempt to load runtime config, repeat N times
	n.runCfg = NewRuntimeConfig(n.log, n.l1Source, &cfg.Rollup)

	confDepth := cfg.Driver.VerifierConfDepth
	reload := func(ctx context.Context) (eth.L1BlockRef, error) {
		fetchCtx, fetchCancel := context.WithTimeout(ctx, time.Second*10)
		l1Head, err := n.l1Source.L1BlockRefByLabel(fetchCtx, eth.Unsafe)
		fetchCancel()
		if err != nil {
			n.log.Error("failed to fetch L1 head for runtime config initialization", "err", err)
			return eth.L1BlockRef{}, err
		}

		// Apply confirmation-distance
		blNum := l1Head.Number
		if blNum >= confDepth {
			blNum -= confDepth
		}
		fetchCtx, fetchCancel = context.WithTimeout(ctx, time.Second*10)
		confirmed, err := n.l1Source.L1BlockRefByNumber(fetchCtx, blNum)
		fetchCancel()
		if err != nil {
			n.log.Error("failed to fetch confirmed L1 block for runtime config loading", "err", err, "number", blNum)
			return eth.L1BlockRef{}, err
		}

		fetchCtx, fetchCancel = context.WithTimeout(ctx, time.Second*10)
		err = n.runCfg.Load(fetchCtx, confirmed)
		fetchCancel()
		if err != nil {
			n.log.Error("failed to fetch runtime config data", "err", err)
			return l1Head, err
		}

		err = n.handleProtocolVersionsUpdate(ctx)
		return l1Head, err
	}

	// initialize the runtime config before unblocking
	if _, err := retry.Do(ctx, 5, retry.Fixed(time.Second*10), func() (eth.L1BlockRef, error) {
		ref, err := reload(ctx)
		if errors.Is(err, errNodeHalt) { // don't retry on halt error
			err = nil
		}
		return ref, err
	}); err != nil {
		return fmt.Errorf("failed to load runtime configuration repeatedly, last error: %w", err)
	}

	// start a background loop, to keep reloading it at the configured reload interval
	reloader := func(ctx context.Context, reloadInterval time.Duration) {
		if reloadInterval <= 0 {
			n.log.Debug("not running runtime-config reloading background loop")
			return
		}
		ticker := time.NewTicker(reloadInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// If the reload fails, we will try again the next interval.
				// Missing a runtime-config update is not critical, and we do not want to overwhelm the L1 RPC.
				l1Head, err := reload(ctx)
				if err != nil {
					if errors.Is(err, errNodeHalt) {
						n.halted.Store(true)
						if n.cancel != nil { // node cancellation is always available when started as CLI app
							n.cancel(errNodeHalt)
							return
						} else {
							n.log.Debug("opted to halt, but cannot halt node", "l1_head", l1Head)
						}
					} else {
						n.log.Warn("failed to reload runtime config", "err", err)
					}
				} else {
					n.log.Debug("reloaded runtime config", "l1_head", l1Head)
				}
			case <-ctx.Done():
				return
			}
		}
	}

	n.runtimeConfigReloaderDone = make(chan struct{})
	// Manages the lifetime of reloader. In order to safely Close the OpNode
	go func(ctx context.Context, reloadInterval time.Duration) {
		reloader(ctx, reloadInterval)
		close(n.runtimeConfigReloaderDone)
	}(n.resourcesCtx, cfg.RuntimeConfigReloadInterval) // this keeps running after initialization
	return nil
}

func (n *OpNode) initL1Blob(ctx context.Context, cfg *Config) error {
	// If Ecotone upgrade is not scheduled yet, then there is no need for a Blob API.
	if cfg.Rollup.EcotoneTime == nil {
		return nil
	}
	// Once the Ecotone upgrade is scheduled, we must have initialized the Blob API settings.
	if cfg.L1Blob == nil {
		return fmt.Errorf("missing L1 Blob Endpoint configuration: this API is mandatory for Ecotone upgrade at t=%d", *cfg.Rollup.EcotoneTime)
	}

	rpcClients, err := cfg.L1Blob.Setup(ctx, n.log)
	if err != nil {
		return fmt.Errorf("failed to setup L1 blob client: %w", err)
	}
	instrumentedClients := make([]client.RPC, 0)
	for _, rpc := range rpcClients {
		instrumentedClients = append(instrumentedClients, client.NewInstrumentedRPC(rpc, &n.metrics.RPCClientMetrics))
	}
	n.l1Blob = sources.NewBSCBlobClient(instrumentedClients)
	return nil
}

func (n *OpNode) initL2(ctx context.Context, cfg *Config, snapshotLog log.Logger) error {
	rpcClient, rpcCfg, err := cfg.L2.Setup(ctx, n.log, &cfg.Rollup)
	if err != nil {
		return fmt.Errorf("failed to setup L2 execution-engine RPC client: %w", err)
	}

	n.l2Source, err = sources.NewEngineClient(
		client.NewInstrumentedRPC(rpcClient, &n.metrics.RPCClientMetrics), n.log, n.metrics.L2SourceCache, rpcCfg,
	)
	if err != nil {
		return fmt.Errorf("failed to create Engine client: %w", err)
	}

	if err := cfg.Rollup.ValidateL2Config(ctx, n.l2Source, cfg.Sync.SyncMode == sync.ELSync); err != nil {
		return err
	}

	var sequencerConductor conductor.SequencerConductor = &conductor.NoOpConductor{}
	if cfg.ConductorEnabled {
		sequencerConductor = NewConductorClient(cfg, n.log, n.metrics)
	}

	// if plasma is not explicitly activated in the node CLI, the config + any error will be ignored.
	rpCfg, err := cfg.Rollup.GetOPPlasmaConfig()
	if cfg.Plasma.Enabled && err != nil {
		return fmt.Errorf("failed to get plasma config: %w", err)
	}
	plasmaDA := plasma.NewPlasmaDA(n.log, cfg.Plasma, rpCfg, n.metrics.PlasmaMetrics)
	if cfg.SafeDBPath != "" {
		n.log.Info("Safe head database enabled", "path", cfg.SafeDBPath)
		safeDB, err := safedb.NewSafeDB(n.log, cfg.SafeDBPath)
		if err != nil {
			return fmt.Errorf("failed to create safe head database at %v: %w", cfg.SafeDBPath, err)
		}
		n.safeDB = safeDB
	} else {
		n.safeDB = safedb.Disabled
	}
	n.l2Driver = driver.NewDriver(&cfg.Driver, &cfg.Rollup, n.l2Source, n.l1Source, n.l1Blob, n, n, n.log, snapshotLog, n.metrics, cfg.ConfigPersistence, n.safeDB, &cfg.Sync, sequencerConductor, plasmaDA)
	return nil
}

func (n *OpNode) initRPCServer(cfg *Config) error {
	server, err := newRPCServer(&cfg.RPC, &cfg.Rollup, n.l2Source.L2Client, n.l2Driver, n.safeDB, n.log, n.appVersion, n.metrics)
	if err != nil {
		return err
	}
	if n.p2pNode != nil {
		server.EnableP2P(p2p.NewP2PAPIBackend(n.p2pNode, n.log, n.metrics))
	}
	if cfg.RPC.EnableAdmin {
		server.EnableAdminAPI(NewAdminAPI(n.l2Driver, n.metrics, n.log))
		n.log.Info("Admin RPC enabled")
	}
	n.log.Info("Starting JSON-RPC server")
	if err := server.Start(); err != nil {
		return fmt.Errorf("unable to start RPC server: %w", err)
	}
	n.server = server
	return nil
}

func (n *OpNode) initMetricsServer(cfg *Config) error {
	if !cfg.Metrics.Enabled {
		n.log.Info("metrics disabled")
		return nil
	}
	n.log.Debug("starting metrics server", "addr", cfg.Metrics.ListenAddr, "port", cfg.Metrics.ListenPort)
	metricsSrv, err := n.metrics.StartServer(cfg.Metrics.ListenAddr, cfg.Metrics.ListenPort)
	if err != nil {
		return fmt.Errorf("failed to start metrics server: %w", err)
	}
	n.log.Info("started metrics server", "addr", metricsSrv.Addr())
	n.metricsSrv = metricsSrv
	return nil
}

func (n *OpNode) initHeartbeat(cfg *Config) {
	if !cfg.Heartbeat.Enabled {
		return
	}
	var peerID string
	if cfg.P2P.Disabled() {
		peerID = "disabled"
	} else {
		peerID = n.P2P().Host().ID().String()
	}

	payload := &heartbeat.Payload{
		Version: version.Version,
		Meta:    version.Meta,
		Moniker: cfg.Heartbeat.Moniker,
		PeerID:  peerID,
		ChainID: cfg.Rollup.L2ChainID.Uint64(),
	}

	go func(url string) {
		if err := heartbeat.Beat(n.resourcesCtx, n.log, url, payload); err != nil {
			log.Error("heartbeat goroutine crashed", "err", err)
		}
	}(cfg.Heartbeat.URL)
}

func (n *OpNode) initPProf(cfg *Config) error {
	n.pprofService = oppprof.New(
		cfg.Pprof.ListenEnabled,
		cfg.Pprof.ListenAddr,
		cfg.Pprof.ListenPort,
		cfg.Pprof.ProfileType,
		cfg.Pprof.ProfileDir,
		cfg.Pprof.ProfileFilename,
	)

	if err := n.pprofService.Start(); err != nil {
		return fmt.Errorf("failed to start pprof service: %w", err)
	}

	return nil
}

func (n *OpNode) initP2P(ctx context.Context, cfg *Config) error {
	if cfg.P2P != nil {
		// TODO(protocol-quest/97): Use EL Sync instead of CL Alt sync for fetching missing blocks in the payload queue.
		p2pNode, err := p2p.NewNodeP2P(n.resourcesCtx, &cfg.Rollup, n.log, cfg.P2P, n, n.l2Source, n.runCfg, n.metrics, false)
		if err != nil || p2pNode == nil {
			return err
		}
		n.p2pNode = p2pNode
		if n.p2pNode.Dv5Udp() != nil {
			go n.p2pNode.DiscoveryProcess(n.resourcesCtx, n.log, &cfg.Rollup, cfg.P2P.TargetPeers())
		}
	}
	return nil
}

func (n *OpNode) initP2PSigner(ctx context.Context, cfg *Config) error {
	// the p2p signer setup is optional
	if cfg.P2PSigner == nil {
		return nil
	}
	// p2pSigner may still be nil, the signer setup may not create any signer, the signer is optional
	var err error
	n.p2pSigner, err = cfg.P2PSigner.SetupSigner(ctx)
	return err
}

func (n *OpNode) Start(ctx context.Context) error {
	n.log.Info("Starting execution engine driver")
	// start driving engine: sync blocks by deriving them from L1 and driving them into the engine
	if err := n.l2Driver.Start(); err != nil {
		n.log.Error("Could not start a rollup node", "err", err)
		return err
	}

	// Optionally defer enabling gossip until op-geth's unsafe head has caught up to the live tip
	// via the L1 derivation pipeline. This avoids the driver/alt-sync livelock that occurs when
	// gossip floods in with payloads whose parent does not match op-geth's stale unsafe head.
	// Disabled by default; enable via --catch-up for RPC / verifier nodes.
	// See waitForOpGethCatchUp for full background.
	if n.catchUpEnabled {
		if err := n.waitForOpGethCatchUp(ctx); err != nil {
			// Catch-up failed (e.g. timeout, L1 derivation unhealthy). Enable gossip anyway
			// to avoid blocking the node forever; the system degrades to the pre-fix behavior.
			n.log.Warn("startup catch-up did not complete cleanly; enabling gossip anyway", "err", err)
		}
		n.gossipReady.Store(true)
		n.log.Info("gossip enabled; op-node fully active")
	}
	// If catch-up is disabled, gossipReady was already set to true in New(),
	// so OnUnsafeL2Payload behaves identically to the pre-fix code path.

	log.Info("Rollup node started")
	return nil
}

func (n *OpNode) OnNewL1Head(ctx context.Context, sig eth.L1BlockRef) {
	n.tracer.OnNewL1Head(ctx, sig)

	if n.l2Driver == nil {
		return
	}
	// Pass on the event to the L2 Engine
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := n.l2Driver.OnL1Head(ctx, sig); err != nil {
		n.log.Warn("failed to notify engine driver of L1 head change", "err", err)
	}
}

func (n *OpNode) OnNewL1Safe(ctx context.Context, sig eth.L1BlockRef) {
	if n.l2Driver == nil {
		return
	}
	// Pass on the event to the L2 Engine
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := n.l2Driver.OnL1Safe(ctx, sig); err != nil {
		n.log.Warn("failed to notify engine driver of L1 safe block change", "err", err)
	}
}

func (n *OpNode) OnNewL1Finalized(ctx context.Context, sig eth.L1BlockRef) {
	if n.l2Driver == nil {
		return
	}
	// Pass on the event to the L2 Engine
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := n.l2Driver.OnL1Finalized(ctx, sig); err != nil {
		n.log.Warn("failed to notify engine driver of L1 finalized block change", "err", err)
	}
}

func (n *OpNode) PublishL2Payload(ctx context.Context, envelope *eth.ExecutionPayloadEnvelope) error {
	n.tracer.OnPublishL2Payload(ctx, envelope)

	// publish to p2p, if we are running p2p at all
	if n.p2pNode != nil {
		payload := envelope.ExecutionPayload
		if n.p2pSigner == nil {
			return fmt.Errorf("node has no p2p signer, payload %s cannot be published", payload.ID())
		}
		n.log.Info("Publishing signed execution payload on p2p", "id", payload.ID())
		return n.p2pNode.GossipOut().PublishL2Payload(ctx, envelope, n.p2pSigner)
	}
	// if p2p is not enabled then we just don't publish the payload
	return nil
}

func (n *OpNode) OnUnsafeL2Payload(ctx context.Context, from peer.ID, envelope *eth.ExecutionPayloadEnvelope) error {
	// Drop gossip payloads received during the startup catch-up phase.
	// While op-geth's unsafe head is still catching up via L1 derivation, accepting
	// real-time gossip payloads would fill the clSync queue with orphan payloads
	// (parent != op-geth.UnsafeL2Head) and trigger the driver/alt-sync livelock.
	// Gossip is re-enabled once waitForOpGethCatchUp completes.
	// Any payloads dropped here are recovered by gossipsub mesh re-broadcasts and
	// alt-sync backfill once gossipReady is set.
	//
	// Exception: the very first payload is always let through, regardless of catch-up
	// state. In ELSync mode this is required to trigger the WillStartEL → FinishedEL
	// transition inside InsertUnsafePayload (the "Skipping EL sync ..." finalized-block
	// check). Without this, IsEngineSyncing() stays true and derivation is blocked
	// from running during catch-up, defeating the purpose of the wait. In CLSync mode
	// this single payload simply enters the clSync queue and is harmless.
	if !n.gossipReady.Load() {
		if !n.firstPayloadAllowed.Swap(true) {
			n.log.Info("allowing first gossip payload through during catch-up to unblock engine sync state",
				"id", envelope.ExecutionPayload.ID(), "peer", from)
			// fall through to the regular processing path below
		} else {
			n.log.Debug("dropping gossip payload during startup catch-up phase",
				"id", envelope.ExecutionPayload.ID(), "peer", from)
			return nil
		}
	}

	// ignore if it's from ourselves
	if n.p2pNode != nil && from == n.p2pNode.Host().ID() {
		return nil
	}

	n.tracer.OnUnsafeL2Payload(ctx, from, envelope)

	n.log.Info("Received signed execution payload from p2p", "id", envelope.ExecutionPayload.ID(), "peer", from)

	// Pass on the event to the L2 Engine
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if err := n.l2Driver.OnUnsafeL2Payload(ctx, envelope); err != nil {
		n.log.Warn("failed to notify engine driver of new L2 payload", "err", err, "id", envelope.ExecutionPayload.ID())
	}

	return nil
}

func (n *OpNode) RequestL2Range(ctx context.Context, start, end eth.L2BlockRef) error {
	if n.p2pNode != nil && n.p2pNode.AltSyncEnabled() {
		if unixTimeStale(start.Time, 12*time.Hour) {
			n.log.Debug("ignoring request to sync L2 range, timestamp is too old for p2p", "start", start, "end", end, "start_time", start.Time)
			return nil
		}
		return n.p2pNode.RequestL2Range(ctx, start, end)
	}
	n.log.Debug("ignoring request to sync L2 range, no sync method available", "start", start, "end", end)
	return nil
}

// unixTimeStale returns true if the unix timestamp is before the current time minus the supplied duration.
func unixTimeStale(timestamp uint64, duration time.Duration) bool {
	return time.Unix(int64(timestamp), 0).Before(time.Now().Add(-1 * duration))
}

func (n *OpNode) P2P() p2p.Node {
	return n.p2pNode
}

func (n *OpNode) RuntimeConfig() ReadonlyRuntimeConfig {
	return n.runCfg
}

// Stop stops the node and closes all resources.
// If the provided ctx is expired, the node will accelerate the stop where possible, but still fully close.
func (n *OpNode) Stop(ctx context.Context) error {
	if n.closed.Load() {
		return ErrAlreadyClosed
	}

	var result *multierror.Error

	if n.server != nil {
		if err := n.server.Stop(ctx); err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to close RPC server: %w", err))
		}
	}
	if n.p2pNode != nil {
		if err := n.p2pNode.Close(); err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to close p2p node: %w", err))
		}
	}
	if n.p2pSigner != nil {
		if err := n.p2pSigner.Close(); err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to close p2p signer: %w", err))
		}
	}

	if n.resourcesClose != nil {
		n.resourcesClose()
	}

	// stop L1 heads feed
	if n.l1HeadsSub != nil {
		n.l1HeadsSub.Unsubscribe()
	}
	// stop polling for L1 safe-head changes
	if n.l1SafeSub != nil {
		n.l1SafeSub.Unsubscribe()
	}
	// stop polling for L1 finalized-head changes
	if n.l1FinalizedSub != nil {
		n.l1FinalizedSub.Unsubscribe()
	}

	// close L2 driver
	if n.l2Driver != nil {
		if err := n.l2Driver.Close(); err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to close L2 engine driver cleanly: %w", err))
		}
	}

	if n.safeDB != nil {
		if err := n.safeDB.Close(); err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to close safe head db: %w", err))
		}
	}

	// Wait for the runtime config loader to be done using the data sources before closing them
	if n.runtimeConfigReloaderDone != nil {
		<-n.runtimeConfigReloaderDone
	}

	// close L2 engine RPC client
	if n.l2Source != nil {
		n.l2Source.Close()
	}

	// close L1 data source
	if n.l1Source != nil {
		n.l1Source.Close()
	}

	if result == nil { // mark as closed if we successfully fully closed
		n.closed.Store(true)
	}

	if n.halted.Load() {
		// if we had a halt upon initialization, idle for a while, with open metrics, to prevent a rapid restart-loop
		tim := time.NewTimer(time.Minute * 5)
		n.log.Warn("halted, idling to avoid immediate shutdown repeats")
		defer tim.Stop()
		select {
		case <-tim.C:
		case <-ctx.Done():
		}
	}

	// Close metrics and pprof only after we are done idling
	if n.pprofService != nil {
		if err := n.pprofService.Stop(ctx); err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to close pprof server: %w", err))
		}
	}
	if n.metricsSrv != nil {
		if err := n.metricsSrv.Stop(ctx); err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to close metrics server: %w", err))
		}
	}

	return result.ErrorOrNil()
}

func (n *OpNode) Stopped() bool {
	return n.closed.Load()
}

func (n *OpNode) HTTPEndpoint() string {
	if n.server == nil {
		return ""
	}
	return fmt.Sprintf("http://%s", n.server.Addr().String())
}
