package loss_detection

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/general252/gout/uarray"
)

func TestNewSeqLossCheck(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var wg = &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	var cq = NewSeqLossCheck(wg, ctx, func(lossSeq uint32) {

	})

	const TEST = BitSize - 10

	var arr []uint64
	arr = append(arr, 5)
	arr = append(arr, 6)
	arr = append(arr, 500)
	arr = append(arr, 1000)
	arr = append(arr, 1003)
	arr = append(arr, 1056)
	arr = append(arr, 1033)
	arr = append(arr, 3000)
	arr = append(arr, 6000)
	arr = append(arr, 98295)
	arr = append(arr, 98298)
	arr = append(arr, uint64(uint32(98302%BitSize)))

	log.Printf("需要丢的包: %v\n\n", arr)

	for i := TEST; i < TEST+200000; i++ {
		var item = uint32(i % BitSize)
		if uarray.ContainsUint(arr, uint64(item)) >= 0 {
			log.Printf("remote %v", item)
		} else {
			cq.Add(item)
		}

		if 0 == i%10 {
			time.Sleep(time.Millisecond * 5)
		}
	}
}
