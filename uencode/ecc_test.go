package uencode

import (
	"crypto/elliptic"
	"fmt"
	"log"
	"strings"
)

func ExampleEccGenerateKey() {
	pri, pub, _ := EccGenerateKey(elliptic.P256())
	fmt.Println(pri)
	fmt.Println(pub)

	// output:
	//-----BEGIN ECC Private Key-----
	//MHcCAQEEIKxexVTP94xVS2LFoo1BLZwgw258ihSoK3JNbcqr/fMBoAoGCCqGSM49
	//AwEHoUQDQgAEu1Bk0NxQrfLI6yGaY2ZvXum1I6r1DVHEy6xH2Np/YHudWDo+8+gS
	//LCgSpfktoCZUzlzjMQGPVv7IbrEFAMnoDg==
	//-----END ECC Private Key-----
	//
	//-----BEGIN ECC Public Key-----
	//MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEu1Bk0NxQrfLI6yGaY2ZvXum1I6r1
	//DVHEy6xH2Np/YHudWDo+8+gSLCgSpfktoCZUzlzjMQGPVv7IbrEFAMnoDg==
	//-----END ECC Public Key-----
}

var (
	pri = `-----BEGIN ECC Private Key-----
MHcCAQEEINE3HKnaHF1h1flTRgpy2nXOosnJwIVPk2fZlmFyrz3boAoGCCqGSM49
AwEHoUQDQgAEe4Utbc0fIkywI+dvnTM+nOcG/gWe5MIvu0Eh63pmDIGGhcyFRab4
2/irUGw2fUoBycpGvLdav0ftaKmWhgisPQ==
-----END ECC Private Key-----`
	pub = `-----BEGIN ECC Public Key-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEe4Utbc0fIkywI+dvnTM+nOcG/gWe
5MIvu0Eh63pmDIGGhcyFRab42/irUGw2fUoBycpGvLdav0ftaKmWhgisPQ==
-----END ECC Public Key-----`
)

func ExampleEccSign() {
	var msg = []byte(strings.Repeat("a", 1000))
	r, s, err := EccSign(msg, pri)
	if err != nil {
		return
	}

	log.Println(string(r))
	log.Println(string(s))

	r = []byte("53239160152256391266514960151342461391297860554311123424745057892820584220996")
	s = []byte("90267341205982732265789688485691605571964324220642226112561266862518196681918")

	v := EccVerifySign(msg, pub, r, s)
	log.Println(v)

	// output:
}
