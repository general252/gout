package udp_packet

import (
	"container/list"
	"context"
	"fmt"
	"github.com/general252/gout/laboratory/loss_detection"
	"github.com/general252/gout/ulog"
	"github.com/general252/gout/unetio"
	"net"
	"sync"
	"time"
)

func Serialization(data []byte) []UdpPacket {
	var seq = unetio.GetSeq()

	if len(data) <= unetio.UdpPacketPayloadMaxSize {
		var pkt = UdpPacket{
			PktSeq:   seq,
			PktCount: 1,
			PktIndex: 0,
			PktType:  unetio.PktTypeData,
		}
		_, _ = pkt.Payload.Write(data)

		return []UdpPacket{pkt}
	} else {
		var dataLength = len(data)

		var n = dataLength / unetio.UdpPacketPayloadMaxSize
		if dataLength%unetio.UdpPacketPayloadMaxSize > 0 {
			n += 1
		}

		var result = make([]UdpPacket, 0, n)
		for i := 0; i < n; i++ {
			var offset = i * unetio.UdpPacketPayloadMaxSize
			var end = offset + unetio.UdpPacketPayloadMaxSize
			if end > dataLength {
				end = dataLength
			}

			var pkt = UdpPacket{
				PktSeq:   seq,
				PktCount: uint16(n),
				PktIndex: uint16(i),
			}
			_, _ = pkt.Payload.Write(data[offset:end])

			result = append(result, pkt)
		}

		return result
	}
}

func Deserialization(pktData []byte) (*UdpPacket, error) {
	var pkt UdpPacket
	if err := pkt.UnPacket(pktData); err != nil {
		return nil, err
	} else {
		return &pkt, nil
	}
}

type stUdpPacketBuffer struct {
	addr   net.UDPAddr
	buffer []byte
}

type PayloadHandle func(payload []byte)

type UDPPacketFactory struct {
	ctx              context.Context
	handle           PayloadHandle
	mulUdpPacketList sync.Map // "ip:port_pktSeq" -> *mulUdpPacket

	udpPktBufferItemList      *list.List // 缓存区队列, 不直接处理接收的数据, 目的是防止解包耗时导致丢包
	mutexUdpPktBufferItemList sync.Mutex // 缓冲区锁

	lossCheck *loss_detection.SeqCheck
}

func NewFactoryPacket(ctx context.Context, handle PayloadHandle) *UDPPacketFactory {
	var rs = &UDPPacketFactory{
		ctx:                  ctx,
		handle:               handle,
		udpPktBufferItemList: list.New(),
	}

	rs.lossCheck = loss_detection.NewSeqLossCheck(nil, ctx, func(lossSeq uint32) {

	})

	go rs.routine()

	return rs
}

func (c *UDPPacketFactory) showInfo(info string) {
	fmt.Println(info)
	c.mulUdpPacketList.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
}

func (c *UDPPacketFactory) getMulPacket(connSeqId string) *mulUdpPacket {
	tmpObj, ok := c.mulUdpPacketList.Load(connSeqId)
	if !ok {
		return nil
	}

	obj, ok := tmpObj.(*mulUdpPacket)
	if !ok || obj == nil {
		ulog.ErrorF("inner error. %v", tmpObj)
		return nil
	}

	return obj
}

func (c *UDPPacketFactory) delMulPacket(connSeqId string) {
	if obj := c.getMulPacket(connSeqId); obj != nil {
		//ulog.InfoF("del packet %v\n", connSeqId)
	} else {
		ulog.InfoF("del packet %v", connSeqId)
	}

	c.mulUdpPacketList.Delete(connSeqId)
}

func (c *UDPPacketFactory) innerPushMulUdpPacket(pkt *UdpPacket) *mulUdpPacket {
	var connSeq = pkt.ConnSeqId()

	var findObj = c.getMulPacket(connSeq)
	if findObj != nil {
		findObj.RecvPacket(pkt)

		return findObj
	} else {
		var addObj = &mulUdpPacket{
			addTime: time.Now(),
		}
		addObj.RecvPacket(pkt)

		c.mulUdpPacketList.Store(connSeq, addObj)
		c.lossCheck.Add(pkt.Seq())

		return addObj
	}
}

func (c *UDPPacketFactory) mPushPacketData(addr net.UDPAddr, pktData []byte) {
	if pkt, err := Deserialization(pktData); err != nil {
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
			c.handle(pkt.PayloadData())
			c.lossCheck.Add(pkt.Seq())
		}
	}
}

func (c *UDPPacketFactory) PushPacketData(addr net.UDPAddr, pktData []byte) {
	c.mutexUdpPktBufferItemList.Lock()
	defer c.mutexUdpPktBufferItemList.Unlock()

	var tmpBuffer = make([]byte, len(pktData))
	copy(tmpBuffer, pktData)

	c.udpPktBufferItemList.PushBack(&stUdpPacketBuffer{
		addr:   addr,
		buffer: tmpBuffer,
	})
}

func (c *UDPPacketFactory) mHandleAddItemBuffer() {
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
func (c *UDPPacketFactory) checkTimeout() {
	var now = time.Now()
	var keys []interface{}

	c.mulUdpPacketList.Range(func(key, value interface{}) bool {
		obj, ok := value.(*mulUdpPacket)
		if !ok || nil == obj {
			keys = append(keys, key)
			return true
		}

		// 分包接收完整
		if obj.IsRecvAllPacket() {
			c.handle(obj.GetPacketData())
			keys = append(keys, key)
			return true
		}

		// 数据接收超时
		if obj.IsTimeout(now) {
			keys = append(keys, key)
			ulog.ErrorF("timeout. %v, value: \n%v", key, obj)
			return true
		}

		return true
	})

	for _, key := range keys {
		connSeqId, ok := key.(string)
		if !ok {
			ulog.ErrorF("inner error. %v", key)
		} else {
			c.delMulPacket(connSeqId)
		}
	}
}

func (c *UDPPacketFactory) routine() {
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
