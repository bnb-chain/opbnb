package actions

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	plasma "github.com/ethereum-optimism/optimism/op-plasma"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func setupFallbackClientTest(t Testing, sd *e2eutils.SetupData, log log.Logger, l1Url string) (*L1Miner, *L1Replica, *L1Replica, *L2Engine, *L2Sequencer, *sources.FallbackClient) {
	jwtPath := e2eutils.WriteDefaultJWT(t)

	miner := NewL1MinerWithPort(t, log, sd.L1Cfg, 8545)
	l1_2 := NewL1ReplicaWithPort(t, log, sd.L1Cfg, 8546)
	l1_3 := NewL1ReplicaWithPort(t, log, sd.L1Cfg, 8547)
	isMultiUrl, urlList := client.MultiUrlParse(l1Url)
	require.True(t, isMultiUrl)
	opts := []client.RPCOption{
		client.WithHttpPollInterval(0),
		client.WithDialBackoff(10),
	}
	rpc, err := client.NewRPC(t.Ctx(), log, urlList[0], opts...)
	require.NoError(t, err)
	fallbackClient := sources.NewFallbackClient(t.Ctx(), rpc, urlList, log, sd.RollupCfg.L1ChainID, sd.RollupCfg.Genesis.L1, func(url string) (client.RPC, error) {
		return client.NewRPC(t.Ctx(), log, url, opts...)
	})
	l1F, err := sources.NewL1Client(fallbackClient, log, nil, sources.L1ClientDefaultConfig(sd.RollupCfg, false, sources.RPCKindBasic))
	require.NoError(t, err)
	l1Blob := sources.NewBSCBlobClient([]client.RPC{rpc})
	engine := NewL2Engine(t, log, sd.L2Cfg, sd.RollupCfg.Genesis.L1, jwtPath)
	l2Cl, err := sources.NewEngineClient(engine.RPCClient(), log, nil, sources.EngineClientDefaultConfig(sd.RollupCfg))
	require.NoError(t, err)

	sequencer := NewL2Sequencer(t, log, l1F, l1Blob, plasma.Disabled, l2Cl, sd.RollupCfg, 0)
	return miner, l1_2, l1_3, engine, sequencer, fallbackClient.(*sources.FallbackClient)
}

func TestL1FallbackClient_SwitchUrl(gt *testing.T) {
	t := NewDefaultTesting(gt)
	p := &e2eutils.TestParams{
		MaxSequencerDrift:   300,
		SequencerWindowSize: 200,
		ChannelTimeout:      120,
		L1BlockTime:         3,
	}
	dp := e2eutils.MakeDeployParams(t, p)
	sd := e2eutils.Setup(t, dp, defaultAlloc)
	logT := testlog.Logger(t, log.LvlDebug)
	miner, l1_2, _, engine, sequencer, fallbackClient := setupFallbackClientTest(t, sd, logT, "http://127.0.0.1:8545,http://127.0.0.1:8546,http://127.0.0.1:8547")
	miner.ActL1SetFeeRecipient(common.Address{'A'})

	sequencer.ActL2PipelineFull(t)

	signer := types.LatestSigner(sd.L2Cfg.Config)
	cl := engine.EthClient()
	aliceTx := func() {
		n, err := cl.PendingNonceAt(t.Ctx(), dp.Addresses.Alice)
		require.NoError(t, err)
		tx := types.MustSignNewTx(dp.Secrets.Alice, signer, &types.DynamicFeeTx{
			ChainID:   sd.L2Cfg.Config.ChainID,
			Nonce:     n,
			GasTipCap: big.NewInt(2 * params.GWei),
			GasFeeCap: new(big.Int).Add(miner.l1Chain.CurrentBlock().BaseFee, big.NewInt(2*params.GWei)),
			Gas:       params.TxGas,
			To:        &dp.Addresses.Bob,
			Value:     e2eutils.Ether(2),
		})
		require.NoError(gt, cl.SendTransaction(t.Ctx(), tx))
	}
	makeL2BlockWithAliceTx := func() {
		aliceTx()
		sequencer.ActL2StartBlock(t)
		engine.ActL2IncludeTx(dp.Addresses.Alice)(t) // include a test tx from alice
		sequencer.ActL2EndBlock(t)
	}

	errRpc := miner.RPCClient().CallContext(t.Ctx(), nil, "admin_stopHTTP")
	require.NoError(t, errRpc)

	l2BlockCount := 0
	for i := 0; i < 8; i++ {
		miner.ActL1StartBlock(3)(t)
		miner.ActL1EndBlock(t)
		newBlock := miner.l1Chain.GetBlockByHash(miner.l1Chain.CurrentBlock().Hash())
		_, err := l1_2.l1Chain.InsertChain([]*types.Block{newBlock})
		require.NoError(t, err)

		sequencer.L2Verifier.l1State.HandleNewL1HeadBlock(eth.L1BlockRef{
			Hash:       newBlock.Hash(),
			Number:     newBlock.NumberU64(),
			ParentHash: newBlock.ParentHash(),
			Time:       newBlock.Time(),
		})
		origin := miner.l1Chain.CurrentBlock()

		for sequencer.SyncStatus().UnsafeL2.Time+sd.RollupCfg.BlockTime < origin.Time {
			makeL2BlockWithAliceTx()
			//require.Equal(t, uint64(i), sequencer.SyncStatus().UnsafeL2.L1Origin.Number, "no L1 origin change before time matches")
			l2BlockCount++
			if l2BlockCount == 11 {
				require.Equal(t, 1, fallbackClient.GetCurrentIndex(), "fallback client should switch url to second url")
				errRpc2 := miner.RPCClient().CallContext(t.Ctx(), nil, "admin_startHTTP", "127.0.0.1", 8545, "*", "eth,net,web3,debug,admin,txpool", "*")
				require.NoError(t, errRpc2)
			}
			if l2BlockCount == 17 {
				require.Equal(t, 0, fallbackClient.GetCurrentIndex(), "fallback client should recover url to first url")
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}
