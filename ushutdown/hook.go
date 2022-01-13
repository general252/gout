package ushutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type StopHandler func()

type Hook interface {
	WithSignals(signals ...os.Signal) Hook
	WithStopHandler(handles ...StopHandler) Hook

	GracefulStop()
	GracefulStopContext(ctx context.Context)
}

func NewHook() Hook {
	return innerNewHook()
}

type hook struct {
	mux sync.Mutex

	quitChan chan os.Signal
	handles  []StopHandler
}

func innerNewHook() Hook {
	c := &hook{
		quitChan: make(chan os.Signal, 5),
	}

	return c.WithSignals(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
}

func (c *hook) WithSignals(signals ...os.Signal) Hook {
	for _, s := range signals {
		signal.Notify(c.quitChan, s)
	}

	return c
}

func (c *hook) WithStopHandler(handles ...StopHandler) Hook {
	c.handles = append(c.handles, handles...)

	return c
}

func (c *hook) GracefulStop() {
	c.innerGracefulStopContext(context.TODO())
}

func (c *hook) GracefulStopContext(ctx context.Context) {
	c.innerGracefulStopContext(ctx)
}

func (c *hook) innerGracefulStopContext(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-c.quitChan:
	}

	signal.Stop(c.quitChan)

	for _, f := range c.handles {
		f()
	}
}
