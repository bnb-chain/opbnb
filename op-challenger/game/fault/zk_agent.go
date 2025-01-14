package fault

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/contracts"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/trace/outputs"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
	"github.com/ethereum-optimism/optimism/op-challenger/metrics"
	"github.com/ethereum-optimism/optimism/op-service/clock"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
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
	}
}

func (z *ZkAgent) Act(ctx context.Context) error {
	if z.tryResolve(ctx) {
		return nil
	}

	start := z.systemClock.Now()
	defer func() {
		z.metrics.RecordGameActTime(z.systemClock.Since(start).Seconds())
	}()

	z.shouldChallenge(ctx)
	return nil
}

func (z *ZkAgent) tryResolve(ctx context.Context) bool {
	return true
}

func (z *ZkAgent) shouldChallenge(ctx context.Context) (bool, error) {
	var targetIdx int
	claimsHashNotEqual := false
	rootClaim, err := z.zkFaultDisputeGame.GetRootClaim(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get root claim: %w", err)
	}
	endBlockOutputRootResp, err := z.outputCacheLoader.LoadOne(z.endBlock)
	if err != nil {
		return false, fmt.Errorf("failed to load end block outputRoot resp: %w", err)
	}
	if rootClaim != common.Hash(endBlockOutputRootResp.OutputRoot) {
		targetIdx = int(z.endBlock-z.startBlock) / 3
		claimsHashNotEqual = true
	} else {
		claimsHash, err := z.zkFaultDisputeGame.GetClaimsHash(ctx)
		if err != nil {
			return false, fmt.Errorf("failed to get claims hash: %w", err)
		}
		outputRootResps := z.outputCacheLoader.Load(z.startBlock, z.endBlock)
		var outputRoots []eth.Bytes32
		for _, outputRootResp := range outputRootResps {
			outputRoots = append(outputRoots, outputRootResp.OutputRoot)
		}
		realOutputRootHash := calHash(outputRoots)
		claimsHashNotEqual = realOutputRootHash != claimsHash
		if claimsHashNotEqual {
			idx, originClaims, err := z.findTargetChallengeIdx(outputRootResps)
			if err != nil {
				return false, err
			}
			targetIdx = idx
			z.originClaims = originClaims
		}
	}

	claims, err := z.zkFaultDisputeGame.GetChallengedClaims(ctx, targetIdx)
	if err != nil {
		return false, fmt.Errorf("failed to get challenged claims: %w", err)
	}
	if claims {
		z.log.Trace("Challenged claims", "targetIdx", targetIdx)
		return false, nil
	}

	detectFaultDeadline := z.createAt.Add(z.maxDetectFaultDuration)
	if time.Now().After(detectFaultDeadline) {
		if claimsHashNotEqual {
			z.log.Error("discovered a game that should be challenged, but the challenge period has passed!")
		}
		return false, nil
	}
	z.targetChallengeIdx = targetIdx
	return true, nil
}

func (z *ZkAgent) findTargetChallengeIdx(outputResponses []*eth.OutputResponse) (int, []eth.Bytes32, error) {
	return 0, nil, nil
}

func calHash(cache []eth.Bytes32) common.Hash {
	var encodeBytes [][]byte
	for _, one := range cache {
		encodeBytes = append(encodeBytes, one[:])
	}
	return crypto.Keccak256Hash(encodeBytes...)
}
