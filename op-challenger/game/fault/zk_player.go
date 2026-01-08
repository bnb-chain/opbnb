package fault

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/contracts"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/trace/outputs"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/trace/zk"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
	gameTypes "github.com/ethereum-optimism/optimism/op-challenger/game/types"
	"github.com/ethereum-optimism/optimism/op-challenger/metrics"
	"github.com/ethereum-optimism/optimism/op-service/clock"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type ZKGamePlayer struct {
	act               actor
	status            gameTypes.GameStatus
	loader            GameInfo
	logger            log.Logger
	syncValidator     SyncValidator
	gameL1Head        eth.BlockID
	outputCacheLoader *outputs.OutputCacheLoader
}

func NewZKGamePlayer(
	ctx context.Context,
	clock clock.Clock,
	l1Clock types.ClockReader,
	logger log.Logger,
	m metrics.Metricer,
	addr common.Address,
	loader contracts.ZKFaultDisputeGame,
	l1HeaderSource L1HeaderSource,
	syncValidator SyncValidator,
	outputCacheLoader *outputs.OutputCacheLoader,
	factory contracts.DisputeGameFactory,
	sender TxSender,
	index uint64,
	challengeByProof bool,
	dir string,
	responseChallengeByProof bool,
	responseClaimants []common.Address,
) (*ZKGamePlayer, error) {
	logger = logger.New("zkgame", addr, "idx", index)

	status, err := loader.GetStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load zk game status: %w", err)
	}
	if status != gameTypes.GameStatusInProgress {
		logger.Info("Game already resolved", "status", status)
		// Game is already complete so skip creating the trace provider, loading game inputs etc.
		return &ZKGamePlayer{
			logger: logger,
			loader: loader,
			status: status,
			// Act function does nothing because the game is already complete
			act: func(ctx context.Context) error {
				return nil
			},
		}, nil
	}
	logger.Debug("Game will play")
	l1HeadHash, err := loader.GetL1Head(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load game L1 head: %w", err)
	}
	l1Header, err := l1HeaderSource.HeaderByHash(ctx, l1HeadHash)
	if err != nil {
		return nil, fmt.Errorf("failed to load L1 header %v: %w", l1HeadHash, err)
	}
	l1Head := eth.HeaderBlockID(l1Header)
	detectFaultDuration, err := loader.GetMaxDetectFaultDuration(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load max detect fault duration: %w", err)
	}
	maxClockDuration, err := loader.GetMaxClockDuration(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail get max clock duration:%w", err)
	}
	maxGenerateProofDuration, err := loader.GetMaxGenerateProofDuration(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail get max generate proof duration:%w", err)
	}
	createAt, err := loader.GetCreatedAt(ctx)
	if err != nil {
		return nil, err
	}
	config, err := loader.GetConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	blockDistance, err := config.GetBlockDistance(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load block distance: %w", err)
	}

	startBlock, endBlock, err := loader.GetBlockRange(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load game block range: %w", err)
	}

	proofAccessor, err := zk.NewProofAccessor(dir, l1Head)
	if err != nil {
		return nil, fmt.Errorf("fail new proofAccessor: %w", err)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("could not create proofs directory %v: %w", dir, err)
	}

	gameFactory, ok := factory.(*contracts.ZkGameFactory)
	if !ok {
		return nil, fmt.Errorf("factory should be zkGameFactory,gameAddr:%s", addr)
	}
	agent := NewZkAgent(m, clock, l1Clock, loader, outputCacheLoader, startBlock, endBlock, detectFaultDuration, maxGenerateProofDuration, maxClockDuration,
		createAt, logger, l1Head, l1HeaderSource, addr, gameFactory, sender, challengeByProof, proofAccessor, blockDistance, responseChallengeByProof, responseClaimants)

	return &ZKGamePlayer{
		logger:            logger,
		loader:            loader,
		status:            status,
		act:               agent.Act,
		syncValidator:     syncValidator,
		gameL1Head:        l1Head,
		outputCacheLoader: outputCacheLoader,
	}, nil
}

func (z *ZKGamePlayer) ValidatePrestate(ctx context.Context) error {
	return nil
}

func (z *ZKGamePlayer) ProgressGame(ctx context.Context) gameTypes.GameStatus {
	z.logger.Debug("progress game", "status", z.status)
	if z.status != gameTypes.GameStatusInProgress {
		// Game is already complete so don't try to perform further actions.
		z.logger.Trace("Skipping completed game")
		return z.status
	}
	if err := z.syncValidator.ValidateNodeSynced(ctx, z.gameL1Head); errors.Is(err, ErrNotInSync) {
		z.logger.Warn("Local node not sufficiently up to date", "err", err)
		return z.status
	} else if err != nil {
		z.logger.Error("Could not check local node was in sync", "err", err)
		return z.status
	}
	z.logger.Trace("Checking if actions are required")
	if err := z.act(ctx); err != nil {
		z.logger.Error("Error when acting on game", "err", err)
	}
	status, err := z.loader.GetStatus(ctx)
	if err != nil {
		z.logger.Error("Unable to retrieve game status", "err", err)
		return gameTypes.GameStatusInProgress
	}
	z.logGameStatus(ctx, status)
	z.status = status
	return status
}

func (z *ZKGamePlayer) logGameStatus(ctx context.Context, status gameTypes.GameStatus) {
	if status == gameTypes.GameStatusInProgress {
		claimCount, err := z.loader.GetClaimCount(ctx)
		if err != nil {
			z.logger.Error("Failed to get claim count for in progress game", "err", err)
			return
		}
		z.logger.Info("Game info", "claims", claimCount, "status", status)
		return
	}
	z.logger.Info("Game resolved", "status", status)
}

func (z *ZKGamePlayer) Status() gameTypes.GameStatus {
	return z.status
}
