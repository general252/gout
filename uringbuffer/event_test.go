package uringbuffer

import (
	"testing"
	"time"
)

func TestSign(t *testing.T) {
	e := newEvent()

	go func() {
		e.wait()
	}()

	time.Sleep(time.Second * 5)
	e.signal()

	time.Sleep(time.Second)
}
