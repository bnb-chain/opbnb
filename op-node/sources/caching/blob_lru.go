package caching

import (
	"math"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
)

type SizeFn func(value any) uint64

// SizeConstrainedCache is a cache where capacity is in bytes (instead of item count). When the cache
// is at capacity, and a new item is added, older items are evicted until the size
// constraint is met.
//
// OBS: This cache assumes that items are content-addressed: keys are unique per content.
// In other words: two Add(..) with the same key K, will always have the same value V.
type SizeConstrainedCache[K comparable, V any] struct {
	m       Metrics
	label   string
	size    uint64
	maxSize uint64
	sizeFn  SizeFn
	lru     *lru.Cache[K, V]
	lock    sync.Mutex
}

// NewSizeConstrainedCache creates a new size-constrained LRU cache.
func NewSizeConstrainedCache[K comparable, V any](m Metrics, label string, maxSize uint64, sizeFn SizeFn) *SizeConstrainedCache[K, V] {
	cache, _ := lru.New[K, V](math.MaxInt)
	return &SizeConstrainedCache[K, V]{
		m:       m,
		label:   label,
		size:    0,
		maxSize: maxSize,
		sizeFn:  sizeFn,
		lru:     cache,
	}
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
// OBS: This cache assumes that items are content-addressed: keys are unique per content.
// In other words: two Add(..) with the same key K, will always have the same value V.
// OBS: The value is _not_ copied on Add, so the caller must not modify it afterwards.
func (c *SizeConstrainedCache[K, V]) Add(key K, value V) (evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// If there is already a item with the same key in the cache,
	// remove the old element first, then write in the new item.
	if oldValue, ok := c.lru.Get(key); ok {
		c.lru.Remove(key)
		c.size -= c.sizeFn(oldValue)
	}

	// Evict the oldest items if cache is overflow
	targetSize := c.size + c.sizeFn(value)
	for targetSize > c.maxSize {
		evicted = true
		_, v, ok := c.lru.RemoveOldest()
		if !ok {
			// list is now empty. Break
			break
		}
		targetSize -= c.sizeFn(v)
	}
	c.size = targetSize

	c.lru.Add(key, value)
	if c.m != nil {
		c.m.CacheAdd(c.label, c.lru.Len(), evicted)
	}

	return evicted
}

// Get looks up a key's value from the cache.
func (c *SizeConstrainedCache[K, V]) Get(key K) (V, bool) {
	c.lock.Lock()
	value, ok := c.lru.Get(key)
	c.lock.Unlock()

	if c.m != nil {
		c.m.CacheGet(c.label, ok)
	}
	return value, ok
}
