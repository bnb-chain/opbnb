package batcher

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/ethereum-optimism/optimism/op-batcher/metrics"
	"github.com/ethereum-optimism/optimism/op-service/dial"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum-optimism/optimism/op-service/testutils"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	txmgrmocks "github.com/ethereum-optimism/optimism/op-service/txmgr/mocks"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockL2EndpointProvider struct {
	ethClient       *testutils.MockL2Client
	ethClientErr    error
	rollupClient    *testutils.MockRollupClient
	rollupClientErr error
}

func TestBatchSubmitter_SendTransaction_GasLimit(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		useBlobs bool
		wantGas  uint64
	}{
		{
			name:    "empty frame",
			wantGas: 53_004,
		},
		{
			name:    "intrinsic dominates",
			data:    make([]byte, 3_000),
			wantGas: 65_004,
		},
		{
			name:    "just below floor crossover",
			data:    make([]byte, 5_332),
			wantGas: 74_332,
		},
		{
			name:    "just above floor crossover",
			data:    make([]byte, 5_333),
			wantGas: 74_340,
		},
		{
			name:    "floor dominates zero data",
			data:    make([]byte, 10_000),
			wantGas: 121_010,
		},
		{
			name:    "intrinsic dominates nonzero data",
			data:    bytes.Repeat([]byte{1}, 10_000),
			wantGas: 733_004,
		},
		{
			name:    "floor dominates mixed data",
			data:    append(make([]byte, 10_000), bytes.Repeat([]byte{1}, 1_000)...),
			wantGas: 161_010,
		},
		{
			name:     "blob",
			data:     make([]byte, 10_000),
			useBlobs: true,
			wantGas:  53_000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, _ := setup(t)
			bs.Config.UseBlobs = tt.useBlobs
			manager := txmgrmocks.NewTxManager(t)
			var candidate txmgr.TxCandidate
			manager.On("Send", mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					candidate = args.Get(1).(txmgr.TxCandidate)
				}).
				Return(&types.Receipt{}, nil).Once()

			ctx := context.Background()
			queue := txmgr.NewQueue[txID](ctx, manager, 1)
			receiptsCh := make(chan txmgr.TxReceipt[txID], 1)
			data := singleFrameTxData(frameData{data: tt.data})
			require.NoError(t, bs.sendTransaction(ctx, data, queue, receiptsCh))
			queue.Wait()
			require.Len(t, receiptsCh, 1)
			require.NoError(t, (<-receiptsCh).Err)

			// Calldata includes a zero version byte. Expected limits retain
			// the existing conservative intrinsic-gas schedule.
			require.Equal(t, tt.wantGas, candidate.GasLimit)
		})
	}
}

func newEndpointProvider() *mockL2EndpointProvider {
	return &mockL2EndpointProvider{
		ethClient:    new(testutils.MockL2Client),
		rollupClient: new(testutils.MockRollupClient),
	}
}

func (p *mockL2EndpointProvider) EthClient(context.Context) (dial.EthClientInterface, error) {
	return p.ethClient, p.ethClientErr
}

func (p *mockL2EndpointProvider) RollupClient(context.Context) (dial.RollupClientInterface, error) {
	return p.rollupClient, p.rollupClientErr
}

func (p *mockL2EndpointProvider) Close() {}

const genesisL1Origin = uint64(123)

func setup(t *testing.T) (*BatchSubmitter, *mockL2EndpointProvider) {
	ep := newEndpointProvider()

	cfg := defaultTestRollupConfig
	cfg.Genesis.L1.Number = genesisL1Origin

	return NewBatchSubmitter(DriverSetup{
		Log:              testlog.Logger(t, log.LevelDebug),
		Metr:             metrics.NoopMetrics,
		RollupConfig:     &cfg,
		EndpointProvider: ep,
	}), ep
}

func TestBatchSubmitter_SafeL1Origin(t *testing.T) {
	bs, ep := setup(t)

	tests := []struct {
		name                   string
		currentSafeOrigin      uint64
		failsToFetchSyncStatus bool
		expectResult           uint64
		expectErr              bool
	}{
		{
			name:              "ExistingSafeL1Origin",
			currentSafeOrigin: 999,
			expectResult:      999,
		},
		{
			name:              "NoExistingSafeL1OriginUsesGenesis",
			currentSafeOrigin: 0,
			expectResult:      genesisL1Origin,
		},
		{
			name:                   "ErrorFetchingSyncStatus",
			failsToFetchSyncStatus: true,
			expectErr:              true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.failsToFetchSyncStatus {
				ep.rollupClient.ExpectSyncStatus(&eth.SyncStatus{}, errors.New("failed to fetch sync status"))

			} else {
				ep.rollupClient.ExpectSyncStatus(&eth.SyncStatus{
					SafeL2: eth.L2BlockRef{
						L1Origin: eth.BlockID{
							Number: tt.currentSafeOrigin,
						},
					},
				}, nil)
			}

			id, err := bs.safeL1Origin(context.Background())

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectResult, id.Number)
			}
		})
	}

}

func TestBatchSubmitter_SafeL1Origin_FailsToResolveRollupClient(t *testing.T) {
	bs, ep := setup(t)

	ep.rollupClientErr = errors.New("failed to resolve rollup client")

	_, err := bs.safeL1Origin(context.Background())
	require.Error(t, err)
}
