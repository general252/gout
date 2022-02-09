package uencode

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

func publicKey(priKey interface{}) interface{} {
	switch k := priKey.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func getCert(priKey interface{}, certInfo *CertInfo) (certPubKey string, err error) {
	defer func() {
		if er := recover(); er != nil {
			certPubKey = ""
			err = fmt.Errorf("%v", er)
		}
	}()

	var outPubKeyCert bytes.Buffer

	//保存私钥cert
	if certInfo == nil {
		return "", fmt.Errorf("")
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", err
	}

	keyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment

	template := x509.Certificate{
		SerialNumber: serialNumber,

		Subject: pkix.Name{CommonName: certInfo.CommonName},

		NotBefore: certInfo.NotBefore,
		NotAfter:  certInfo.NotAfter,
		IsCA:      certInfo.IsCA,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	bytesX509PrivateKey, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priKey), priKey)
	if err != nil {
		return "", err
	}

	// 保存
	if err := pem.Encode(&outPubKeyCert, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: bytesX509PrivateKey,
	}); err != nil {
		return "", err
	}

	return outPubKeyCert.String(), nil
}

// Cert2Pem cert格式转pem格式(公钥)
// pemType: "RSA Public Key"/"ECC Public Key"
func Cert2Pem(cert string, pemType string) (pemPubKey string, err error) {
	defer func() {
		if er := recover(); er != nil {
			pemPubKey = ""
			err = fmt.Errorf("%v", er)
		}
	}()

	block, _ := pem.Decode([]byte(cert))
	if block == nil || len(block.Bytes) <= 0 {
		return "", fmt.Errorf("decode pem fail")
	}

	objCertificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}

	//X509对公钥编码
	bytesX509PublicKey, err := x509.MarshalPKIXPublicKey(objCertificate.PublicKey)
	if err != nil {
		return "", err
	}

	var outPubKey bytes.Buffer
	//pem格式编码
	if err := pem.Encode(&outPubKey, &pem.Block{
		Type:  pemType,
		Bytes: bytesX509PublicKey,
	}); err != nil {
		return "", err
	}

	return outPubKey.String(), nil
}

//////////////////////////////////////////////////////////////////

type CertInfo struct {
	CommonName          string
	NotBefore, NotAfter time.Time
	IsCA                bool
}

// GenerateCert 创建Cert(issuerCert: 颁发者cert, issuerKey: 颁发者秘钥)
func GenerateCert(certInfo *CertInfo, issuerCert string, issuerKey string) (cert string, key string, err error) {
	priKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{CommonName: certInfo.CommonName},
		NotBefore:    certInfo.NotBefore,
		NotAfter:     certInfo.NotAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  certInfo.IsCA,
	}

	var objIssuer *x509.Certificate
	var objIssuerKey crypto.PrivateKey
	if len(issuerCert) == 0 {
		objIssuer = template
		objIssuerKey = priKey
	} else {
		objIssuer, err = parseCert([]byte(issuerCert))
		if err != nil {
			return "", "", err
		}
		objIssuerKey, err = parseKey([]byte(issuerKey))
		if err != nil {
			return "", "", err
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, objIssuer, priKey.Public(), objIssuerKey)
	if err != nil {
		return "", "", err
	}
	cert, err = certToString(derBytes)
	key, err = keyToString(priKey)

	return cert, key, nil
}

// CertCheckSignature 判断 childCert是parentCert颁发
func CertCheckSignature(childCert string, parentCert string) error {
	objChildCert, err := parseCert([]byte(childCert))
	if err != nil {
		return err
	}
	objParentCert, err := parseCert([]byte(parentCert))
	if err != nil {
		return err
	}

	return objChildCert.CheckSignatureFrom(objParentCert)
}

func parseCert(certBytes []byte) (*x509.Certificate, error) {
	blk, _ := pem.Decode(certBytes)
	cert, err := x509.ParseCertificate(blk.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func parseKey(keyBytes []byte) (crypto.PrivateKey, error) {
	blk, _ := pem.Decode(keyBytes)
	key, err := x509.ParseECPrivateKey(blk.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func certToString(cert []byte) (string, error) {
	var out bytes.Buffer
	if err := pem.Encode(&out, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}); err != nil {
		return "", err
	}

	return out.String(), nil
}

func keyToString(key *ecdsa.PrivateKey) (string, error) {
	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err := pem.Encode(&out, &pem.Block{
		Type:  "ECC PRIVATE KEY",
		Bytes: keyBytes,
	}); err != nil {
		return "", err
	}

	return out.String(), nil
}
