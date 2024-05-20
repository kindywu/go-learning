package rwmap

import (
	"errors"
	"sync"
)

type RWMap[K comparable, V any] struct {
	rw sync.RWMutex
	m  map[K]V
}

func NewRWMap[K comparable, V any](n int) *RWMap[K, V] {
	return &RWMap[K, V]{
		m: make(map[K]V),
	}
}

func (m *RWMap[K, V]) Len(k K) int {
	m.rw.RLock()
	defer m.rw.RUnlock()
	return len(m.m)
}

func (m *RWMap[K, V]) Each(fn func(k K, v V) error) error {
	if fn == nil {
		return errors.New("fn can't be nil")
	}
	m.rw.RLock()
	defer m.rw.RUnlock()
	for k, v := range m.m {
		if err := fn(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (m *RWMap[K, V]) Get(k K) (V, bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	v, exist := m.m[k]
	return v, exist
}

func (m *RWMap[K, V]) Set(k K, v V) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.m[k] = v
}

func (m *RWMap[K, V]) Delete(k K) {
	m.rw.Lock()
	defer m.rw.Unlock()
	delete(m.m, k)
}
