package uencode

import (
	"encoding/base64"
	"testing"
)

func TestGenerateRSAKey(t *testing.T) {
	priKey, pubKey, err := RSAGenerateKey(384)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(RSAValidatePemPrivateKey(priKey))
	t.Log(RSAValidatePemPublicKey(pubKey))

	var src []byte
	for i := 0; i < 58; i++ {
		var x = '0' + i%10
		src = append(src, byte(x))
	}

	temp, err := RSAEncrypt(src, pubKey)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(src))
	t.Log(len(temp))
	t.Log(base64.StdEncoding.EncodeToString(temp))

	dst, err := RSADecrypt(temp, priKey)
	t.Log(string(dst))
}
