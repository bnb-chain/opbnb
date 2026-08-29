package proposer

import (
	"context"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/log"
)

type OutputRootCacheHandler struct {
	ctx                 context.Context
	log                 log.Logger
	isStart             *atomic.Bool
	readyChan           chan *outputRootBatchData
	outputRootFetchFunc func(ctx context.Context, block []*big.Int) ([]*eth.OutputResponse, bool, error)
	batchSize           uint64
}

type outputRootBatchData struct {
	outputRootList  []eth.Bytes32
	lastSyncStatus  *eth.SyncStatus
	lastL2BlockRef  *eth.L2BlockRef
	parentGameIndex *big.Int
	l2BlockNumber   *big.Int
}

func newOutputRootCacheHandler(
	ctx context.Context,
	log log.Logger,
	outputRootFetchFunc func(ctx context.Context, block []*big.Int) ([]*eth.OutputResponse, bool, error),
	batchSize uint64,
) *OutputRootCacheHandler {
	return &OutputRootCacheHandler{
		ctx:                 ctx,
		log:                 log,
		isStart:             &atomic.Bool{},
		readyChan:           make(chan *outputRootBatchData, 1),
		outputRootFetchFunc: outputRootFetchFunc,
		batchSize:           batchSize,
	}
}

func (h *OutputRootCacheHandler) startFrom(parentGame *GameInformation, distance *big.Int) {
	if !h.isStart.CompareAndSwap(false, true) {
		return
	}
	go h.loop(parentGame, distance)
}

func (h *OutputRootCacheHandler) loop(parentGame *GameInformation, distance *big.Int) {
	currentBlockNumber := new(big.Int).Add(parentGame.extraData.endL2BlockNumber, distance)
	endBlockNumber := new(big.Int).Add(parentGame.extraData.endL2BlockNumber, new(big.Int).SetUint64(h.batchSize))
	outputRootList := make([]eth.Bytes32, 0, h.batchSize/distance.Uint64())
	h.log.Debug("outputRoot loop", "parentGame index", parentGame.game.Index.Uint64(),
		"currentBlockNumber", currentBlockNumber, "endBlockNumber", endBlockNumber, "stepBigInt", distance)
	var blockList []*big.Int
	var lastSyncStatus *eth.SyncStatus
	var lastBlockRef *eth.L2BlockRef
	for currentBlockNumber.Cmp(endBlockNumber) <= 0 {
		blockList = append(blockList, currentBlockNumber)
		currentBlockNumber = new(big.Int).Add(currentBlockNumber, distance)
	}

	for {
		outputRootResponses, shouldPropose, err := h.outputRootFetchFunc(h.ctx, blockList)
		if err != nil {
			h.log.Error("failed to fetch outputRoot", "err", err, "blockList", blockList)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if !shouldPropose {
			h.log.Warn("the block number cannot submit the output root yet; it needs to wait for a while")
			time.Sleep(1 * time.Second)
			continue
		}
		for _, outputRootResp := range outputRootResponses {
			outputRootList = append(outputRootList, outputRootResp.OutputRoot)
			lastSyncStatus = outputRootResp.Status
			lastBlockRef = &outputRootResp.BlockRef
		}
		break
	}
	h.readyChan <- &outputRootBatchData{
		outputRootList:  outputRootList,
		lastSyncStatus:  lastSyncStatus,
		lastL2BlockRef:  lastBlockRef,
		parentGameIndex: parentGame.game.Index,
		l2BlockNumber:   endBlockNumber,
	}
	h.isStart.CompareAndSwap(true, false)
	h.log.Debug("outputRoot loop end", "currentBlockNumber", currentBlockNumber, "outputRootListSize", len(outputRootList))
}
