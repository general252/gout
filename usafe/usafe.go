package usafe

import (
	"fmt"
	"runtime/debug"
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
