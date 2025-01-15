package fault

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/contracts"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/trace/outputs"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
	gameTypes "github.com/ethereum-optimism/optimism/op-challenger/game/types"
	"github.com/ethereum-optimism/optimism/op-challenger/metrics"
	"github.com/ethereum-optimism/optimism/op-service/clock"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

type LastAgentState int

const (
	NotHandledYet LastAgentState = iota
	SkipChallenge
	Challenged
	Resolved
)

type ZkAgent struct {
	metrics                metrics.Metricer
	systemClock            clock.Clock
	l1Clock                types.ClockReader
	responder              Responder
	selective              bool
	maxClockDuration       time.Duration
	log                    log.Logger
	zkFaultDisputeGame     contracts.ZKFaultDisputeGame
	outputCacheLoader      *outputs.OutputCacheLoader
	startBlock             uint64
	endBlock               uint64
	maxDetectFaultDuration time.Duration
	createAt               time.Time
	originClaims           []eth.Bytes32
	targetChallengeIdx     int
	l1Head                 eth.BlockID
	l1Source               L1HeaderSource
	gameAddr               common.Address
	factory                *contracts.ZkGameFactory
	lastState              LastAgentState
	txSender               TxSender
}

func NewZkAgent(
	m metrics.Metricer,
	systemClock clock.Clock,
	l1Clock types.ClockReader,
	loader contracts.ZKFaultDisputeGame,
	cacheLoader *outputs.OutputCacheLoader,
	startBlock uint64,
	endBlock uint64,
	maxDetectFaultDuration time.Duration,
	createAt time.Time,
	logger log.Logger,
	l1Head eth.BlockID,
	l1Source L1HeaderSource,
	addr common.Address,
	factory *contracts.ZkGameFactory,
	sender TxSender,
) *ZkAgent {
	return &ZkAgent{
		metrics:                m,
		systemClock:            systemClock,
		l1Clock:                l1Clock,
		zkFaultDisputeGame:     loader,
		outputCacheLoader:      cacheLoader,
		startBlock:             startBlock,
		endBlock:               endBlock,
		maxDetectFaultDuration: maxDetectFaultDuration,
		createAt:               createAt,
		log:                    logger,
		l1Head:                 l1Head,
		l1Source:               l1Source,
		gameAddr:               addr,
		factory:                factory,
		txSender:               sender,
		lastState:              NotHandledYet,
	}
}

func (z *ZkAgent) Act(ctx context.Context) error {
	start := z.systemClock.Now()
	defer func() {
		z.metrics.RecordGameActTime(z.systemClock.Since(start).Seconds())
	}()

	if z.lastState == Resolved {
		z.log.Debug("already resolved,skip")
		return nil
	}

	if z.tryResolve(ctx) {
		z.lastState = Resolved
		z.log.Debug("resolve success")
		return nil
	}

	if z.lastState == SkipChallenge || z.lastState == Challenged {
		z.log.Debug("skip challenge", "lastState", z.lastState)
		return nil
	}
	should, err := z.shouldChallenge(ctx, z.zkFaultDisputeGame, z.gameAddr, z.startBlock, z.endBlock, z.l1Head.Number)
	if err != nil {
		return err
	}
	if should {
		z.log.Debug("found game should be challenged")
		valid, err := z.isParentGamesValid(ctx)
		if err != nil {
			return err
		}
		if !valid {
			z.log.Debug("the parent game is invalid,so we skip the challenge", "gameAddr", z.gameAddr)
			z.lastState = SkipChallenge
			return nil
		}
		detectFaultDeadline := z.createAt.Add(z.maxDetectFaultDuration)
		if time.Now().After(detectFaultDeadline) {
			z.log.Error("discovered a game that should be challenged, but the challenge period has passed!")
			z.lastState = SkipChallenge
			return fmt.Errorf("discovered a game that should be challenged, but the challenge period has passed! deadline:%s", detectFaultDeadline)
		}
		err = z.submitChallengeBySignal(ctx)
		if err != nil {
			return err
		}
		z.lastState = Challenged
		z.log.Debug("challenge success")
	} else {
		z.lastState = SkipChallenge
		z.log.Debug("game is ok, skip challenge")
	}

	return nil
}

