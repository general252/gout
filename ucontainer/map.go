package ucontainer

import (
	"sync"
)

type hashAble interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}

type Map[K hashAble, V any] struct {
	m sync.Map
}

func NewMap[K hashAble, V any]() *Map[K, V] {
	return &Map[K, V]{}
}

func (tis *Map[K, V]) Delete(key K) {
	tis.m.Delete(key)
}

func (tis *Map[K, V]) Store(key K, value V) {
	tis.m.Store(key, value)
}

func (tis *Map[K, V]) Range(f func(key K, value V) bool) {
	tis.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

func (tis *Map[K, V]) Load(key K) (V, bool) {
	v, ok := tis.m.Load(key)
	if !ok {
		return *new(V), false
	}

	value, ok := v.(V)
	if !ok {
		return *new(V), false
	}

	return value, true
}

func (tis *Map[K, V]) IsExist(key K) bool {
	_, ok := tis.m.Load(key)

	return ok
}
