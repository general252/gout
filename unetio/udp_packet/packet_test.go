package udp_packet

import (
	"bytes"
	"encoding/json"
	"github.com/general252/gout/unetio"
	"math"
	"math/rand"
	"testing"
	"time"
)

// 打包解包测试
func TestUdpPacket_Packet(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	data := []byte("hello world")
	count := uint16(rand.Int()%math.MaxInt8) + 23
	var srcPkt = &UdpPacket{
		PktSeq:   rand.Uint32(),
		PktCount: count,
		PktIndex: count - 1,
		PktType:  unetio.PktTypeData,
	}
	_, _ = srcPkt.Payload.Write(data)

	tmpPacketData := srcPkt.Packet()

	var dstPkt = &UdpPacket{}
	if err := dstPkt.UnPacket(tmpPacketData); err != nil {
		t.Error(err)
	}

	srcJson, _ := json.Marshal(srcPkt)
	dstJson, _ := json.Marshal(dstPkt)
	if bytes.Compare(srcJson, dstJson) != 0 {
		t.Errorf("fail")
	}

	t.Logf("success")
}
