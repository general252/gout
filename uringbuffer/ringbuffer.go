// Package ringbuffer contains a ring buffer.
package uringbuffer

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

// RingBuffer is a ring buffer. FROM github.com/aler9/gortsplib/v2/pkg/ringbuffer
type RingBuffer[T any] struct {
	size       uint64
	readIndex  uint64
	writeIndex uint64
	closed     int64
	buffer     []unsafe.Pointer
	event      *event
}

// New allocates a RingBuffer.
func New[T any](size uint64) (*RingBuffer[T], error) {
	// when writeIndex overflows, if size is not a power of
	// two, only a portion of the buffer is used.
	if (size & (size - 1)) != 0 {
		return nil, fmt.Errorf("size must be a power of two")
	}

	return &RingBuffer[T]{
		size:       size,
		readIndex:  1,
		writeIndex: 0,
		buffer:     make([]unsafe.Pointer, size),
		event:      newEvent(),
	}, nil
}

// Close makes Pull() return false.
func (r *RingBuffer[T]) Close() {
	atomic.StoreInt64(&r.closed, 1)
	r.event.signal()
}

// Reset restores Pull() behavior after a Close().
func (r *RingBuffer[T]) Reset() {
	for i := uint64(0); i < r.size; i++ {
		atomic.SwapPointer(&r.buffer[i], nil)
	}
	atomic.SwapUint64(&r.writeIndex, 0)
	r.readIndex = 1
	atomic.StoreInt64(&r.closed, 0)
}

// Push pushes data at the end of the buffer.
func (r *RingBuffer[T]) Push(data T) {
	writeIndex := atomic.AddUint64(&r.writeIndex, 1)
	i := writeIndex % r.size
	atomic.SwapPointer(&r.buffer[i], unsafe.Pointer(&data))
	r.event.signal()
}

// Pull pulls data from the beginning of the buffer.
func (r *RingBuffer[T]) Pull() (value T, ok bool) {
	for {
		i := r.readIndex % r.size
		res := (*T)(atomic.SwapPointer(&r.buffer[i], nil))
		if res == nil {
			if atomic.SwapInt64(&r.closed, 0) == 1 {
				return value, false
			}
			r.event.wait()
			continue
		}

		r.readIndex++
		return *res, true
	}
}
