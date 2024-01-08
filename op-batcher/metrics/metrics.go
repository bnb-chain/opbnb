package metrics

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
	txmetrics "github.com/ethereum-optimism/optimism/op-service/txmgr/metrics"
)

const Namespace = "op_batcher"
const RPCClientSubsystem = "rpc_client"

type Metricer interface {
	RecordInfo(version string)
	RecordUp()

	// Records all L1 and L2 block events
	opmetrics.RefMetricer

	// Record Tx metrics
	txmetrics.TxMetricer

	RecordLatestL1Block(l1ref eth.L1BlockRef)
	RecordL2BlocksLoaded(l2ref eth.L2BlockRef)
	RecordChannelOpened(id derive.ChannelID, numPendingBlocks int)
	RecordL2BlocksAdded(l2ref eth.L2BlockRef, numBlocksAdded, numPendingBlocks, inputBytes, outputComprBytes int)
	RecordL2BlockInPendingQueue(block *types.Block)
	RecordL2BlockInChannel(block *types.Block)
	RecordChannelClosed(id derive.ChannelID, numPendingBlocks int, numFrames int, inputBytes int, outputComprBytes int, reason error)
	RecordChannelFullySubmitted(id derive.ChannelID)
	RecordChannelTimedOut(id derive.ChannelID)

	RecordBatchTxSubmitted()
	RecordBatchTxSuccess()
	RecordBatchTxFailed()

	Document() []opmetrics.DocumentedMetric
	client.Metricer
}

type Metrics struct {
	ns       string
	registry *prometheus.Registry
	factory  opmetrics.Factory

	opmetrics.RefMetrics
	txmetrics.TxMetrics

	info prometheus.GaugeVec
	up   prometheus.Gauge

	// label by openend, closed, fully_submitted, timed_out
	channelEvs opmetrics.EventVec

	pendingBlocksCount        prometheus.GaugeVec
	pendingBlocksBytesTotal   prometheus.Counter
	pendingBlocksBytesCurrent prometheus.Gauge
	blocksAddedCount          prometheus.Gauge

	channelInputBytes       prometheus.GaugeVec
	channelReadyBytes       prometheus.Gauge
	channelOutputBytes      prometheus.Gauge
	channelClosedReason     prometheus.Gauge
	channelNumFrames        prometheus.Gauge
	channelComprRatio       prometheus.Histogram
	channelInputBytesTotal  prometheus.Counter
	channelOutputBytesTotal prometheus.Counter

	batcherTxEvs opmetrics.EventVec

	RPCClientRequestsTotal          *prometheus.CounterVec
	RPCClientRequestDurationSeconds *prometheus.HistogramVec
	RPCClientResponsesTotal         *prometheus.CounterVec
}

var _ Metricer = (*Metrics)(nil)

