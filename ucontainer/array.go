package ucontainer

import (
	"fmt"
	"sync"
)

// defaultArraySize is the default size for a array.
const defaultArraySize = 8

type Array[T any] struct {
	mu   sync.RWMutex
	vals []T
}

func NewArray[T any]() *Array[T] {
	return NewArraySize[T](defaultArraySize)
}

func NewArraySize[T any](size int) *Array[T] {
	return &Array[T]{
		vals: make([]T, 0, size),
	}
}

func (tis *Array[T]) Append(ts ...T) {
	tis.mu.Lock()
	defer tis.mu.Unlock()

	tis.vals = append(tis.vals, ts...)
}

func (tis *Array[T]) Add(value T) {
	tis.mu.Lock()
	defer tis.mu.Unlock()

	tis.vals = append(tis.vals, value)
}

func (tis *Array[T]) Set(index int, value T) error {
	tis.mu.Lock()
	defer tis.mu.Unlock()

	length := len(tis.vals)

	if index < 0 || index >= length {
		return fmt.Errorf("out range")
	}

	tis.vals[index] = value
	return nil
}

func (tis *Array[T]) Get(index int) (value T, ok bool) {
	tis.mu.RLock()
	defer tis.mu.RUnlock()

	length := len(tis.vals)

	if index < 0 || index >= length {
		return value, false
	}

	return tis.vals[index], true
}

// Delete 删除所有与value相等的元素
func (tis *Array[T]) Delete(value T, cmp func(a, b T) (equal bool)) {
	tis.mu.Lock()
	defer tis.mu.Unlock()

	var (
		length = len(tis.vals)
		j      = 0
	)

	for i := 0; i < length; i++ {
		if cmp(tis.vals[i], value) {
			tis.vals[j] = tis.vals[i]
			j++
		}
	}

	tis.vals = tis.vals[:j]
}

// DeleteByIndex 删除索引为index的元素
func (tis *Array[T]) DeleteByIndex(index int) {
	tis.mu.Lock()
	defer tis.mu.Unlock()

	var length = len(tis.vals)

	if index < 0 || index >= length {
		return
	}

	copy(tis.vals[index:], tis.vals[index+1:])

	tis.vals = tis.vals[:len(tis.vals)-1]
}

func (tis *Array[T]) Len() int {
	tis.mu.RLock()
	defer tis.mu.RUnlock()

	return len(tis.vals)
}

func (tis *Array[T]) Cap() int {
	tis.mu.RLock()
	defer tis.mu.RUnlock()

	return cap(tis.vals)
}

func (tis *Array[T]) Range(fn func(index int, value T) bool) {
	tis.mu.RLock()
	defer tis.mu.RUnlock()

	for i, v := range tis.vals {
		if !fn(i, v) {
			return
		}
	}
}

func (tis *Array[T]) ToSlice() []T {
	tis.mu.RLock()
	defer tis.mu.RUnlock()

	res := make([]T, len(tis.vals))
	copy(res, tis.vals)

	return res
}
