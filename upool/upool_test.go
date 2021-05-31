package upool

import (
	"fmt"
	"github.com/general252/gout/unet"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func ExampleNewUPool() {
	pool := NewUPool(WithPoolSize(200), WithSyncWaitTime(time.Second))
	defer pool.Close()

	hook := NewDefaultHook(func(data *PoolItem) error {
		log.Println(string(unet.FormatJsonObject(data)))
		return nil
	})

	pool.AddHook(hook)
	//defer pool.RemoveHook(hook)

	for i := 0; i < 1000; i++ {
		if err := pool.Write(&PoolItem{
			Level: 0,
			Msg:   fmt.Sprintf("msg %v", i),
			When:  time.Now(),
		}); err != nil {
			log.Println(err)
		} else {
			log.Println("ok")
		}
	}

	// output:
	//
}
