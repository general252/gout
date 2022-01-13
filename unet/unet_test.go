package unet

// go test . -v  -count=1

import "testing"

func TestGetHostIP(t *testing.T) {
	ip, err := GetHostIP()
	if err != nil {
		t.Errorf("GetHostIP fail %v", err)
	} else {
		t.Logf("GetHostIP success %v", ip)
	}
}