func (z *ZkAgent) tryResolve(ctx context.Context) bool {
	status, err := z.zkFaultDisputeGame.GetStatus(ctx)
	if err != nil {
		z.log.Error("fail get game status when trying resolve", "err", err)
		return false
	}
	if status != gameTypes.GameStatusInProgress {
		z.log.Debug("the game is not in progress when trying resolve,skip", "status", status)
		return false
	}
	isChallengeSuccess, err := z.zkFaultDisputeGame.IsChallengeSuccess(ctx)
	if err != nil {
		z.log.Error("fail get game IsChallengeSuccess when trying resolve", "err", err)
		return false
	}
	if !isChallengeSuccess {
		err := z.zkFaultDisputeGame.CallResolveClaim(ctx)
		if err != nil {
			z.log.Debug("fail call resolveClaim when trying resolve", "err", err)
		} else {
			candidate, err := z.zkFaultDisputeGame.ResolveClaimTx()
			if err != nil {
				z.log.Error("fail build resolveClaimTx when trying resolve", "err", err)
				return false
			}
			err = z.txSender.SendAndWaitSimple("resolveClaim", candidate)
			if err != nil {
				z.log.Error("fail send resolveClaim when trying resolveClaim", "err", err)
				return false
			}
		}
	}
	gameStatus, err := z.zkFaultDisputeGame.CallResolve(ctx)
	if err != nil {
		z.log.Debug("fail to call resolve, maybe not ready yet", "err", err)
		return false
	}
	z.log.Debug("will submit resolve tx", "status", gameStatus)
	resolveTx, err := z.zkFaultDisputeGame.ResolveTx()
	if err != nil {
		z.log.Error("fail build resolveTx when trying resolve", "err", err)
		return false
	}
	err = z.txSender.SendAndWaitSimple("resolve", resolveTx)
	if err != nil {
		z.log.Error("fail send resolveTx when trying resolve", "err", err)
		return false
	}
	return true
}

func (z *ZkAgent) submitChallengeBySignal(ctx context.Context) error {
	tx, err := z.zkFaultDisputeGame.ChallengeBySignalTx(ctx, z.targetChallengeIdx)
	if err != nil {
		return fmt.Errorf("build challenge by signal tx failed: %w", err)
	}
	err = z.txSender.SendAndWaitSimple("challengeBySignal", tx)
	if err != nil {
		return fmt.Errorf("submit challenge by signal tx failed: %w", err)
	}
	return nil
}

func (z *ZkAgent) shouldChallenge(
	ctx context.Context,
	game contracts.ZKFaultDisputeGame,
	addr common.Address,
	startBlock uint64,
	endBlock uint64,
	l1ParentBlockNumber uint64,
) (bool, error) {
	isChallengeSuccess, err := game.IsChallengeSuccess(ctx)
	if err != nil {
		return false, fmt.Errorf("check if challenge success failed: %w", err)
	}

	if isChallengeSuccess {
		return false, nil
	}

	var targetIdx int
	claimsHashNotEqual := false
	rootClaim, err := game.GetRootClaim(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get root claim: %w", err)
	}
	endBlockOutputRootResp, err := z.outputCacheLoader.LoadOne(endBlock)
	if err != nil {
		return false, fmt.Errorf("failed to load end block outputRoot resp: %w", err)
	}
	if rootClaim != common.Hash(endBlockOutputRootResp.OutputRoot) {
		targetIdx = int(endBlock-startBlock)/3 - 1
		claimsHashNotEqual = true
	} else {
		claimsHash, err := game.GetClaimsHash(ctx)
		if err != nil {
			return false, fmt.Errorf("failed to get claims hash: %w", err)
		}
		outputRootResps := z.outputCacheLoader.Load(startBlock, endBlock)
		var outputRoots []eth.Bytes32
		for _, outputRootResp := range outputRootResps {
			outputRoots = append(outputRoots, outputRootResp.OutputRoot)
		}
		realOutputRootHash := calHash(outputRoots)
		claimsHashNotEqual = realOutputRootHash != claimsHash
		if claimsHashNotEqual {
			idx, originClaims, err := z.findTargetChallengeIdx(ctx, outputRootResps, addr, l1ParentBlockNumber)
			if err != nil {
				return false, err
			}
			targetIdx = idx
			z.originClaims = originClaims
		}
	}

	if !claimsHashNotEqual {
		return false, nil
	}

	claims, err := game.GetChallengedClaims(ctx, targetIdx)
	if err != nil {
		return false, fmt.Errorf("failed to get challenged claims: %w", err)
	}
	if claims {
		z.log.Debug("Claims have already been challenged.", "targetIdx", targetIdx)
		return false, nil
	}

	z.targetChallengeIdx = targetIdx
	return true, nil
}

