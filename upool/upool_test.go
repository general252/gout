package upool

import (
	"github.com/general252/gout/unet"
	"log"
	"time"
)

func ExampleNewUPool() {
	pool := NewUPool(1000)
	defer pool.Close()

	hook := NewDefaultHook(func(data *PoolItem) error {
		log.Println(string(unet.FormatJsonObject(data)))
		return nil
	})

	pool.AddHook(hook)
	defer pool.RemoveHook(hook)

	_ = pool.Write(&PoolItem{
		Level: 0,
		Msg:   "",
		When:  time.Now(),
	})

	time.Sleep(time.Second)

	// output:
	//
}
