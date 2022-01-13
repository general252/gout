package uencode

import (
	"log"
	"time"
)

func ExampleCert2Pem() {
	pemPri, pemPub, certPub, err := RSAGenerateKey(1024, &CertInfo{
		CommonName: "golang",
		NotBefore:  time.Now(),
		NotAfter:   time.Now().AddDate(100, 0, 0),
	})
	if err != nil {
		return
	}

	log.Println(pemPri)
	log.Println(certPub)
	log.Println(pemPub)

	tmpData, _ := Cert2Pem(certPub, "RSA Public Key")
	log.Println(tmpData)

	// output:

}
