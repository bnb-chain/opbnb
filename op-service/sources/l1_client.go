package sources

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/sources/caching"
)

const sequencerConfDepth = 15

type L1ClientConfig struct {
	EthClientConfig

	L1BlockRefsCacheSize int
}

func L1ClientDefaultConfig(config *rollup.Config, trustRPC bool, kind RPCProviderKind) *L1ClientConfig {
	// Cache 3/2 worth of sequencing window of receipts and txs
	span := int(config.SeqWindowSize) * 3 / 2
	fullSpan := span
	if span > 1000 { // sanity cap. If a large sequencing window is configured, do not make the cache too large
		span = 1000
	}
	return &L1ClientConfig{
		EthClientConfig: EthClientConfig{
			// receipts and transactions are cached per block
			ReceiptsCacheSize:     span,
			TransactionsCacheSize: span,
			HeadersCacheSize:      span,
			PayloadsCacheSize:     span,
			MaxRequestsPerBatch:   20, // TODO: tune batch param
			MaxConcurrentRequests: 20,
			TrustRPC:              trustRPC,
			MustBePostMerge:       false,
			RPCProviderKind:       kind,
			MethodResetDuration:   time.Minute,
		},
		// Not bounded by span, to cover find-sync-start range fully for speedy recovery after errors.
		L1BlockRefsCacheSize: fullSpan,
	}
}

// L1Client provides typed bindings to retrieve L1 data from an RPC source,
// with optimized batch requests, cached results, and flag to not trust the RPC
// (i.e. to verify all returned contents against corresponding block hashes).
type L1Client struct {
	*EthClient

	// cache L1BlockRef by hash
	// common.Hash -> eth.L1BlockRef
	l1BlockRefsCache *caching.LRUCache[common.Hash, eth.L1BlockRef]

	//ensure pre-fetch receipts only once
	preFetchReceiptsOnce sync.Once
	//start block for pre-fetch receipts
	preFetchReceiptsStartBlockChan chan uint64
	//max concurrent requests
	maxConcurrentRequests int
	//done chan
	done chan struct{}
}

// NewL1Client wraps a RPC with bindings to fetch L1 data, while logging errors, tracking metrics (optional), and caching.
func NewL1Client(client client.RPC, log log.Logger, metrics caching.Metrics, config *L1ClientConfig) (*L1Client, error) {
	ethClient, err := NewEthClient(client, log, metrics, &config.EthClientConfig, true)
	if err != nil {
		return nil, err
	}

	return &L1Client{
		EthClient:                      ethClient,
		l1BlockRefsCache:               caching.NewLRUCache[common.Hash, eth.L1BlockRef](metrics, "blockrefs", config.L1BlockRefsCacheSize),
		preFetchReceiptsOnce:           sync.Once{},
		preFetchReceiptsStartBlockChan: make(chan uint64, 1),
		maxConcurrentRequests:          config.MaxConcurrentRequests,
		done:                           make(chan struct{}),
	}, nil
}

// L1BlockRefByLabel returns the [eth.L1BlockRef] for the given block label.
// Notice, we cannot cache a block reference by label because labels are not guaranteed to be unique.
func (s *L1Client) L1BlockRefByLabel(ctx context.Context, label eth.BlockLabel) (eth.L1BlockRef, error) {
	info, err := s.BSCInfoByLabel(ctx, label)
	if label == eth.Finalized && err != nil && strings.Contains(err.Error(), "eth_getFinalizedHeader does not exist") {
		// op-e2e not support bsc as L1 currently, so fallback to not use bsc specific method eth_getFinalizedBlock
		info, err = s.InfoByLabel(ctx, label)
	}
	if err != nil {
		// Both geth and erigon like to serve non-standard errors for the safe and finalized heads, correct that.
		// This happens when the chain just started and nothing is marked as safe/finalized yet.
		if strings.Contains(err.Error(), "block not found") || strings.Contains(err.Error(), "Unknown block") {
			err = ethereum.NotFound
		}
		return eth.L1BlockRef{}, fmt.Errorf("failed to fetch head header: %w", err)
	}
	ref := eth.InfoToL1BlockRef(info)
	s.l1BlockRefsCache.Add(ref.Hash, ref)
	return ref, nil
}

// L1BlockRefByNumber returns an [eth.L1BlockRef] for the given block number.
// Notice, we cannot cache a block reference by number because L1 re-orgs can invalidate the cached block reference.
func (s *L1Client) L1BlockRefByNumber(ctx context.Context, num uint64) (eth.L1BlockRef, error) {
	info, err := s.InfoByNumber(ctx, num)
	if err != nil {
		return eth.L1BlockRef{}, fmt.Errorf("failed to fetch header by num %d: %w", num, err)
	}
	ref := eth.InfoToL1BlockRef(info)
	s.l1BlockRefsCache.Add(ref.Hash, ref)
	return ref, nil
}

// L1BlockRefByHash returns the [eth.L1BlockRef] for the given block hash.
// We cache the block reference by hash as it is safe to assume collision will not occur.
func (s *L1Client) L1BlockRefByHash(ctx context.Context, hash common.Hash) (eth.L1BlockRef, error) {
	if v, ok := s.l1BlockRefsCache.Get(hash); ok {
		return v, nil
	}
	info, err := s.InfoByHash(ctx, hash)
	if err != nil {
		return eth.L1BlockRef{}, fmt.Errorf("failed to fetch header by hash %v: %w", hash, err)
	}
	ref := eth.InfoToL1BlockRef(info)
	s.l1BlockRefsCache.Add(ref.Hash, ref)
	return ref, nil
}

