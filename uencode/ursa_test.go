package uencode

import (
	"encoding/base64"
	"log"
	"strings"
	"testing"
)

func TestGenerateRSAKey(t *testing.T) {
	priKey, pubKey, _, err := RSAGenerateKey(384, nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(RSAValidatePemPrivateKey(priKey))
	t.Log(RSAValidatePemPublicKey(pubKey))

	var src = []byte(strings.Repeat("a", 10240))

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

func ExampleRSAGenerateKey() {
	pri, pub, cert, _ := RSAGenerateKey(1024, nil)
	log.Println(pri)
	log.Println(pub)
	_ = cert

	// output:

}

var (
	priKey = `-----BEGIN RSA Private Key-----
MIICXQIBAAKBgQDnkRyIA9sIuBazqKtdMFy2IMiNmTF6DHUhmnt1qXHPKeEBDRVg
QTBPX1N3X6MJ4GgzVg2g3GFIcsz+mO08AWGcPGuM/avLuquATjMak5iYVHDY/BDu
2iX8oYJPc3r1PVR1fbb3gvJbiluPfKbycRYqxa3zWqebCnFW1VNG51I6XQIDAQAB
AoGBANYzRt8SIfQoxOcfKJSk/b2DCcDhagDpsReKXJV0TdBBft6ICbPl2Rgyp3SO
xLOIHxsNiMG52Us41iLTtu6AoJYwP9J9nEk/4xr15Yfd/D09hewGuMgFxX7nKZjE
xXpQEgrpaHAP+fqFjx9Tpv+LPNXaDbvimRGt6qZ8MzcWeMXBAkEA+bLRg5aErXmq
N0b6kUd82GjOLEhyGAHwQmm1rW3rlLNvCWeD61rl2btZdNB+qVzzB1KFL6NMcum9
SUWGoEbZaQJBAO1pJzA0YTUXTfTMobirT8GlIrgMDxIZWmY0kKAT4fcoBt7ZIN9Y
yBBGmPyNgmyVYLGDS3AF8+fJ3RntWyhT5tUCQA171YqIj0Oa5VE02QUNWjWJe1Cy
3M5lFGdRtAjYfbc69U0JtPr5np3iWxNOyvg0V79WenC3HcK60ojpYzq2eLkCQHY1
EB0RR4E+vELyDGe9bHW3ekT3RB233+nZrFT38V+1X05f/90VAHASJqRA9TqJWd6o
x8vcOugi+2KoauX2eI0CQQCA05RUGvQgxKwPFlnw0A5/U24vKos2c0CZTLeezG3o
FpoYzaJAV4wonqB8e+EMkeCo/1497nTcdFFNHHXjNtQD
-----END RSA Private Key-----`

	pubKey = `-----BEGIN RSA Public Key-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDnkRyIA9sIuBazqKtdMFy2IMiN
mTF6DHUhmnt1qXHPKeEBDRVgQTBPX1N3X6MJ4GgzVg2g3GFIcsz+mO08AWGcPGuM
/avLuquATjMak5iYVHDY/BDu2iX8oYJPc3r1PVR1fbb3gvJbiluPfKbycRYqxa3z
WqebCnFW1VNG51I6XQIDAQAB
-----END RSA Public Key-----`
)

func ExampleRSAVerifySign() {
	var msg = []byte(strings.Repeat("v", 10000))
	rSign, _ := RSASign(msg, priKey)
	log.Println(rSign)

	var v = RSAVerifySign(msg, rSign, pubKey)
	log.Println(v)

	// output:

}