func NewMetrics(procName string) *Metrics {
	if procName == "" {
		procName = "default"
	}
	ns := Namespace + "_" + procName

	registry := opmetrics.NewRegistry()
	factory := opmetrics.With(registry)

	return &Metrics{
		ns:       ns,
		registry: registry,
		factory:  factory,

		RefMetrics: opmetrics.MakeRefMetrics(ns, factory),
		TxMetrics:  txmetrics.MakeTxMetrics(ns, factory),

		info: *factory.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "info",
			Help:      "Pseudo-metric tracking version and config info",
		}, []string{
			"version",
		}),
		up: factory.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "up",
			Help:      "1 if the op-batcher has finished starting up",
		}),

		channelEvs: opmetrics.NewEventVec(factory, ns, "", "channel", "Channel", []string{"stage"}),

		pendingBlocksCount: *factory.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "pending_blocks_count",
			Help:      "Number of pending blocks, not added to a channel yet.",
		}, []string{"stage"}),
		pendingBlocksBytesTotal: factory.NewCounter(prometheus.CounterOpts{
			Namespace: ns,
			Name:      "pending_blocks_bytes_total",
			Help:      "Total size of transactions in pending blocks as they are fetched from L2",
		}),
		pendingBlocksBytesCurrent: factory.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "pending_blocks_bytes_current",
			Help:      "Current size of transactions in the pending (fetched from L2 but not in a channel) stage.",
		}),
		blocksAddedCount: factory.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "blocks_added_count",
			Help:      "Total number of blocks added to current channel.",
		}),

		channelInputBytes: *factory.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "input_bytes",
			Help:      "Number of input bytes to a channel.",
		}, []string{"stage"}),
		channelReadyBytes: factory.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "ready_bytes",
			Help:      "Number of bytes ready in the compression buffer.",
		}),
		channelOutputBytes: factory.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "output_bytes",
			Help:      "Number of compressed output bytes from a channel.",
		}),
		channelClosedReason: factory.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "channel_closed_reason",
			Help:      "Pseudo-metric to record the reason a channel got closed.",
		}),
		channelNumFrames: factory.NewGauge(prometheus.GaugeOpts{
			Namespace: ns,
			Name:      "channel_num_frames",
			Help:      "Total number of frames of closed channel.",
		}),
		channelComprRatio: factory.NewHistogram(prometheus.HistogramOpts{
			Namespace: ns,
			Name:      "channel_compr_ratio",
			Help:      "Compression ratios of closed channel.",
			Buckets:   append([]float64{0.1, 0.2}, prometheus.LinearBuckets(0.3, 0.05, 14)...),
		}),
		channelInputBytesTotal: factory.NewCounter(prometheus.CounterOpts{
			Namespace: ns,
			Name:      "input_bytes_total",
			Help:      "Total number of bytes to a channel.",
		}),
		channelOutputBytesTotal: factory.NewCounter(prometheus.CounterOpts{
			Namespace: ns,
			Name:      "output_bytes_total",
			Help:      "Total number of compressed output bytes from a channel.",
		}),

		batcherTxEvs: opmetrics.NewEventVec(factory, ns, "", "batcher_tx", "BatcherTx", []string{"stage"}),

		RPCClientRequestsTotal: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: ns,
			Subsystem: RPCClientSubsystem,
			Name:      "requests_total",
			Help:      "Total RPC requests initiated by the op-batcher's RPC client",
		}, []string{
			"method",
		}),
		RPCClientRequestDurationSeconds: factory.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: ns,
			Subsystem: RPCClientSubsystem,
			Name:      "request_duration_seconds",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			Help:      "Histogram of RPC client request durations",
		}, []string{
			"method",
		}),
		RPCClientResponsesTotal: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: ns,
			Subsystem: RPCClientSubsystem,
			Name:      "responses_total",
			Help:      "Total RPC request responses received by the op-batcher's RPC client",
		}, []string{
			"method",
			"error",
		}),
	}
}

func (m *Metrics) Serve(ctx context.Context, host string, port int) error {
	return opmetrics.ListenAndServe(ctx, m.registry, host, port)
}

func (m *Metrics) Document() []opmetrics.DocumentedMetric {
	return m.factory.Document()
}

func (m *Metrics) StartBalanceMetrics(ctx context.Context,
	l log.Logger, client ethereum.ChainStateReader, account common.Address) {
	opmetrics.LaunchBalanceMetrics(ctx, l, m.registry, m.ns, client, account)
}

// RecordInfo sets a pseudo-metric that contains versioning and
// config info for the op-batcher.
func (m *Metrics) RecordInfo(version string) {
	m.info.WithLabelValues(version).Set(1)
}

// RecordUp sets the up metric to 1.
func (m *Metrics) RecordUp() {
	prometheus.MustRegister()
	m.up.Set(1)
}

const (
	StageLoaded         = "loaded"
	StageOpened         = "opened"
	StageAdded          = "added"
	StageClosed         = "closed"
	StageFullySubmitted = "fully_submitted"
	StageTimedOut       = "timed_out"

	TxStageSubmitted = "submitted"
	TxStageSuccess   = "success"
	TxStageFailed    = "failed"
)

func (m *Metrics) RecordLatestL1Block(l1ref eth.L1BlockRef) {
	m.RecordL1Ref("latest", l1ref)
}

// RecordL2BlocksLoaded should be called when a new L2 block was loaded into the
// channel manager (but not processed yet).
func (m *Metrics) RecordL2BlocksLoaded(l2ref eth.L2BlockRef) {
	m.RecordL2Ref(StageLoaded, l2ref)
}

func (m *Metrics) RecordChannelOpened(id derive.ChannelID, numPendingBlocks int) {
	m.channelEvs.Record(StageOpened)
	m.blocksAddedCount.Set(0) // reset
	m.pendingBlocksCount.WithLabelValues(StageOpened).Set(float64(numPendingBlocks))
}

