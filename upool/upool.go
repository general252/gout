package upool

import (
	"context"
	"fmt"
	"github.com/general252/gout/uencode"
	"sync"
	"time"
)

type PoolHook interface {
	ID() string
	Write(data *PoolItem) error
}

type HookPoolItem func(data *PoolItem) error
type defaultHook struct {
	id string
	h  HookPoolItem
}

func NewDefaultHook(h HookPoolItem) *defaultHook {
	return &defaultHook{
		id: uencode.UUID(),
		h:  h,
	}
}

func (tis *defaultHook) ID() string {
	return tis.id
}

func (tis *defaultHook) Write(data *PoolItem) error {
	if tis.h == nil {
		return fmt.Errorf("handle is nil")
	}

	return tis.h(data)
}

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
}

func NewUPool(s int64) *UPool {
	ctx, cancel := context.WithCancel(context.TODO())
	c := &UPool{
		msgChan: make(chan *PoolItem, s),
		ctx:     ctx,
		cancel:  cancel,
		pool: &sync.Pool{
			New: func() interface{} {
				return &PoolItem{}
			},
		},
	}

	c.wg.Add(1)
	go c.start()

	return c
}

func (tis *UPool) start() {
	defer func() {
		tis.wg.Done()
	}()

	for {
		select {
		case <-tis.ctx.Done():
			return
		case msg := <-tis.msgChan:
			tis.hooks.Range(func(key, value interface{}) bool {
				h, ok := value.(PoolHook)
				if ok && h != nil {
					_ = h.Write(msg)
				}
				return true
			})

			tis.pool.Put(msg)
		}
	}
}

func (tis *UPool) Close() {
	tis.cancel()
	tis.wg.Wait()
}

func (tis *UPool) AddHook(hook PoolHook) {
	tis.hooks.Store(hook.ID(), hook)
}

func (tis *UPool) RemoveHook(hook PoolHook) {
	tis.hooks.Delete(hook.ID())
}

func (tis *UPool) Write(data *PoolItem) error {
	var size = cap(tis.msgChan)
	var count = len(tis.msgChan)
	if count+100 > size {
		return fmt.Errorf("pool full %v / %v", count, size)
	}

	obj := tis.pool.Get().(*PoolItem)
	{
		obj.Msg = data.Msg
		obj.Level = data.Level
		obj.When = data.When
	}
	tis.msgChan <- obj

	return nil
}
