package contracts

import (
	"fmt"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/contracts/metrics"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/sources/batching"
	"github.com/ethereum-optimism/optimism/packages/contracts-bedrock/snapshots"
	"github.com/ethereum/go-ethereum/common"
)

type ZkGameFactory struct {
	DisputeGameFactoryContract
}

func NewZkDisputeGameFactoryContract(
	m metrics.ContractMetricer,
	addr common.Address,
	caller *batching.MultiCaller,
) *ZkGameFactory {
	factoryAbi := snapshots.LoadZkDisputeGameFactoryABI()
	return &ZkGameFactory{
		DisputeGameFactoryContract: DisputeGameFactoryContract{
			metrics:     m,
			multiCaller: caller,
			contract:    batching.NewBoundContract(factoryAbi, addr),
			abi:         factoryAbi,
		},
	}
}

type ZkGameCreateCallData struct {
	GameType      uint32
	Claims        []eth.Bytes32
	ParentGameIdx uint64
	L2BlockNumber uint64
	ExtraData     []byte
}

func (f *ZkGameFactory) DecodeZKGameCreateCallData(data []byte) (*ZkGameCreateCallData, error) {
	name, callResult, err := f.contract.DecodeCall(data)
	if err != nil {
		return nil, err
	}
	if name != "createZkFaultDisputeGame" {
		return nil, fmt.Errorf("invalid name %s", name)
	}
	bytes32Slice := callResult.GetBytes32Slice(1)
	var claims []eth.Bytes32
	for _, oneBytes32 := range bytes32Slice {
		claims = append(claims, oneBytes32)
	}
	result := &ZkGameCreateCallData{
		GameType:      callResult.GetUint32(0),
		Claims:        claims,
		ParentGameIdx: callResult.GetUint64(2),
		L2BlockNumber: callResult.GetUint64(3),
		ExtraData:     callResult.GetBytes(4),
	}
	return result, nil
}
