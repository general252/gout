package loss_detection

import (
	"container/list"
	"context"
	"github.com/willf/bitset"
	"log"
	"math"
	"sync"
	"time"
)

const (
	BitSize    = math.MaxUint16 + math.MaxUint16/2
	InvalidPos = math.MaxUint64 - 10
	SeqIdle    = time.Second * 5
)

type seqAndTime struct {
	seq  uint32
	pos  uint
	time time.Time // 添加的时间
}

type HandLossSeq func(lossSeq uint32)

type SeqCheck struct {
	hand         HandLossSeq
	seqSet       *bitset.BitSet
	seqTimeList  *list.List
	lastCheckPos uint

	lastAddSeq  uint32
	mLastAddSeq uint32 // 辅助, 标记是否有新的seq添加

	mux sync.Mutex
}

// NewSeqLossCheck seq丢失检测, Add(seq uint32)是有序的. 速度限制BitSize32个/SeqIdle
func NewSeqLossCheck(wg *sync.WaitGroup, ctx context.Context, handLoss HandLossSeq) *SeqCheck {
	var rs = &SeqCheck{
		hand:         handLoss,
		seqSet:       bitset.New(BitSize),
		seqTimeList:  list.New(),
		lastCheckPos: InvalidPos,
	}

	if wg != nil {
		wg.Add(1)
	}
	go func() {
		if wg != nil {
			defer wg.Done()
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				rs.Pool()
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()

	return rs
}

func (c *SeqCheck) Add(seq uint32) {
	c.mux.Lock()
	defer c.mux.Unlock()

	var pos = seq % BitSize

	// 检查
	if c.seqSet.Test(uint(pos)) {
		log.Printf("warn, may to fast")
		return
	}

	c.seqSet.Set(uint(pos))
	if c.lastCheckPos == InvalidPos {
		c.lastCheckPos = uint(pos)
	}

	c.lastAddSeq = seq

	c.addMark()
	c.checkMark()
}

func (c *SeqCheck) Pool() {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.lastCheckPos == InvalidPos {
		return
	}

	c.addMark()
	c.checkMark()
}

func (c *SeqCheck) addMark() {
	if c.lastAddSeq == c.mLastAddSeq {
		return
	}

	var seq = c.lastAddSeq
	var pos = seq % BitSize

	var now = time.Now()
	// 判断实现需要记录seq的时间点

	var isAdd = false
	var elm = c.seqTimeList.Back()
	if elm == nil {
		// 首个seq
		isAdd = true
	} else {
		// 到达间隔时间, 添加一个seq时间
		if now.Sub(elm.Value.(*seqAndTime).time) > time.Second {
			isAdd = true
		}
	}

	if isAdd {
		c.mLastAddSeq = c.lastAddSeq

		c.seqTimeList.PushBack(&seqAndTime{
			seq:  seq,
			pos:  uint(pos),
			time: now,
		})
	}
}

func (c *SeqCheck) checkMark() {
	// 丢包检测
	var elm = c.seqTimeList.Front()
	if elm == nil {
		return
	}

	var now = time.Now()
	item := elm.Value.(*seqAndTime)
	if now.Sub(item.time) < SeqIdle {
		return
	}
	// 到时了

	// 再次循环
	if c.lastCheckPos > item.pos {
		for i := c.lastCheckPos; i < BitSize; i++ {
			if c.seqSet.Test(i) {
				c.seqSet.Clear(i)
			} else {
				log.Printf("loss 1 seq: %v", i)
				if c.hand != nil {
					c.hand(uint32(i))
				}
			}
		}
		c.lastCheckPos = 0
	}

	for i := c.lastCheckPos; i <= item.pos; i++ {
		if c.seqSet.Test(i) {
			c.seqSet.Clear(i)
		} else {
			log.Printf("loss 2 seq: %v", i)
			if c.hand != nil {
				c.hand(uint32(i))
			}
		}
	}
	c.lastCheckPos = item.pos + 1

	c.seqTimeList.Remove(elm)
}
