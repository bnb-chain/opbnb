package outputs

import (
	"context"
	"time"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/contracts"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/log"
	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/exp/slices"
)

type OutputCacheLoader struct {
	inner              *lru.Cache[uint64, *eth.OutputResponse]
	rollupClient       fault.RollupClient
	zkFaultDisputeGame contracts.ZKFaultDisputeGame
	logger             log.Logger
	ctx                context.Context
}

func (l *OutputCacheLoader) Load(startBlock uint64, endBlock uint64) []*eth.OutputResponse {
	currentBlock := startBlock
	var needFetchBlock []uint64
	var result []*eth.OutputResponse
	for currentBlock <= endBlock {
		if outputInCache, ok := l.inner.Peek(currentBlock); !ok {
			needFetchBlock = append(needFetchBlock, currentBlock)
		} else {
			result = append(result, outputInCache)
		}
		//todo_welkin add flag for 3
		currentBlock = currentBlock + 3
	}

	if len(needFetchBlock) > 0 {
		for {
			outputResponses, err := l.rollupClient.BatchOutputAtBlock(l.ctx, needFetchBlock)
			if err != nil {
				l.logger.Warn("failed to load output,will retry", "block", currentBlock, "err", err)
				time.Sleep(1 * time.Second)
				continue
			}
			for _, outputResp := range outputResponses {
				l.inner.Add(outputResp.BlockRef.Number, outputResp)
				result = append(result, outputResp)
			}
			break
		}
	}
	if len(result) != int(endBlock-startBlock)/3 {
		l.logger.Warn("maybe miss outputRoot in loader", "len", len(result), "endBlock", endBlock, "startBlock", startBlock)
	}
	slices.SortFunc(result, func(a, b *eth.OutputResponse) int {
		return int(a.BlockRef.Number - b.BlockRef.Number)
	})
	return result
}

func (l *OutputCacheLoader) LoadOne(block uint64) (*eth.OutputResponse, error) {
	if outputResp, ok := l.inner.Peek(block); ok {
		return outputResp, nil
	}
	outputResponse, err := l.rollupClient.OutputAtBlock(l.ctx, block)
	if err != nil {
		return nil, err
	}
	l.inner.Add(outputResponse.BlockRef.Number, outputResponse)
	return outputResponse, nil
}

func NewOutputCacheLoader(
	ctx context.Context,
	rollupClient fault.RollupClient,
	logger log.Logger,
	zkFaultDisputeGame contracts.ZKFaultDisputeGame,
) *OutputCacheLoader {
	lruCache, _ := lru.New[uint64, *eth.OutputResponse](5000)
	return &OutputCacheLoader{
		ctx:                ctx,
		rollupClient:       rollupClient,
		logger:             logger,
		zkFaultDisputeGame: zkFaultDisputeGame,
		inner:              lruCache,
	}
}
