package caching

import (
	"sync"

	"github.com/ethereum/go-ethereum/common/prque"
)

type PreFetchCache[V any] struct {
	m       Metrics
	label   string
	inner   map[uint64]V
	queue   *prque.Prque[uint64, V]
	lock    sync.Mutex
	maxSize int
}

func NewPreFetchCache[V any](m Metrics, label string, maxSize int) *PreFetchCache[V] {
	return &PreFetchCache[V]{
		m:       m,
		label:   label,
		inner:   make(map[uint64]V),
		queue:   prque.New[uint64, V](nil),
		maxSize: maxSize,
	}
}

func (v *PreFetchCache[V]) Add(key uint64, value V) bool {
	defer v.lock.Unlock()
	v.lock.Lock()
	if _, ok := v.inner[key]; ok {
		return false
	}
	v.queue.Push(value, -key)
	v.inner[key] = value
	if v.m != nil {
		v.m.CacheAdd(v.label, v.queue.Size(), false)
	}
	return true
}

func (v *PreFetchCache[V]) AddIfNotFull(key uint64, value V) (success bool, isFull bool) {
	defer v.lock.Unlock()
	v.lock.Lock()
	if _, ok := v.inner[key]; ok {
		return false, false
	}
	if v.queue.Size() >= v.maxSize {
		return false, true
	}
	v.queue.Push(value, -key)
	v.inner[key] = value
	if v.m != nil {
		v.m.CacheAdd(v.label, v.queue.Size(), false)
	}
	return true, false
}

func (v *PreFetchCache[V]) Get(key uint64) (V, bool) {
	defer v.lock.Unlock()
	v.lock.Lock()
	value, ok := v.inner[key]
	if v.m != nil {
		v.m.CacheGet(v.label, ok)
	}
	return value, ok
}

func (v *PreFetchCache[V]) RemoveAll() {
	defer v.lock.Unlock()
	v.lock.Lock()
	v.inner = make(map[uint64]V)
	v.queue.Reset()
}

func (v *PreFetchCache[V]) RemoveLessThan(p uint64) (isRemoved bool) {
	defer v.lock.Unlock()
	v.lock.Lock()
	for !v.queue.Empty() {
		_, qKey := v.queue.Peek()
		if -qKey < p {
			v.queue.Pop()
			delete(v.inner, -qKey)
			isRemoved = true
			continue
		}
		break
	}
	return
}
