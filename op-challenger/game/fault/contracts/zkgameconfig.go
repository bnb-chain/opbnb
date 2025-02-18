package contracts

import (
	"context"
	_ "embed"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/contracts/metrics"
	"github.com/ethereum-optimism/optimism/op-service/sources/batching"
	"github.com/ethereum-optimism/optimism/op-service/sources/batching/rpcblock"
	"github.com/ethereum/go-ethereum/common"
)

var (
	methodBlockDistance = "blockDistance"
)

//go:embed abis/zkgameconfig.json
var zkGameConfigAbi []byte

type ZkGameConfigContract struct {
	metrics     metrics.ContractMetricer
	multiCaller *batching.MultiCaller
	contract    *batching.BoundContract
}

func NewZkGameConfig(
	metrics metrics.ContractMetricer,
	addr common.Address,
	caller *batching.MultiCaller,
) ZkGameConfig {
	contractAbi := mustParseAbi(zkGameConfigAbi)
	return &ZkGameConfigContract{
		metrics:     metrics,
		multiCaller: caller,
		contract:    batching.NewBoundContract(contractAbi, addr),
	}
}

func (c *ZkGameConfigContract) GetBlockDistance(ctx context.Context) (*big.Int, error) {
	defer c.metrics.StartContractRequest("GetBlockDistance")()
	result, err := c.multiCaller.SingleCall(ctx, rpcblock.Latest, c.contract.Call(methodBlockDistance))
	if err != nil {
		return nil, fmt.Errorf("failed to get block distance: %w", err)
	}
	return result.GetBigInt(0), nil
}

type ZkGameConfig interface {
	GetBlockDistance(ctx context.Context) (*big.Int, error)
}
