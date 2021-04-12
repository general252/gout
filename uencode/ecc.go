package uencode

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
)

// ECC 加密: "github.com/wumansgy/goEncrypt"
// RSA 私钥加密: "github.com/wenzhenxi/gorsa"

// EccGenerateKey 生成ECC私钥和公钥. c: 推荐elliptic.P256() ECC 164位的密钥产生的一个安全级相当于RSA 1024位密钥提供的保密强度
func EccGenerateKey(c elliptic.Curve) (priKey, pubKey string, err error) {
	var outPriKey = bytes.Buffer{}
	var outPubKey = bytes.Buffer{}

	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		return "", "", err
	}

	//log.Println("N: ", privateKey.N)
	//log.Println("E: ", privateKey.E)
	//log.Println("D: ", privateKey.D)

	//保存私钥
	{
		// 通过x509标准将得到的ecc私钥序列化为ASN.1 的 DER编码字符串
		bytesX509PrivateKey, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			return "", "", err
		}
		//使用pem格式对x509输出的内容进行编码

		// 保存
		if err := pem.Encode(&outPriKey, &pem.Block{
			Type:  "ECC Private Key",
			Bytes: bytesX509PrivateKey,
		}); err != nil {
			return "", "", err
		}
	}

	//保存公钥
	{
		//获取公钥的数据
		publicKey := privateKey.PublicKey

		//X509对公钥编码
		bytesX509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
		if err != nil {
			return "", "", err
		}

		//pem格式编码
		if err := pem.Encode(&outPubKey, &pem.Block{
			Type:  "ECC Public Key",
			Bytes: bytesX509PublicKey,
		}); err != nil {
			return "", "", err
		}
	}

	return outPriKey.String(), outPubKey.String(), nil
}

// EccSign 签名
func EccSign(msg []byte, pemPriKey string) (r []byte, s []byte, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = fmt.Errorf("EccSign error: %v", er)
		}
	}()

	block, _ := pem.Decode([]byte(pemPriKey))
	if block == nil || len(block.Bytes) <= 0 {
		return nil, nil, fmt.Errorf("decode pem fail")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	hash := sha256.New()
	hash.Write(msg)
	resultHash := hash.Sum(nil)

	rt, st, err := ecdsa.Sign(rand.Reader, privateKey, resultHash)
	if err != nil {
		return nil, nil, err
	}

	r, err = rt.MarshalText()
	if err != nil {
		return nil, nil, err
	}
	s, err = st.MarshalText()
	if err != nil {
		return nil, nil, err
	}
	return r, s, nil
}

// EccVerifySign 验证签名
func EccVerifySign(msg []byte, pemPubKey string, r, s []byte) (result bool) {
	defer func() {
		if err := recover(); err != nil {
			result = false
			log.Println(err)
		}
	}()

	block, _ := pem.Decode([]byte(pemPubKey))
	if block == nil || len(block.Bytes) <= 0 {
		return false
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}

	publicKey := publicKeyInterface.(*ecdsa.PublicKey)

	hash := sha256.New()
	hash.Write(msg)
	resultHash := hash.Sum(nil)

	var rt, st big.Int
	_ = rt.UnmarshalText(r)
	_ = st.UnmarshalText(s)

	result = ecdsa.Verify(publicKey, resultHash, &rt, &st)
	return result
}
