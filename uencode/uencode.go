package uencode

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
)

func UUID() string {
	return uuid.NewV4().String()
}

func MD5Bit16(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))[8:24]
}

func MD5Bit32(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func SHA1(date []byte) string {
	h := sha1.New()
	h.Write(date)
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(date []byte) string {
	h := sha256.New()
	h.Write(date)
	return hex.EncodeToString(h.Sum(nil))
}

func SHA512(date []byte) string {
	h := sha512.New()
	h.Write(date)
	return hex.EncodeToString(h.Sum(nil))
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func HmacSHA1(key []byte, data []byte) string {
	mac := hmac.New(sha1.New, key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

func HmacSHA256(key []byte, data []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

func HmacSHA512(key []byte, data []byte) string {
	mac := hmac.New(sha512.New, key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}
