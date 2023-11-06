package dict

import (
	"sync"

	"golang.org/x/exp/maps"
)

type Dict[K comparable, T any] struct {
	sync.RWMutex
	data map[K]T
}

func NewDict[K comparable, T any]() *Dict[K, T] {
	return &Dict[K, T]{data: map[K]T{}}
}

func (d *Dict[K, T]) Del(key K) (item T) {
	d.Lock()
	defer d.Unlock()

	var ok bool
	if item, ok = d.data[key]; ok {
		delete(d.data, key)
	}
	return
}

func (d *Dict[K, T]) Set(key K, item T) {
	d.Lock()
	defer d.Unlock()

	d.data[key] = item
}

func (d *Dict[K, T]) Get(key K) (item T, ok bool) {
	d.RLock()
	defer d.RUnlock()

	item, ok = d.data[key]
	return
}

func (d *Dict[K, T]) Len() int {
	d.RLock()
	defer d.RUnlock()

	return len(d.data)
}

func (d *Dict[K, T]) Has(key K) bool {
	_, ok := d.Get(key)

	return ok
}

func (d *Dict[K, T]) Loop(fn func(K, T) bool) {
	d.RLock()
	data := d.data
	d.RUnlock()

	for key, value := range data {
		if !fn(key, value) {
			return
		}
	}
}

func (d *Dict[K, T]) Values() []T {
	d.RLock()
	defer d.RUnlock()

	return maps.Values(d.data)
}

func (d *Dict[K, T]) Keys() []K {
	d.RLock()
	defer d.RUnlock()

	return maps.Keys(d.data)
}

func (d *Dict[K, T]) Clear() {
	d.Lock()
	defer d.Unlock()

	clear(d.data)
}
