package derive

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"

	"github.com/ethereum-optimism/optimism/op-node/metrics"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/sync"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/testutils"
)

func TestInsertUnsafePayloadELSyncTrigger(t *testing.T) {
	genesis := eth.L2BlockRef{Hash: common.Hash{0x01}}
	engineUnsafe := eth.L2BlockRef{Hash: common.Hash{0x64}, Number: 100}
	engineSafe := eth.L2BlockRef{Hash: common.Hash{0x62}, Number: 98}
	engineFinalized := eth.L2BlockRef{Hash: common.Hash{0x5f}, Number: 95}
	cfg := &rollup.Config{Genesis: rollup.Genesis{L2: genesis.ID()}}

	for _, tc := range []struct {
		name          string
		payloadNumber uint64
		finalized     eth.L2BlockRef
		finalizedErr  error
		wantELSync    bool
	}{
		{name: "replayed", payloadNumber: 90, finalized: engineFinalized},
		{name: "same_height", payloadNumber: 100, finalized: engineFinalized},
		{name: "within_gap", payloadNumber: 109, finalized: engineFinalized},
		{name: "at_gap", payloadNumber: 110, finalized: engineFinalized},
		{name: "above_gap", payloadNumber: 111, finalized: engineFinalized, wantELSync: true},
		{name: "missing_finalized", payloadNumber: 105, finalizedErr: ethereum.NotFound, wantELSync: true},
		{name: "genesis_finalized", payloadNumber: 105, finalized: genesis, wantELSync: true},
	} {
		for _, cached := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/cached_%t", tc.name, cached), func(t *testing.T) {
				engine := &testutils.MockEngine{}
				engine.Test(t)
				t.Cleanup(func() { engine.AssertExpectations(t) })
				controller := NewEngineController(engine, log.NewLogger(log.DiscardHandler()), metrics.NoopMetrics, cfg,
					&sync.Config{SyncMode: sync.ELSync, ELTriggerGap: 10}, false)
				if cached {
					controller.SetUnsafeHead(engineUnsafe)
					controller.SetSafeHead(engineSafe)
					controller.SetFinalizedHead(engineFinalized)
				}
				initialUnsafe := controller.UnsafeL2Head()
				initialSafe := controller.SafeL2Head()
				initialFinalized := controller.Finalized()

				ref := eth.L2BlockRef{Hash: common.Hash{byte(tc.payloadNumber)}, Number: tc.payloadNumber}
				payload := &eth.ExecutionPayload{BlockHash: ref.Hash, BlockNumber: eth.Uint64Quantity(ref.Number)}
				engine.ExpectL2BlockRefByLabel(eth.Finalized, tc.finalized, tc.finalizedErr)
				if !cached {
					// Before the initial derivation reset, the real unsafe head is only in the engine.
					engine.ExpectL2BlockRefByLabel(eth.Unsafe, engineUnsafe, nil)
				}
				if tc.wantELSync {
					engine.ExpectNewPayload(payload, nil, &eth.PayloadStatusV1{Status: eth.ExecutionSyncing}, nil)
					engine.ExpectForkchoiceUpdate(&eth.ForkchoiceState{
						HeadBlockHash:      ref.Hash,
						SafeBlockHash:      initialSafe.Hash,
						FinalizedBlockHash: initialFinalized.Hash,
					}, nil, &eth.ForkchoiceUpdatedResult{
						PayloadStatus: eth.PayloadStatusV1{Status: eth.ExecutionSyncing},
					}, nil)
				}
				// Skipped EL sync must not submit a payload or a forkchoice update.
				require.NoError(t, controller.InsertUnsafePayload(context.Background(),
					&eth.ExecutionPayloadEnvelope{ExecutionPayload: payload}, ref))
				require.Equal(t, tc.wantELSync, controller.IsEngineSyncing())
				require.Equal(t, initialUnsafe, controller.UnsafeL2Head())
				require.Equal(t, initialSafe, controller.SafeL2Head())
				require.Equal(t, initialFinalized, controller.Finalized())
			})
		}
	}
}
