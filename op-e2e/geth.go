package op_e2e

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/catalyst"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/node"
)

func waitForL1OriginOnL2(l1BlockNum uint64, client *ethclient.Client, timeout time.Duration) (*types.Block, error) {
	timeoutCh := time.After(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	headChan := make(chan *types.Header, 100)
	headSub, err := client.SubscribeNewHead(ctx, headChan)
	if err != nil {
		return nil, err
	}
	defer headSub.Unsubscribe()

	for {
		select {
		case head := <-headChan:
			block, err := client.BlockByNumber(ctx, head.Number)
			if err != nil {
				return nil, err
			}
			l1Info, err := derive.L1InfoDepositTxData(block.Transactions()[0].Data())
			if err != nil {
				return nil, err
			}
			if l1Info.Number >= l1BlockNum {
				return block, nil
			}

		case err := <-headSub.Err():
			return nil, fmt.Errorf("error in head subscription: %w", err)
		case <-timeoutCh:
			return nil, errors.New("timeout")
		}
	}
}

func waitForTransaction(hash common.Hash, client *ethclient.Client, timeout time.Duration) (*types.Receipt, error) {
	timeoutCh := time.After(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		receipt, err := client.TransactionReceipt(ctx, hash)
		if receipt != nil && err == nil {
			return receipt, nil
		} else if err != nil && !errors.Is(err, ethereum.NotFound) {
			return nil, err
		}

		select {
		case <-timeoutCh:
			return nil, errors.New("timeout")
		case <-ticker.C:
		}
	}
}

func waitForBlock(number *big.Int, client *ethclient.Client, timeout time.Duration) (*types.Block, error) {
	timeoutCh := time.After(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	headChan := make(chan *types.Header, 100)
	headSub, err := client.SubscribeNewHead(ctx, headChan)
	if err != nil {
		return nil, err
	}
	defer headSub.Unsubscribe()

	for {
		select {
		case head := <-headChan:
			if head.Number.Cmp(number) >= 0 {
				return client.BlockByNumber(ctx, number)
			}
		case err := <-headSub.Err():
			return nil, fmt.Errorf("error in head subscription: %w", err)
		case <-timeoutCh:
			return nil, errors.New("timeout")
		}
	}
}

func initL1Geth(cfg *SystemConfig, genesis *core.Genesis, opts ...GethOption) (*node.Node, *eth.Ethereum, error) {
	ethConfig := &ethconfig.Config{
		NetworkId: cfg.DeployConfig.L1ChainID,
		Genesis:   genesis,
	}
	nodeConfig := &node.Config{
		Name:        "l1-geth",
		HTTPHost:    "127.0.0.1",
		HTTPPort:    0,
		WSHost:      "127.0.0.1",
		WSPort:      0,
		WSModules:   []string{"debug", "admin", "eth", "txpool", "net", "rpc", "web3", "personal", "engine"},
		HTTPModules: []string{"debug", "admin", "eth", "txpool", "net", "rpc", "web3", "personal", "engine"},
	}

	l1Node, l1Eth, err := createGethNode(false, nodeConfig, ethConfig, []*ecdsa.PrivateKey{cfg.Secrets.CliqueSigner}, opts...)
	if err != nil {
		return nil, nil, err
	}
	// Activate merge
	l1Eth.Merger().FinalizePoS()

	// Instead of running a whole beacon node, we run this fake-proof-of-stake sidecar that sequences L1 blocks using the Engine API.
	l1Node.RegisterLifecycle(&fakePoS{
		eth:       l1Eth,
		log:       log.Root(), // geth logger is global anyway. Would be nice to replace with a local logger though.
		blockTime: cfg.DeployConfig.L1BlockTime,
		// for testing purposes we make it really fast, otherwise we don't see it finalize in short tests
		finalizedDistance: 8,
		safeDistance:      4,
		engineAPI:         catalyst.NewConsensusAPI(l1Eth),
	})

	return l1Node, l1Eth, nil
}

// fakePoS is a testing-only utility to attach to Geth,
// to build a fake proof-of-stake L1 chain with fixed block time and basic lagging safe/finalized blocks.
type fakePoS struct {
	eth       *eth.Ethereum
	log       log.Logger
	blockTime uint64

	finalizedDistance uint64
	safeDistance      uint64

	engineAPI *catalyst.ConsensusAPI
	sub       ethereum.Subscription
}

func (f *fakePoS) Start() error {
	f.sub = event.NewSubscription(func(quit <-chan struct{}) error {
		// poll every half a second: enough to catch up with any block time when ticks are missed
		t := time.NewTicker(time.Second / 2)
		for {
			select {
			case now := <-t.C:
				chain := f.eth.BlockChain()
				head := chain.CurrentBlock()
				finalized := chain.CurrentFinalBlock()
				if finalized == nil { // fallback to genesis if nothing is finalized
					finalized = chain.Genesis().Header()
				}
				safe := chain.CurrentSafeBlock()
				if safe == nil { // fallback to finalized if nothing is safe
					safe = finalized
				}
				if head.Number.Uint64() > f.finalizedDistance { // progress finalized block, if we can
					finalized = f.eth.BlockChain().GetHeaderByNumber(head.Number.Uint64() - f.finalizedDistance)
				}
				if head.Number.Uint64() > f.safeDistance { // progress safe block, if we can
					safe = f.eth.BlockChain().GetHeaderByNumber(head.Number.Uint64() - f.safeDistance)
				}
				// start building the block as soon as we are past the current head time
				if head.Time >= uint64(now.Unix()) {
					continue
				}
				res, err := f.engineAPI.ForkchoiceUpdatedV1(engine.ForkchoiceStateV1{
					HeadBlockHash:      head.Hash(),
					SafeBlockHash:      safe.Hash(),
					FinalizedBlockHash: finalized.Hash(),
				}, &engine.PayloadAttributes{
					Timestamp:             head.Time + f.blockTime,
					Random:                common.Hash{},
					SuggestedFeeRecipient: head.Coinbase,
				})
				if err != nil {
					f.log.Error("failed to start building L1 block", "err", err)
					continue
				}
				if res.PayloadID == nil {
					f.log.Error("failed to start block building", "res", res)
					continue
				}
				// wait with sealing, if we are not behind already
				delay := time.Until(time.Unix(int64(head.Time+f.blockTime), 0))
				tim := time.NewTimer(delay)
				select {
				case <-tim.C:
					// no-op
				case <-quit:
					tim.Stop()
					return nil
				}
				payload, err := f.engineAPI.GetPayloadV1(*res.PayloadID)
				if err != nil {
					f.log.Error("failed to finish building L1 block", "err", err)
					continue
				}
				if _, err := f.engineAPI.NewPayloadV1(*payload); err != nil {
					f.log.Error("failed to insert built L1 block", "err", err)
					continue
				}
				if _, err := f.engineAPI.ForkchoiceUpdatedV1(engine.ForkchoiceStateV1{
					HeadBlockHash:      payload.BlockHash,
					SafeBlockHash:      safe.Hash(),
					FinalizedBlockHash: finalized.Hash(),
				}, nil); err != nil {
					f.log.Error("failed to make built L1 block canonical", "err", err)
					continue
				}
			case <-quit:
				return nil
			}
		}
	})
	return nil
}

func (f *fakePoS) Stop() error {
	f.sub.Unsubscribe()
	return nil
}

func defaultNodeConfig(name string, jwtPath string) *node.Config {
	return &node.Config{
		Name:        name,
		WSHost:      "127.0.0.1",
		WSPort:      0,
		AuthAddr:    "127.0.0.1",
		AuthPort:    0,
		HTTPHost:    "127.0.0.1",
		HTTPPort:    0,
		WSModules:   []string{"debug", "admin", "eth", "txpool", "net", "rpc", "web3", "personal", "engine"},
		HTTPModules: []string{"debug", "admin", "eth", "txpool", "net", "rpc", "web3", "personal", "engine"},
		JWTSecret:   jwtPath,
	}
}

type GethOption func(ethCfg *ethconfig.Config, nodeCfg *node.Config) error

// init a geth node.
func initL2Geth(name string, l2ChainID *big.Int, genesis *core.Genesis, jwtPath string, opts ...GethOption) (*node.Node, *eth.Ethereum, error) {
	ethConfig := &ethconfig.Config{
		NetworkId: l2ChainID.Uint64(),
		Genesis:   genesis,
		Miner: miner.Config{
			Etherbase:         common.Address{},
			Notify:            nil,
			NotifyFull:        false,
			ExtraData:         nil,
			GasFloor:          0,
			GasCeil:           0,
			GasPrice:          nil,
			Recommit:          0,
			Noverify:          false,
			NewPayloadTimeout: 0,
		},
	}
	nodeConfig := defaultNodeConfig(fmt.Sprintf("l2-geth-%v", name), jwtPath)
	return createGethNode(true, nodeConfig, ethConfig, nil, opts...)
}

// createGethNode creates an in-memory geth node based on the configuration.
// The private keys are added to the keystore and are unlocked.
// If the node is l2, catalyst is enabled.
// The node should be started and then closed when done.
func createGethNode(l2 bool, nodeCfg *node.Config, ethCfg *ethconfig.Config, privateKeys []*ecdsa.PrivateKey, opts ...GethOption) (*node.Node, *eth.Ethereum, error) {
	for i, opt := range opts {
		if err := opt(ethCfg, nodeCfg); err != nil {
			return nil, nil, fmt.Errorf("failed to apply geth option %d: %w", i, err)
		}
	}
	ethCfg.NoPruning = true // force everything to be an archive node
	n, err := node.New(nodeCfg)
	if err != nil {
		n.Close()
		return nil, nil, err
	}

	keydir := n.KeyStoreDir()
	scryptN := 2
	scryptP := 1
	n.AccountManager().AddBackend(keystore.NewKeyStore(keydir, scryptN, scryptP))
	ks := n.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	password := "foobar"
	for _, pk := range privateKeys {
		act, err := ks.ImportECDSA(pk, password)
		if err != nil {
			n.Close()
			return nil, nil, err
		}
		err = ks.Unlock(act, password)
		if err != nil {
			n.Close()
			return nil, nil, err
		}
	}

	backend, err := eth.New(n, ethCfg)
	if err != nil {
		n.Close()
		return nil, nil, err

	}

	// PR 25459 changed this to only default in CLI, but not in default programmatic RPC selection.
	// PR 25642 fixed it for the mobile version only...
	utils.RegisterFilterAPI(n, backend.APIBackend, ethCfg)

	n.RegisterAPIs(tracers.APIs(backend.APIBackend))

	// Enable catalyst if l2
	if l2 {
		if err := catalyst.Register(n, backend); err != nil {
			n.Close()
			return nil, nil, err
		}
	}
	return n, backend, nil

}
