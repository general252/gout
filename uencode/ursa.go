package uencode

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// RSAGenerateKey 生成RSA私钥和公钥. bits 证书大小
func RSAGenerateKey(bits int, certInfo *CertInfo) (priKey, pubKey string, certPubKey string, err error) {
	var outPriKey = bytes.Buffer{}
	var outPubKey = bytes.Buffer{}
	var outPubKeyCert = bytes.Buffer{}

	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", "", err
	}

	//log.Println("N: ", privateKey.N)
	//log.Println("E: ", privateKey.E)
	//log.Println("D: ", privateKey.D)

	//保存私钥
	{
		// 通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
		bytesX509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
		//使用pem格式对x509输出的内容进行编码

		// 保存
		if err := pem.Encode(&outPriKey, &pem.Block{
			Type:  "RSA Private Key",
			Bytes: bytesX509PrivateKey,
		}); err != nil {
			return "", "", "", err
		}
	}

	//保存公钥
	{
		//获取公钥的数据
		publicKey := privateKey.PublicKey

		//X509对公钥编码
		bytesX509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
		if err != nil {
			return "", "", "", err
		}

		//pem格式编码
		if err := pem.Encode(&outPubKey, &pem.Block{
			Type:  "RSA Public Key",
			Bytes: bytesX509PublicKey,
		}); err != nil {
			return "", "", "", err
		}
	}

	//保存公钥cert
	if certInfo != nil {
		data, err := getCert(privateKey, certInfo)
		if err != nil {
			return "", "", "", err
		}

		outPubKeyCert.WriteString(data)
	}

	return outPriKey.String(), outPubKey.String(), outPubKeyCert.String(), nil
}

// RSAEncrypt RSA加密
// plainText 要加密的数据
// pemPubKey pem格式公钥
// 一个分组的大小 key size - 42. 例如bits: 384, 分组大小: 384/8-42=6
func RSAEncrypt(plainText []byte, pemPubKey string) (rText []byte, err error) {
	defer func() {
		if er := recover(); er != nil {
			rText = nil
			err = fmt.Errorf("%v", er)
		}
	}()

	blk, _ := pem.Decode([]byte(pemPubKey))
	if blk == nil || len(blk.Bytes) <= 0 {
		return nil, fmt.Errorf("decode pem fail")
	}

	objPubKeyInterface, err := x509.ParsePKIXPublicKey(blk.Bytes)
	if err != nil {
		return nil, err
	}

	objPubKey, ok := objPubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("get public key fail")
	}

	var hash = sha1.New() // hash.Size() is 20

	var oneMaxSize = objPubKey.Size() - 2*hash.Size() - 2
	if oneMaxSize <= 0 {
		return nil, fmt.Errorf(`public key length to small. key size: %v, bits > (20*2+2)x8. key_size > 2*hash_size+2`,
			objPubKey.Size())
	}

	var plainTextLength = len(plainText)
	var count = plainTextLength / oneMaxSize
	if plainTextLength%oneMaxSize > 0 {
		count += 1
	}

	var result bytes.Buffer
	for i := 0; i < count; i++ {
		var offset = i * oneMaxSize
		var end = offset + oneMaxSize
		if end > plainTextLength {
			end = plainTextLength
		}
		srcData := plainText[offset:end]

		if data, err := rsa.EncryptOAEP(hash, rand.Reader, objPubKey, srcData, nil); err != nil {
			// len(srcData) max value is: objPubKey.Size()-2*sha1.New().Size()-2
			// data, err := rsa.EncryptOAEP(hash, rand.Reader, objPubKey, srcData, nil)

			// len(srcData) max value is: objPubKey.Size()-11
			// data, err = rsa.EncryptPKCS1v15(rand.Reader, objPubKey, srcData)
			return nil, err
		} else {
			result.Write(data)
		}
	}

	return result.Bytes(), nil
}

// RSADecrypt RSA解密
// cipherText 需要解密的byte数据
// pemPriKey pem格式私钥
func RSADecrypt(cipherText []byte, pemPriKey string) (rText []byte, err error) {
	defer func() {
		if er := recover(); er != nil {
			rText = nil
			err = fmt.Errorf("%v", er)
		}
	}()

	// pem解码
	block, _ := pem.Decode([]byte(pemPriKey))
	if block == nil || len(block.Bytes) <= 0 {
		return nil, fmt.Errorf("decode pem fail")
	}

	// X509解码
	objPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	var hash = sha1.New() // hash.Size() is 20

	var plainTextLength = len(cipherText)
	var oneMaxSize = objPrivateKey.Size()
	if plainTextLength%oneMaxSize != 0 {
		return nil, fmt.Errorf("plainText length error, %v", oneMaxSize)
	}
	if plainTextLength < oneMaxSize {
		return nil, fmt.Errorf("plainText length error, < %v", oneMaxSize)
	}

	var count = plainTextLength / oneMaxSize
	var result bytes.Buffer
	for i := 0; i < count; i++ {
		var offset = i * oneMaxSize
		var end = offset + oneMaxSize

		srcData := cipherText[offset:end]
		if data, err := rsa.DecryptOAEP(hash, rand.Reader, objPrivateKey, srcData, nil); err != nil {
			//plainText, err = rsa.DecryptPKCS1v15(rand.Reader, objPrivateKey, cipherText)
			return nil, err
		} else {
			result.Write(data)
		}
	}

	// 返回明文
	return result.Bytes(), nil
}

// RSAValidatePemPublicKey 检查pem格式公钥
func RSAValidatePemPublicKey(pemPubKey string) bool {
	blk, _ := pem.Decode([]byte(pemPubKey))
	if blk == nil || len(blk.Bytes) <= 0 {
		return false
	}

	objPubKeyInterface, err := x509.ParsePKIXPublicKey(blk.Bytes)
	if err != nil {
		return false
	}

	objPublicKey, ok := objPubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return false
	}

	if objPublicKey.N == nil {
		return false
	}

	if objPublicKey.E < 2 {
		return false
	}
	if objPublicKey.E > 1<<31-1 {
		return false
	}
	return true
}

// RSAValidatePemPrivateKey 检查pem格式私钥
func RSAValidatePemPrivateKey(pemPriKey string) bool {
	// pem解码
	block, _ := pem.Decode([]byte(pemPriKey))
	if block == nil || len(block.Bytes) <= 0 {
		return false
	}

	// X509解码
	objPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return false
	}

	return objPrivateKey.Validate() == nil
}

// RSASign 签名
func RSASign(msg []byte, pemPriKey string) (rText []byte, err error) {
	defer func() {
		if er := recover(); er != nil {
			rText = nil
			err = fmt.Errorf("%v", er)
		}
	}()

	block, _ := pem.Decode([]byte(pemPriKey))
	if block == nil || len(block.Bytes) <= 0 {
		return nil, fmt.Errorf("decode pem fail")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	hash := sha256.New()
	hash.Write(msg)
	hashed := hash.Sum(nil)

	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return nil, err
	}

	return sign, nil
}

// RSAVerifySign 验证签名
func RSAVerifySign(msg []byte, sign []byte, pemPubKey string) (rc bool) {
	defer func() {
		if er := recover(); er != nil {
			rc = false
		}
	}()

	block, _ := pem.Decode([]byte(pemPubKey))
	if block == nil || len(block.Bytes) <= 0 {
		return false
	}

	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}

	publicKey := publicInterface.(*rsa.PublicKey)

	hash := sha256.New()
	hash.Write(msg)
	hashed := hash.Sum(nil)

	result := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed, sign)
	return result == nil
}
