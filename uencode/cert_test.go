package uencode

import (
	"crypto/x509/pkix"
	"log"
	"time"
)

func ExampleCert2Pem() {
	pemPri, pemPub, certPub, err := RSAGenerateKey(1024, &CertInfo{
		Issuer: pkix.Name{
		},
		Subject: pkix.Name{
			Organization:       []string{"招商"},
			OrganizationalUnit: []string{"China"},
			CommonName:         "golang",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(100, 0, 0),
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