func (s *L1Client) GoOrUpdatePreFetchReceipts(ctx context.Context, l1Start uint64) error {
	s.preFetchReceiptsStartBlockChan <- l1Start
	s.preFetchReceiptsOnce.Do(func() {
		s.log.Info("pre-fetching receipts start", "startBlock", l1Start)
		go func() {
			var currentL1Block uint64
			var parentHash common.Hash
			for {
				select {
				case <-s.done:
					s.log.Info("pre-fetching receipts done")
					return
				case currentL1Block = <-s.preFetchReceiptsStartBlockChan:
					s.log.Debug("pre-fetching receipts currentL1Block changed", "block", currentL1Block)
					s.recProvider.GetReceiptsCache().RemoveAll()
					parentHash = common.Hash{}
				default:
					blockRef, err := s.L1BlockRefByLabel(ctx, eth.Unsafe)
					if err != nil {
						s.log.Debug("failed to fetch latest block ref", "err", err)
						time.Sleep(3 * time.Second)
						continue
					}

					if currentL1Block > blockRef.Number {
						s.log.Debug("current block height exceeds the latest block height of l1, will wait for a while.", "currentL1Block", currentL1Block, "l1Latest", blockRef.Number)
						time.Sleep(3 * time.Second)
						continue
					}

					var taskCount int
					maxConcurrent := s.maxConcurrentRequests / 2
					if blockRef.Number-currentL1Block >= uint64(maxConcurrent) {
						taskCount = maxConcurrent
					} else {
						taskCount = int(blockRef.Number-currentL1Block) + 1
					}

					blockInfoChan := make(chan eth.L1BlockRef, taskCount)
					oldestFetchBlockNumber := currentL1Block

					var wg sync.WaitGroup
					for i := 0; i < taskCount; i++ {
						wg.Add(1)
						go func(ctx context.Context, blockNumber uint64) {
							defer wg.Done()
							for {
								select {
								case <-s.done:
									return
								default:
									blockInfo, err := s.L1BlockRefByNumber(ctx, blockNumber)
									if err != nil {
										s.log.Debug("failed to fetch block ref", "err", err, "blockNumber", blockNumber)
										time.Sleep(1 * time.Second)
										continue
									}
									pair, ok := s.recProvider.GetReceiptsCache().Get(blockNumber, false)
									if ok && pair.blockHash == blockInfo.Hash {
										blockInfoChan <- blockInfo
										return
									}

									isSuccess, err := s.PreFetchReceipts(ctx, blockInfo.Hash)
									if err != nil {
										s.log.Warn("failed to pre-fetch receipts", "err", err)
										time.Sleep(1 * time.Second)
										continue
									}
									if !isSuccess {
										s.log.Debug("pre fetch receipts fail without error,need retry", "blockHash", blockInfo.Hash, "blockNumber", blockNumber)
										time.Sleep(1 * time.Second)
										continue
									}
									s.log.Debug("pre-fetching receipts done", "block", blockInfo.Number, "hash", blockInfo.Hash)
									blockInfoChan <- blockInfo
									return
								}
							}
						}(ctx, currentL1Block)
						currentL1Block = currentL1Block + 1
					}
					wg.Wait()
					close(blockInfoChan)

					//try to find out l1 reOrg and return to an earlier block height for re-prefetching
					var latestBlockHash common.Hash
					latestBlockNumber := uint64(0)
					var oldestBlockParentHash common.Hash
					for l1BlockInfo := range blockInfoChan {
						if l1BlockInfo.Number > latestBlockNumber {
							latestBlockHash = l1BlockInfo.Hash
							latestBlockNumber = l1BlockInfo.Number
						}
						if l1BlockInfo.Number == oldestFetchBlockNumber {
							oldestBlockParentHash = l1BlockInfo.ParentHash
						}
					}

					s.log.Debug("pre-fetching receipts hash", "latestBlockHash", latestBlockHash, "latestBlockNumber", latestBlockNumber, "oldestBlockNumber", oldestFetchBlockNumber, "oldestBlockParentHash", oldestBlockParentHash)
					if parentHash != (common.Hash{}) && oldestBlockParentHash != (common.Hash{}) && oldestBlockParentHash != parentHash && currentL1Block >= sequencerConfDepth+uint64(taskCount) {
						currentL1Block = currentL1Block - sequencerConfDepth - uint64(taskCount)
						s.log.Warn("pre-fetching receipts found l1 reOrg, return to an earlier block height for re-prefetching", "recordParentHash", parentHash, "unsafeParentHash", oldestBlockParentHash, "number", oldestFetchBlockNumber, "backToNumber", currentL1Block)
						parentHash = common.Hash{}
						continue
					}
					parentHash = latestBlockHash
				}
			}
		}()
	})
	return nil
}

func (s *L1Client) ClearReceiptsCacheBefore(blockNumber uint64) {
	s.log.Debug("clear receipts cache before", "blockNumber", blockNumber)
	s.recProvider.GetReceiptsCache().RemoveLessThan(blockNumber)
}

func (s *L1Client) Close() {
	close(s.done)
	s.EthClient.Close()
}
