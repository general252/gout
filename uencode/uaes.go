package uencode

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
)

// 填充数据
func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

// 去掉填充数据
func unPadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

// AESEncrypt 加密
func AESEncrypt(src []byte, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	src = padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(src, src)

	return src, nil
}

// AESDecrypt 解密
func AESDecrypt(src []byte, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(src, src)
	src = unPadding(src)

	return src, nil
}

func getKeyIv(password string) (key, iv []byte) {
	h256 := sha256.New()
	h256.Write([]byte("salt_537B" + password))
	key = h256.Sum(nil)

	h := sha1.New()
	h.Write([]byte("salt_4CD1" + password))
	iv = h.Sum(nil)

	iv = iv[:16]

	return
}

// AESEncryptV2 加密
func AESEncryptV2(data []byte, password string) ([]byte, error) {
	key, iv := getKeyIv(password)
	return AESEncrypt(data, key, iv)
}

// AESDecryptV2 解码
func AESDecryptV2(data []byte, password string) ([]byte, error) {
	key, iv := getKeyIv(password)
	return AESDecrypt(data, key, iv)
}
