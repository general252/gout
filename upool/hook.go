package upool

import (
	"fmt"

	"github.com/general252/gout/uencode"
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
