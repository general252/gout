package uencode

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
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

// 加密
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

// 解密
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
