package udp_packet

import (
	"bytes"
	"math"
	"testing"
)

func TestSerialization(t *testing.T) {
	var srcData = make([]byte, 2020)
	for i := 0; i < len(srcData); i++ {
		srcData[i] = uint8(i % math.MaxUint8)
	}

	var dstPktList []*UdpPacket

	srcPktList := Serialization(srcData)
	for i := 0; i < len(srcPktList); i++ {
		tmpPkt, err := Deserialization(srcPktList[i].Packet())
		if err != nil {
			t.Error(err)
		}

		dstPktList = append(dstPktList, tmpPkt)
	}

	var dstData bytes.Buffer
	for i := 0; i < len(dstPktList); i++ {
		dstData.Write(dstPktList[i].PayloadData())
	}

	if bytes.Compare(srcData, dstData.Bytes()) != 0 {
		t.Errorf("fail")
	}

	t.Logf("success")
}
