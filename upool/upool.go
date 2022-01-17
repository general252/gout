package upool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type PoolItem struct {
	Level    int         `json:"level"`
	Msg      string      `json:"msg"`
	When     time.Time   `json:"when"`
	UserData interface{} `json:"user_data"`
}

type UPool struct {
	pool    *sync.Pool
	msgChan chan *PoolItem

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	hooks sync.Map

	opt poolOptions
}

// NewUPool new pool, s: pool buffer size, 1000
func NewUPool(opts ...PoolOption) *UPool {
	var poolOpt = poolOptions{
		poolSize:      defaultPoolSize,
		syncWriteIdle: defaultSyncWriteIdle,
	}
	for _, opt := range opts {
		opt.apply(&poolOpt)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	c := &UPool{
		msgChan: make(chan *PoolItem, poolOpt.poolSize),
		ctx:     ctx,
		cancel:  cancel,
		pool: &sync.Pool{
			New: func() interface{} {
				return &PoolItem{}
			},
		},
		opt: poolOpt,
	}

	c.wg.Add(1)
	go c.start()

	return c
}

func (tis *UPool) Close() {
	tis.cancel()
	tis.wg.Wait()
}

// AddHook add hook
func (tis *UPool) AddHook(hook PoolHook) {
	tis.hooks.Store(hook.ID(), hook)
}

// RemoveHook remove hook by hook id
func (tis *UPool) RemoveHook(hook PoolHook) {
	tis.hooks.Delete(hook.ID())
}

// Write if pool is full will error
func (tis *UPool) Write(data *PoolItem) error {
	obj := tis.pool.Get().(*PoolItem)
	{
		obj.Msg = data.Msg
		obj.Level = data.Level
		obj.When = data.When
	}

	select {
	case tis.msgChan <- obj:
		// ok
	default:
		tis.pool.Put(obj)
		return fmt.Errorf("write fail")
	}

	return nil
}

// WriteSync write item, if pool is full will wait chan
func (tis *UPool) WriteSync(data *PoolItem) error {
	obj := tis.pool.Get().(*PoolItem)
	{
		obj.Msg = data.Msg
		obj.Level = data.Level
		obj.When = data.When
	}

	select {
	case tis.msgChan <- obj:
		// ok
		return nil
	case <-time.After(tis.opt.syncWriteIdle):
		tis.pool.Put(obj)
		return fmt.Errorf("timeout")
	}
}

func (tis *UPool) write(data *PoolItem) {
	tis.hooks.Range(func(key, value interface{}) bool {
		h, ok := value.(PoolHook)
		if ok && h != nil {
			_ = h.Write(data)
		}
		return true
	})
}

// start go routine
func (tis *UPool) start() {
	defer func() {
		tis.wg.Done()
	}()

	var handleMsg = func() {
		// 处理剩余的信息
		n := len(tis.msgChan)
		for i := 0; i < n; i++ {
			msg := <-tis.msgChan

			tis.write(msg)
			tis.pool.Put(msg)
		}
	}

	for {
		select {
		case <-tis.ctx.Done():
			handleMsg()
			return
		case msg := <-tis.msgChan:
			tis.write(msg)
			tis.pool.Put(msg)

			handleMsg()
		}
	}
}
