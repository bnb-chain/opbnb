package contracts

import (
	"context"
	_ "embed"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/contracts/metrics"
	gameTypes "github.com/ethereum-optimism/optimism/op-challenger/game/types"
	"github.com/ethereum-optimism/optimism/op-service/sources/batching"
	"github.com/ethereum-optimism/optimism/op-service/sources/batching/rpcblock"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
)

//go:embed abis/ZKFaultDisputeGame.json
var zkFaultDisputeGameAbi []byte

type ZKFaultDisputeGameContract struct {
	metrics     metrics.ContractMetricer
	multiCaller *batching.MultiCaller
	contract    *batching.BoundContract
}

func NewZKFaultDisputeGameContract(
	ctx context.Context,
	metrics metrics.ContractMetricer,
	addr common.Address,
	caller *batching.MultiCaller,
) (ZKFaultDisputeGame, error) {
	contractAbi := mustParseAbi(zkFaultDisputeGameAbi)

	result, err := caller.SingleCall(ctx, rpcblock.Latest, batching.NewContractCall(contractAbi, addr, methodVersion))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve version of dispute game %v: %w", addr, err)
	}
	version := result.GetString(0)
	if strings.HasPrefix(version, "1.2.") {
		// Detected an older version of contracts, use a compatibility shim.
		return &ZKFaultDisputeGameContract{
			metrics:     metrics,
			multiCaller: caller,
			contract:    batching.NewBoundContract(contractAbi, addr),
		}, nil
	} else {
		return nil, fmt.Errorf("zk fault dispute game contract version %v does not start with '1.2.',addr:%s", version, addr)
	}
}

func (z *ZKFaultDisputeGameContract) GetStatus(ctx context.Context) (gameTypes.GameStatus, error) {
	defer z.metrics.StartContractRequest("GetStatus")()
	result, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, z.contract.Call(methodStatus))
	if err != nil {
		return 0, fmt.Errorf("failed to fetch status: %w", err)
	}
	return gameTypes.GameStatusFromUint8(result.GetUint8(0))
}

func (z *ZKFaultDisputeGameContract) GetClaimCount(ctx context.Context) (uint64, error) {
	defer z.metrics.StartContractRequest("GetClaimCount")()
	result, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, z.contract.Call(methodClaimLen))
	if err != nil {
		return 0, fmt.Errorf("failed to fetch claim count: %w", err)
	}
	return result.GetBigInt(0).Uint64(), nil
}

func (z *ZKFaultDisputeGameContract) GetL1Head(ctx context.Context) (common.Hash, error) {
	defer z.metrics.StartContractRequest("GetL1Head")()
	result, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, z.contract.Call(methodL1Head))
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to fetch L1 head: %w", err)
	}
	return result.GetHash(0), nil
}

func (z *ZKFaultDisputeGameContract) GetBlockRange(ctx context.Context) (prestateBlock uint64, poststateBlock uint64, retErr error) {
	defer z.metrics.StartContractRequest("GetBlockRange")()
	results, err := z.multiCaller.Call(ctx, rpcblock.Latest,
		z.contract.Call(methodStartingBlockNumber),
		z.contract.Call(methodL2BlockNumber))
	if err != nil {
		retErr = fmt.Errorf("failed to retrieve game block range: %w", err)
		return
	}
	if len(results) != 2 {
		retErr = fmt.Errorf("expected 2 results but got %v", len(results))
		return
	}
	prestateBlock = results[0].GetBigInt(0).Uint64()
	poststateBlock = results[1].GetBigInt(0).Uint64()
	return
}

func (z *ZKFaultDisputeGameContract) GetMaxDetectFaultDuration(ctx context.Context) (time.Duration, error) {
	defer z.metrics.StartContractRequest("GetMaxDetectFaultDuration")()
	result, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, z.contract.Call(methodMaxDetectFaultDuration))
	if err != nil {
		return time.Duration(0), fmt.Errorf("failed to fetch maxDetectFaultDuration: %w", err)
	}
	return time.Duration(result.GetUint64(0)), nil
}

func (z *ZKFaultDisputeGameContract) GetCreatedAt(ctx context.Context) (time.Time, error) {
	defer z.metrics.StartContractRequest("GetCreatedAt")()
	result, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, z.contract.Call(methodCreateAt))
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to fetch createAt: %w", err)
	}

	return time.Unix(int64(result.GetUint64(0)), 0), nil
}

func (z *ZKFaultDisputeGameContract) GetClaimsHash(ctx context.Context) (common.Hash, error) {
	defer z.metrics.StartContractRequest("GetClaimsHash")()
	result, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, z.contract.Call(methodClaimsHash))
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to fetch claims hash: %w", err)
	}
	return result.GetHash(0), nil
}

func (z *ZKFaultDisputeGameContract) GetRootClaim(ctx context.Context) (common.Hash, error) {
	defer z.metrics.StartContractRequest("GetRootClaim")()
	result, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, z.contract.Call(methodRootClaim))
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to fetch claims hash: %w", err)
	}
	return result.GetHash(0), nil
}

func (z *ZKFaultDisputeGameContract) GetChallengedClaims(ctx context.Context, targetIdx int) (bool, error) {
	defer z.metrics.StartContractRequest("GetChallengedClaims")()
	result, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, z.contract.Call(methodChallengedClaims))
	if err != nil {
		return false, fmt.Errorf("failed to fetch challengedClaims: %w,targetIdx:%d", err, targetIdx)
	}
	return result.GetBool(0), nil
}

