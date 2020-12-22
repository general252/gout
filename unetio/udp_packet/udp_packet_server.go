package udp_packet

import (
	"container/list"
	"context"
	"fmt"
	"github.com/general252/gout/laboratory/loss_detection"
	"github.com/general252/gout/ulog"
	"net"
	"sort"
	"sync"
	"time"
)

type stUdpPacketBuffer struct {
	addr   net.UDPAddr
	buffer []byte
}

type PayloadHandle func(pktHeadInfo *MulUdpPacket, payload []byte)
type TimeoutHandle func(pktHeadInfo *MulUdpPacket)

type UdpPacketServer struct {
	ctx           context.Context
	handlePayload PayloadHandle
	handleTimeout TimeoutHandle

	mulUdpPacketList sync.Map // "ip:port_pktSeq" -> *MulUdpPacket

	udpPktBufferItemList      *list.List // 缓存区队列, 不直接处理接收的数据, 目的是防止解包耗时导致丢包
	mutexUdpPktBufferItemList sync.Mutex // 缓冲区锁

	lossCheck *loss_detection.SeqCheck
}

func NewUdpPacketServer(ctx context.Context, handlePayload PayloadHandle, handleTime TimeoutHandle, handleLoss loss_detection.HandLossSeq) *UdpPacketServer {
	var rs = &UdpPacketServer{
		ctx:                  ctx,
		handlePayload:        handlePayload,
		handleTimeout:        handleTime,
		udpPktBufferItemList: list.New(),
	}

	rs.lossCheck = loss_detection.NewSeqLossCheck(nil, ctx, func(lossSeq uint32) {
		if handleLoss != nil {
			handleLoss(lossSeq)
		}
	})

	go rs.routine()

	return rs
}

func (c *UdpPacketServer) showInfo(info string) {
	fmt.Println(info)
	c.mulUdpPacketList.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
}

func (c *UdpPacketServer) getMulPacket(connSeqId string) *MulUdpPacket {
	tmpObj, ok := c.mulUdpPacketList.Load(connSeqId)
	if !ok {
		return nil
	}

	obj, ok := tmpObj.(*MulUdpPacket)
	if !ok || obj == nil {
		ulog.ErrorF("inner error. %v", tmpObj)
		return nil
	}

	return obj
}

func (c *UdpPacketServer) delMulPacket(connSeqId string) {
	if obj := c.getMulPacket(connSeqId); obj != nil {
		//ulog.InfoF("del packet %v\n", connSeqId)
	} else {
		ulog.InfoF("del packet %v", connSeqId)
	}

	c.mulUdpPacketList.Delete(connSeqId)
}

func (c *UdpPacketServer) innerPushMulUdpPacket(pkt *UdpPacket) *MulUdpPacket {
	var connSeq = pkt.ConnSeqId()

	var findObj = c.getMulPacket(connSeq)
	if findObj != nil {
		findObj.RecvPacket(pkt)

		return findObj
	} else {
		var addObj = &MulUdpPacket{
			addTime: time.Now(),
		}
		addObj.RecvPacket(pkt)

		c.mulUdpPacketList.Store(connSeq, addObj)
		c.lossCheck.Add(pkt.Seq())

		return addObj
	}
}

func (c *UdpPacketServer) mPushPacketData(addr net.UDPAddr, pktData []byte) {
	if pkt, err := c.Deserialization(pktData); err != nil {
		ulog.ErrorF("Deserialization fail. %v %v", addr, err)
	} else {
		pkt.SetRemoteAddr(addr)

		if pkt.IsMulPacket() {
			var mulPkt = c.innerPushMulUdpPacket(pkt)
			_ = mulPkt

			// 检查是否接受完成
			//if mulPkt.IsRecvAllPacket() {
			//	// 通知
			//	c.handle(mulPkt.GetPacketData())
			//	// 删除
			//	c.delMulPacket(pkt.ConnSeqId())
			//}
		} else {
			if true {
				// 所有的包都需要缓存一下, 纠正包顺序
				var mulPkt = c.innerPushMulUdpPacket(pkt)
				_ = mulPkt
			} else {
				var objMulPkt = &MulUdpPacket{addTime: time.Now()}
				objMulPkt.RecvPacket(pkt)
				c.handlePayload(objMulPkt, pkt.PayloadData())

				c.lossCheck.Add(pkt.Seq())
			}
		}
	}
}

func (c *UdpPacketServer) PushPacketData(addr net.UDPAddr, pktData []byte) {
	c.mutexUdpPktBufferItemList.Lock()
	defer c.mutexUdpPktBufferItemList.Unlock()

	var tmpBuffer = make([]byte, len(pktData))
	copy(tmpBuffer, pktData)

	c.udpPktBufferItemList.PushBack(&stUdpPacketBuffer{
		addr:   addr,
		buffer: tmpBuffer,
	})
}

func (c *UdpPacketServer) mHandleAddItemBuffer() {
	c.mutexUdpPktBufferItemList.Lock()
	defer c.mutexUdpPktBufferItemList.Unlock()

	for c.udpPktBufferItemList.Len() > 0 {
		var elm = c.udpPktBufferItemList.Front()
		if elm == nil {
			break
		}

		obj, ok := elm.Value.(*stUdpPacketBuffer)
		if ok {
			c.mPushPacketData(obj.addr, obj.buffer)
		} else {
			ulog.ErrorF("inner error. not *stUdpPacketBuffer")
		}

		c.udpPktBufferItemList.Remove(elm)
	}
}

// 检查分包接收完成或接收超时
func (c *UdpPacketServer) checkTimeout() {
	var now = time.Now()
	var keys []interface{}
	var vlus []*MulUdpPacket

	c.mulUdpPacketList.Range(func(key, value interface{}) bool {
		obj, ok := value.(*MulUdpPacket)
		if !ok || nil == obj {
			keys = append(keys, key)
			return true
		}

		// 分包接收完整
		if obj.IsRecvAllPacket() && obj.IsTimeBuffer(now) {
			keys = append(keys, key)
			vlus = append(vlus, obj)
			// c.handlePayload(obj, obj.GetPacketData())
			return true
		}

		// 数据接收超时
		if obj.IsTimeout(now) {
			keys = append(keys, key)
			c.handleTimeout(obj)
			return true
		}

		return true
	})

	sort.Slice(vlus, func(i, j int) bool {
		return vlus[i].addTime.Sub(vlus[j].addTime) > 0
	})

	for i := 0; i < len(vlus); i++ {
		c.handlePayload(vlus[i], vlus[i].GetPacketData())
	}

	for _, key := range keys {
		connSeqId, ok := key.(string)
		if !ok {
			ulog.ErrorF("inner error. %v", key)
		} else {
			c.delMulPacket(connSeqId)
		}
	}
}

func (c *UdpPacketServer) routine() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.mHandleAddItemBuffer()
			c.checkTimeout()

			time.Sleep(time.Millisecond * 5)
		}
	}
}

func (c *UdpPacketServer) Deserialization(pktData []byte) (*UdpPacket, error) {
	var pkt UdpPacket
	if err := pkt.UnPacket(pktData); err != nil {
		return nil, err
	} else {
		return &pkt, nil
	}
}
