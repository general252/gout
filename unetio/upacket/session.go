package upacket

import (
	"bytes"
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/general252/gout/uerror"
	"github.com/general252/gout/ulog"
	"github.com/general252/gout/usafe"
)

const (
	PacketArraySize = 2048
)

type UdpSession struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	server     *UdpServer
	remoteAddr string

	pktArray []*MulUdpPacket

	bufferList *list.List
	muxBuffer  sync.Mutex // 缓冲区锁

	currentIndex          int // 当前访问的索引pktArray
	currentPacketWaitTime time.Time
	haveRecvFirstPacket   bool // 是否收到首个包
}

func NewUdpSession(server *UdpServer, clientAddr string) *UdpSession {
	rs := &UdpSession{
		server:     server,
		remoteAddr: clientAddr,
		bufferList: list.New(),
		pktArray:   make([]*MulUdpPacket, PacketArraySize),
		wg:         &sync.WaitGroup{},
	}
	rs.ctx, rs.cancel = context.WithCancel(context.TODO())

	rs.wg.Add(1)
	usafe.Go(func() {
		rs.routine()
	})

	return rs
}

func (c *UdpSession) Open() {

}

func (c *UdpSession) Close() {
	if c.cancel != nil {
		c.cancel()
		c.wg.Wait()
	}
}

func (c *UdpSession) routine() {
	defer c.wg.Done()
	for {
		var now = time.Now()
		select {
		case <-c.ctx.Done():
			return
		default:
			var sleep = true
			for i := 0; i < 1000; i++ {
				if c.mHandlePacketBuffer() != 0 {
					sleep = false
				} else {
					break
				}
			}
			for i := 0; i < 1000; i++ {
				if c.checkPacket(now) != 0 {
					sleep = false
				} else {
					break
				}
			}

			if sleep {
				time.Sleep(time.Millisecond * 2)
			}
		}
	}
}

func (c *UdpSession) mPopPacketBuffer() (*bytes.Buffer, error) {
	c.muxBuffer.Lock()
	defer c.muxBuffer.Unlock()

	if c.bufferList.Len() == 0 {
		return nil, nil
	}
	ulog.InfoFWithTimes("recv_packet_", 200, "udp buffer len ============%v, remote: %v",
		c.bufferList.Len(), c.remoteAddr)
	//ulog.InfoF("=====================%v", c.bufferList.Len())

	var elm = c.bufferList.Front()
	if elm == nil {
		return nil, uerror.WithMessage("list error. elm is nil")
	}

	obj, ok := elm.Value.(*bytes.Buffer)
	if !ok || obj == nil {
		return nil, fmt.Errorf("list elm error")
	}

	c.bufferList.Remove(elm)
	return obj, nil
}

func (c *UdpSession) mHandlePacketBuffer() int {
	buff, err := c.mPopPacketBuffer()
	if err != nil {
		ulog.ErrorF("mPopPacketBuffer fail %v", err)
	}
	if buff == nil {
		// 没有需要处理的数据
		return 0
	}

	// 解包, 保存
	for {
		var pkt UdpPacket
		if err := pkt.UnPacket(buff.Bytes()); err != nil {
			ulog.ErrorF("un packet fail %v", err)
			break
		}

		c.onParsePacket(&pkt)

		var index = pkt.PktSeq % PacketArraySize
		var objMulPacket = c.pktArray[index]
		if objMulPacket == nil {
			// 添加一个复合包
			objMulPacket = NewMulUdpPacket(int(pkt.PktCount))
			c.pktArray[index] = objMulPacket
		} else {
			// 包序号不一致, 存的包没有被读取, 新包被丢弃
			if objMulPacket.CheckPacket(&pkt) == false {
				ulog.WarnF("[###]packet array fill..., drop packet %v", pkt.PktSeq)
				break
			}
		}

		if err := objMulPacket.AddPacket(&pkt); err != nil {
			ulog.ErrorF("mul packet add packet fail %v", err)
			break
		}

		// 成功接收包
		c.onParsePacketSuccess(&pkt)
		break
	}

	return 1
}

func (c *UdpSession) onParsePacket(pkt *UdpPacket) {
	// 接收到一个包

}

func (c *UdpSession) onParsePacketSuccess(pkt *UdpPacket) {
	// 成功接收到一个包
	if c.haveRecvFirstPacket == false {
		c.currentIndex = int(pkt.PktSeq) % PacketArraySize
		c.currentPacketWaitTime = time.Now()

		c.haveRecvFirstPacket = true
	}
}

func (c *UdpSession) OnRecvPacket(packet *MulUdpPacket) {
	payload, _ := packet.GetPacketData()

	c.server.handRecvPacket(packet, payload)
}

func (c *UdpSession) OnPacketTimeout(packet *MulUdpPacket) {
	c.server.handTimeoutPacket(packet)
}

func (c *UdpSession) OnPacketLoss(index int) {
	ulog.WarnF("#### skip packet seq index %v", c.currentIndex)
	c.server.handLossPacket(index)
}

func (c *UdpSession) getIndex(index int) int {
	return index % PacketArraySize
}

func (c *UdpSession) checkPacket(now time.Time) (n int) {
	n = 0
	if c.haveRecvFirstPacket == false {
		// 还没有接收到数据包
		return
	}

	var objMulPacket = c.pktArray[c.currentIndex]
	if objMulPacket == nil {
		if c.pktArray[c.getIndex(c.currentIndex+1)] != nil && now.Sub(c.currentPacketWaitTime) > time.Second*4 {
			// 当前包超时(没有接收到), 访问下一个数据包
			c.OnPacketLoss(c.currentIndex)

			c.currentIndex = (c.currentIndex + 1) % PacketArraySize
			c.currentPacketWaitTime = now
			n = 1
		}
		return
	}

	for {
		if objMulPacket.IsRecvAllPacket() {
			// 接收完成
			c.OnRecvPacket(objMulPacket)
			break
		}

		if objMulPacket.IsTimeout(now, time.Second*3) {
			// 超时
			c.OnPacketTimeout(objMulPacket)
			break
		}

		return
	}

	c.pktArray[c.currentIndex] = nil
	c.currentIndex = (c.currentIndex + 1) % PacketArraySize
	c.currentPacketWaitTime = now
	n = 1
	return
}

func (c *UdpSession) OnRecvData(data []byte) {
	c.muxBuffer.Lock()
	defer c.muxBuffer.Unlock()

	var buff bytes.Buffer
	buff.Write(data)

	c.bufferList.PushBack(&buff)
}
