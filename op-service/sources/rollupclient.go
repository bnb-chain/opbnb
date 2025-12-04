package sources

import (
	"context"
	"io"

	"github.com/ethereum-optimism/optimism/op-service/sources/batching"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/exp/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

type RollupClient struct {
	rpc client.RPC
}

func NewRollupClient(rpc client.RPC) *RollupClient {
	return &RollupClient{rpc}
}

func (r *RollupClient) OutputAtBlock(ctx context.Context, blockNum uint64) (*eth.OutputResponse, error) {
	var output *eth.OutputResponse
	err := r.rpc.CallContext(ctx, &output, "optimism_outputAtBlock", hexutil.Uint64(blockNum))
	return output, err
}

func (r *RollupClient) BatchOutputAtBlock(ctx context.Context, blocks []uint64) ([]*eth.OutputResponse, error) {
	var result []*eth.OutputResponse
	batchCall := batching.NewIterativeBatchCall[uint64, *eth.OutputResponse](blocks, func(block uint64) (*eth.OutputResponse, rpc.BatchElem) {
		var response eth.OutputResponse
		elem := rpc.BatchElem{
			Method: "optimism_outputAtBlock",
			Args:   []interface{}{hexutil.Uint64(block)},
			Result: &response,
		}
		result = append(result, &response)
		return &response, elem
	}, r.rpc.BatchCallContext, r.rpc.CallContext, 100)
	for {
		if err := batchCall.Fetch(ctx); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	return batchCall.Result()
}

func (r *RollupClient) SafeHeadAtL1Block(ctx context.Context, blockNum uint64) (*eth.SafeHeadResponse, error) {
	var output *eth.SafeHeadResponse
	err := r.rpc.CallContext(ctx, &output, "optimism_safeHeadAtL1Block", hexutil.Uint64(blockNum))
	return output, err
}

func (r *RollupClient) SyncStatus(ctx context.Context) (*eth.SyncStatus, error) {
	var output *eth.SyncStatus
	err := r.rpc.CallContext(ctx, &output, "optimism_syncStatus")
	return output, err
}

func (r *RollupClient) RollupConfig(ctx context.Context) (*rollup.Config, error) {
	var output *rollup.Config
	err := r.rpc.CallContext(ctx, &output, "optimism_rollupConfig")
	return output, err
}

func (r *RollupClient) Version(ctx context.Context) (string, error) {
	var output string
	err := r.rpc.CallContext(ctx, &output, "optimism_version")
	return output, err
}

func (r *RollupClient) StartSequencer(ctx context.Context, unsafeHead common.Hash) error {
	return r.rpc.CallContext(ctx, nil, "admin_startSequencer", unsafeHead)
}

func (r *RollupClient) StopSequencer(ctx context.Context) (common.Hash, error) {
	var result common.Hash
	err := r.rpc.CallContext(ctx, &result, "admin_stopSequencer")
	return result, err
}

func (r *RollupClient) SequencerActive(ctx context.Context) (bool, error) {
	var result bool
	err := r.rpc.CallContext(ctx, &result, "admin_sequencerActive")
	return result, err
}

func (r *RollupClient) PostUnsafePayload(ctx context.Context, payload *eth.ExecutionPayloadEnvelope) error {
	return r.rpc.CallContext(ctx, nil, "admin_postUnsafePayload", payload)
}

func (r *RollupClient) SetLogLevel(ctx context.Context, lvl slog.Level) error {
	return r.rpc.CallContext(ctx, nil, "admin_setLogLevel", lvl.String())
}

func (r *RollupClient) Close() {
	r.rpc.Close()
}
