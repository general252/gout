package uencode

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"regexp"

	"github.com/google/uuid"
)

func UUID() string {
	return uuid.New().String()
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

// IsInvalidBase64 检查数据是无效的base64数据
// true: data不是base64数据
// false: data符合base64(注意: 不一定是base64数据, 也可能是hex数据)
func IsInvalidBase64(str string) bool {
	pattern := "^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{4}|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)$"
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return true
	}

	if !(len(str)%4 == 0 && matched) {
		return true
	}

	_, err = base64.StdEncoding.DecodeString(str)
	if err != nil {
		return true
	}
	return false
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

func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