// RecordL2BlocksAdded should be called when L2 block were added to the channel
// builder, with the latest added block.
func (m *Metrics) RecordL2BlocksAdded(l2ref eth.L2BlockRef, numBlocksAdded, numPendingBlocks, inputBytes, outputComprBytes int) {
	m.RecordL2Ref(StageAdded, l2ref)
	m.blocksAddedCount.Add(float64(numBlocksAdded))
	m.pendingBlocksCount.WithLabelValues(StageAdded).Set(float64(numPendingBlocks))
	m.channelInputBytes.WithLabelValues(StageAdded).Set(float64(inputBytes))
	m.channelReadyBytes.Set(float64(outputComprBytes))
}

func (m *Metrics) RecordChannelClosed(id derive.ChannelID, numPendingBlocks int, numFrames int, inputBytes int, outputComprBytes int, reason error) {
	m.channelEvs.Record(StageClosed)
	m.pendingBlocksCount.WithLabelValues(StageClosed).Set(float64(numPendingBlocks))
	m.channelNumFrames.Set(float64(numFrames))
	m.channelInputBytes.WithLabelValues(StageClosed).Set(float64(inputBytes))
	m.channelOutputBytes.Set(float64(outputComprBytes))
	m.channelInputBytesTotal.Add(float64(inputBytes))
	m.channelOutputBytesTotal.Add(float64(outputComprBytes))

	var comprRatio float64
	if inputBytes > 0 {
		comprRatio = float64(outputComprBytes) / float64(inputBytes)
	}
	m.channelComprRatio.Observe(comprRatio)

	m.channelClosedReason.Set(float64(ClosedReasonToNum(reason)))
}

func (m *Metrics) RecordL2BlockInPendingQueue(block *types.Block) {
	size := float64(estimateBatchSize(block))
	m.pendingBlocksBytesTotal.Add(size)
	m.pendingBlocksBytesCurrent.Add(size)
}

func (m *Metrics) RecordL2BlockInChannel(block *types.Block) {
	size := float64(estimateBatchSize(block))
	m.pendingBlocksBytesCurrent.Add(-1 * size)
	// Refer to RecordL2BlocksAdded to see the current + count of bytes added to a channel
}

func ClosedReasonToNum(reason error) int {
	// CLI-3640
	return 0
}

func (m *Metrics) RecordChannelFullySubmitted(id derive.ChannelID) {
	m.channelEvs.Record(StageFullySubmitted)
}

func (m *Metrics) RecordChannelTimedOut(id derive.ChannelID) {
	m.channelEvs.Record(StageTimedOut)
}

func (m *Metrics) RecordBatchTxSubmitted() {
	m.batcherTxEvs.Record(TxStageSubmitted)
}

func (m *Metrics) RecordBatchTxSuccess() {
	m.batcherTxEvs.Record(TxStageSuccess)
}

func (m *Metrics) RecordBatchTxFailed() {
	m.batcherTxEvs.Record(TxStageFailed)
}

func (m *Metrics) RecordRPCClientRequest(method string) func(err error) {
	m.RPCClientRequestsTotal.WithLabelValues(method).Inc()
	timer := prometheus.NewTimer(m.RPCClientRequestDurationSeconds.WithLabelValues(method))
	return func(err error) {
		m.RecordRPCClientResponse(method, err)
		timer.ObserveDuration()
	}
}

// RecordRPCClientResponse records an RPC response. It will
// convert the passed-in error into something metrics friendly.
// Nil errors get converted into <nil>, RPC errors are converted
// into rpc_<error code>, HTTP errors are converted into
// http_<status code>, and everything else is converted into
// <unknown>.
func (m *Metrics) RecordRPCClientResponse(method string, err error) {
	var errStr string
	var rpcErr rpc.Error
	var httpErr rpc.HTTPError
	if err == nil {
		errStr = "<nil>"
	} else if errors.As(err, &rpcErr) {
		errStr = fmt.Sprintf("rpc_%d", rpcErr.ErrorCode())
	} else if errors.As(err, &httpErr) {
		errStr = fmt.Sprintf("http_%d", httpErr.StatusCode)
	} else if errors.Is(err, ethereum.NotFound) {
		errStr = "<not found>"
	} else {
		errStr = "<unknown>"
	}
	m.RPCClientResponsesTotal.WithLabelValues(method, errStr).Inc()
}

// estimateBatchSize estimates the size of the batch
func estimateBatchSize(block *types.Block) uint64 {
	size := uint64(70) // estimated overhead of batch metadata
	for _, tx := range block.Transactions() {
		// Don't include deposit transactions in the batch.
		if tx.IsDepositTx() {
			continue
		}
		// Add 2 for the overhead of encoding the tx bytes in a RLP list
		size += tx.Size() + 2
	}
	return size
}
