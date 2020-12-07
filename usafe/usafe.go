package usafe

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"
)

type HandleRecover func(err interface{})

var defaultHandleRecover HandleRecover = func(err interface{}) {
	fmt.Printf("Error in Go routine: %s\n", err)
	fmt.Printf("Stack: %s\n", debug.Stack())
}

func SetDefaultHandleRecover(handleRecover HandleRecover) {
	defaultHandleRecover = handleRecover
}

// Go starts a recoverable goroutine.
func Go(goroutine func()) {
	GoWithRecover(goroutine, defaultHandleRecover)
}

// GoWithRecover starts a recoverable goroutine using given customRecover() function.
func GoWithRecover(goroutine func(), handleRecover HandleRecover) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				handleRecover(err)
			}
		}()
		goroutine()
	}()
}

// GoRetry 重试funRetry
func GoRetry(ctx context.Context, interval time.Duration, retryTimes int, funRetry func() error, funResult func(err error)) {
	Go(func() {
		// 等待ticker时间
		var waitSignal = func(cont context.Context, timer *time.Ticker) (timeout bool) {
			select {
			case <-cont.Done():
				return false
			case <-timer.C:
				return true
			}
		}

		t := time.NewTicker(interval)
		defer func() {
			t.Stop()
		}()

		var err error
		for i := 0; i < retryTimes; i++ {
			if waitSignal(ctx, t) == false {
				return
			}

			if err = funRetry(); err != nil {
				continue
			}

			break
		}

		funResult(err)
	})
}
