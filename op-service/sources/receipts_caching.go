package sources

import (
	"context"
	"sync"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/sources/caching"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// A CachingReceiptsProvider caches successful receipt fetches from the inner
// ReceiptsProvider. It also avoids duplicate in-flight requests per block hash.
type CachingReceiptsProvider struct {
	inner InnerReceiptsProvider
	cache *caching.PreFetchCache[*ReceiptsHashPair]
	// lock fetching process for each block hash to avoid duplicate requests
	fetching   map[common.Hash]*sync.Mutex
	fetchingMu sync.Mutex // only protects map
}

func NewCachingReceiptsProvider(inner InnerReceiptsProvider, m caching.Metrics, cacheSize int) *CachingReceiptsProvider {
	return &CachingReceiptsProvider{
		inner:    inner,
		cache:    caching.NewPreFetchCache[*ReceiptsHashPair](m, "receipts", cacheSize),
		fetching: make(map[common.Hash]*sync.Mutex),
	}
}

func NewCachingRPCReceiptsProvider(client rpcClient, log log.Logger, config RPCReceiptsConfig, m caching.Metrics, cacheSize int) *CachingReceiptsProvider {
	return NewCachingReceiptsProvider(NewRPCReceiptsFetcher(client, log, config), m, cacheSize)
}

func (p *CachingReceiptsProvider) getOrCreateFetchingLock(blockHash common.Hash) *sync.Mutex {
	p.fetchingMu.Lock()
	defer p.fetchingMu.Unlock()
	if mu, ok := p.fetching[blockHash]; ok {
		return mu
	}
	mu := new(sync.Mutex)
	p.fetching[blockHash] = mu
	return mu
}

func (p *CachingReceiptsProvider) deleteFetchingLock(blockHash common.Hash) {
	p.fetchingMu.Lock()
	defer p.fetchingMu.Unlock()
	delete(p.fetching, blockHash)
}

// FetchReceipts fetches receipts for the given block and transaction hashes
// it expects that the inner FetchReceipts implementation handles validation
func (p *CachingReceiptsProvider) FetchReceipts(ctx context.Context, blockInfo eth.BlockInfo, txHashes []common.Hash, isForPreFetch bool) (types.Receipts, error, bool) {
	block := eth.ToBlockID(blockInfo)
	var isFull bool

	if v, ok := p.cache.Get(block.Number, !isForPreFetch); ok && v.blockHash == block.Hash {
		return v.receipts, nil, isFull
	}

	mu := p.getOrCreateFetchingLock(block.Hash)
	mu.Lock()
	defer mu.Unlock()
	// Other routine might have fetched in the meantime
	if v, ok := p.cache.Get(block.Number, !isForPreFetch); ok && v.blockHash == block.Hash {
		// we might have created a new lock above while the old
		// fetching job completed.
		p.deleteFetchingLock(block.Hash)
		return v.receipts, nil, isFull
	}

	isFull = p.cache.IsFull()
	if isForPreFetch && isFull {
		return nil, nil, true
	}

	r, err := p.inner.FetchReceipts(ctx, blockInfo, txHashes)
	if err != nil {
		return nil, err, isFull
	}

	p.cache.AddIfNotFull(block.Number, &ReceiptsHashPair{blockHash: block.Hash, receipts: r})
	// result now in cache, can delete fetching lock
	p.deleteFetchingLock(block.Hash)
	return r, nil, isFull
}

func (p *CachingReceiptsProvider) GetReceiptsCache() *caching.PreFetchCache[*ReceiptsHashPair] {
	return p.cache
}