func (z *ZkAgent) findTargetChallengeIdx(
	ctx context.Context,
	outputResponses []*eth.OutputResponse,
	addr common.Address,
	l1ParentBlockNumber uint64,
) (int, []eth.Bytes32, error) {
	createCallData, err := z.getGameCreateCallData(ctx, l1ParentBlockNumber, addr)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get create call data: %w", err)
	}
	if len(outputResponses) != len(createCallData.Claims) {
		return 0, nil, fmt.Errorf("the size of outputResponses is inconsistent with the claims in calldata, outputResp size: %d, calldata claims size: %d", len(outputResponses), len(createCallData.Claims))
	}
	for idx, claim := range createCallData.Claims {
		if outputResponses[idx].OutputRoot != claim {
			return idx, createCallData.Claims, nil
		}
	}

	return 0, nil, errors.New("no challenged outputs found")
}

func (z *ZkAgent) getGameCreateCallData(
	ctx context.Context,
	l1ParentBlockNumber uint64,
	addr common.Address,
) (*contracts.ZkGameCreateCallData, error) {
	receipts, err := z.l1Source.BlockReceipts(ctx, rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(l1ParentBlockNumber+1)))
	if err != nil {
		return nil, fmt.Errorf("fail get l1 block receipts: %w,number:%d", err, l1ParentBlockNumber+1)
	}
	var gameCreateTxHash *common.Hash
	for _, receipt := range receipts {
		address, _, _, err := z.factory.DecodeDisputeGameCreatedLog(receipt)
		if err != nil && !errors.Is(err, contracts.ErrEventNotFound) {
			return nil, fmt.Errorf("fail decode dispute game created log: %w", err)
		}
		if errors.Is(err, contracts.ErrEventNotFound) {
			continue
		}
		if address == addr {
			gameCreateTxHash = &receipt.TxHash
		}
	}
	if gameCreateTxHash == nil {
		return nil, fmt.Errorf("we did not find the hash of the transaction that created the game.addr:%s,blockNumber:%d", addr, l1ParentBlockNumber)
	}
	transactionByHash, _, err := z.l1Source.TransactionByHash(ctx, *gameCreateTxHash)
	if err != nil {
		return nil, fmt.Errorf("fail get transaction by hash(%s): %w", *gameCreateTxHash, err)
	}
	callData := transactionByHash.Data()
	gameCreateCallData, err := z.factory.DecodeZKGameCreateCallData(callData)
	if err != nil {
		return nil, fmt.Errorf("fail decode game create call data: %w,callData:%s", err, hexutil.Encode(callData))
	}
	return gameCreateCallData, nil
}

func (z *ZkAgent) isParentGamesValid(ctx context.Context) (bool, error) {
	parent, err := z.zkFaultDisputeGame.GetParentGame(ctx)
	if err != nil {
		return false, fmt.Errorf("fail get parent games proxy for challenge: %w", err)
	}
	if parent == nil {
		return true, nil
	}
	for {
		addr := parent.GetAddr()
		status, err := parent.GetStatus(ctx)
		if err != nil {
			return false, fmt.Errorf("fail get parent(%s) status: %w", addr, err)
		}
		if status == gameTypes.GameStatusDefenderWon {
			return true, nil
		}
		if status == gameTypes.GameStatusChallengerWon {
			return false, nil
		}
		start, end, err := parent.GetBlockRange(ctx)
		if err != nil {
			return false, fmt.Errorf("fail get parent game(%s) range for challenge: %w", addr, err)
		}
		l1Hash, err := parent.GetL1Head(ctx)
		if err != nil {
			return false, fmt.Errorf("fail get parent game(%s) l1 head for challenge: %w", addr, err)
		}
		l1Header, err := z.l1Source.HeaderByHash(ctx, l1Hash)
		if err != nil {
			return false, fmt.Errorf("fail get l1 head by hash(%s) for parent game(%s) check: %w", l1Hash, addr, err)
		}
		should, err := z.shouldChallenge(ctx, parent, addr, start, end, l1Header.Number.Uint64())
		if err != nil {
			return false, fmt.Errorf("fail check if the parent game(%s) should be challenged: %w", addr, err)
		}
		if should {
			return false, nil
		}
		parentGame, err := parent.GetParentGame(ctx)
		if err != nil {
			return false, fmt.Errorf("fail get parent game(%s) parent: %w", addr, err)
		}
		if parentGame == nil {
			break
		}
		parent = parentGame
	}
	return true, nil
}

func calHash(cache []eth.Bytes32) common.Hash {
	var encodeBytes []byte
	for _, one := range cache {
		encodeBytes = append(encodeBytes, one[:]...)
	}
	return crypto.Keccak256Hash(encodeBytes)
}