func (z *ZKFaultDisputeGameContract) GetParentGame(ctx context.Context) (ZKFaultDisputeGame, error) {
	defer z.metrics.StartContractRequest("GetParentGame")()
	result, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, z.contract.Call(methodParentGameProxy))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch parent game proxy: %w", err)
	}
	address := result.GetAddress(0)
	if address == (common.Address{}) {
		return nil, nil
	}
	contract, err := NewZKFaultDisputeGameContract(ctx, z.metrics, address, z.multiCaller)
	if err != nil {
		return nil, fmt.Errorf("failed to new parent game contract(%s): %w", address, err)
	}
	return contract, nil
}

func (z *ZKFaultDisputeGameContract) GetAddr() common.Address {
	return z.contract.Addr()
}

func (z *ZKFaultDisputeGameContract) ChallengeBySignalTx(
	ctx context.Context,
	challengeIdx int,
) (txmgr.TxCandidate, error) {
	call := z.contract.Call(methodChallengeBySignal, big.NewInt(int64(challengeIdx)))
	return z.txWithBond(ctx, call)
}

func (z *ZKFaultDisputeGameContract) txWithBond(
	ctx context.Context,
	call *batching.ContractCall,
) (txmgr.TxCandidate, error) {
	tx, err := call.ToTxCandidate()
	if err != nil {
		return txmgr.TxCandidate{}, fmt.Errorf("failed to create transaction: %w", err)
	}
	tx.Value, err = z.GetChallengerBond(ctx)
	if err != nil {
		return txmgr.TxCandidate{}, fmt.Errorf("failed to fetch required bond: %w", err)
	}
	return tx, nil
}

func (z *ZKFaultDisputeGameContract) GetChallengerBond(ctx context.Context) (*big.Int, error) {
	defer z.metrics.StartContractRequest("GetChallengerBond")()
	bond, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, z.contract.Call(methodChallengerBond))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve required challenger bond: %w", err)
	}
	return bond.GetBigInt(0), nil
}

func (z *ZKFaultDisputeGameContract) IsChallengeSuccess(ctx context.Context) (bool, error) {
	defer z.metrics.StartContractRequest("IsChallengeSuccess")()
	result, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, z.contract.Call(methodIsChallengeSuccess))
	if err != nil {
		return false, fmt.Errorf("failed to retrieve isChallengeSuccess: %w", err)
	}
	return result.GetBool(0), nil
}

func (z *ZKFaultDisputeGameContract) ResolveClaimTx() (txmgr.TxCandidate, error) {
	call := z.resolveClaimCall()
	return call.ToTxCandidate()
}

func (z *ZKFaultDisputeGameContract) CallResolveClaim(ctx context.Context) error {
	defer z.metrics.StartContractRequest("CallResolveClaim")()
	call := z.resolveClaimCall()
	_, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, call)
	if err != nil {
		return fmt.Errorf("failed to call resolve claim: %w", err)
	}
	return nil
}

func (z *ZKFaultDisputeGameContract) resolveClaimCall() *batching.ContractCall {
	return z.contract.Call(methodResolveClaim)
}

func (z *ZKFaultDisputeGameContract) CallResolve(ctx context.Context) (gameTypes.GameStatus, error) {
	defer z.metrics.StartContractRequest("CallResolve")()
	call := z.resolveCall()
	result, err := z.multiCaller.SingleCall(ctx, rpcblock.Latest, call)
	if err != nil {
		return gameTypes.GameStatusInProgress, fmt.Errorf("failed to call resolve claim: %w", err)
	}
	return gameTypes.GameStatusFromUint8(result.GetUint8(0))
}

func (z *ZKFaultDisputeGameContract) ResolveTx() (txmgr.TxCandidate, error) {
	call := z.resolveCall()
	return call.ToTxCandidate()
}

func (z *ZKFaultDisputeGameContract) resolveCall() *batching.ContractCall {
	return z.contract.Call(methodResolve)
}

type ZKFaultDisputeGame interface {
	GetStatus(ctx context.Context) (gameTypes.GameStatus, error)
	GetClaimCount(context.Context) (uint64, error)
	GetL1Head(ctx context.Context) (common.Hash, error)
	GetBlockRange(ctx context.Context) (prestateBlock uint64, poststateBlock uint64, retErr error)
	GetMaxDetectFaultDuration(ctx context.Context) (time.Duration, error)
	GetCreatedAt(ctx context.Context) (time.Time, error)
	GetClaimsHash(ctx context.Context) (common.Hash, error)
	GetRootClaim(ctx context.Context) (common.Hash, error)
	GetChallengedClaims(ctx context.Context, targetIdx int) (bool, error)
	GetParentGame(ctx context.Context) (ZKFaultDisputeGame, error)
	GetAddr() common.Address
	ChallengeBySignalTx(ctx context.Context, challengeIdx int) (txmgr.TxCandidate, error)
	IsChallengeSuccess(ctx context.Context) (bool, error)
	CallResolveClaim(ctx context.Context) error
	ResolveClaimTx() (txmgr.TxCandidate, error)
	CallResolve(ctx context.Context) (gameTypes.GameStatus, error)
	ResolveTx() (txmgr.TxCandidate, error)
}
