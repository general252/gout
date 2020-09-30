package uencode

import (
	"encoding/base64"
	"testing"
)

func TestAESEncrypt(t *testing.T) {
	d := make([]byte, 128)
	for i := 0; i < len(d); i++ {
		d[i] = byte(i%26) + 'a'
	}
	key := []byte("hgfedcba87654321")
	iv :=  []byte("abcdefkxlikalkgi")
	x1, err := AESEncrypt(d, key, iv)
	if err != nil {
		t.Error(err)
	}

	t.Logf("加密后: %v", len(x1))
	t.Logf("%v", base64.StdEncoding.EncodeToString(x1))
	x2, err := AESDecrypt(x1, key, iv)
	if err != nil {
		t.Error(err)
	}
	t.Logf("解密后: %v %v", len(x2), string(x2))
}
