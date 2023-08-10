package uringbuffer

import (
	"sync"
)

type event struct {
	mutex sync.Mutex
	cond  *sync.Cond
	value bool
}

func newEvent() *event {
	cv := &event{}
	cv.cond = sync.NewCond(&cv.mutex)
	return cv
}

func (cv *event) set() {
	cv.cond.L.Lock()
	defer cv.cond.L.Unlock()

	cv.value = true
}

func (cv *event) signal() {
	if !cv.value {
		cv.set()
	}

	cv.cond.Broadcast()
}

func (cv *event) wait() {
	cv.cond.L.Lock()
	defer cv.cond.L.Unlock()

	if !cv.value {
		cv.cond.Wait()
	}

	cv.value = false
}
