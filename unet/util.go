package unet

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// FormatJson 格式化json字符串
func FormatJson(data []byte) []byte {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil
	}

	if data, err := json.MarshalIndent(obj, "", "  "); err != nil {
		return nil
	} else {
		return data
	}
}

// FormatJsonString 格式化json字符串
func FormatJsonString(data []byte) string {
	return string(FormatJson(data))
}

// FormatJsonObject 格式化json对象
func FormatJsonObject(object interface{}) []byte {
	if data, err := json.MarshalIndent(object, "", "  "); err != nil {
		return nil
	} else {
		return data
	}
}

// FormatJsonObjectString 格式化json对象
func FormatJsonObjectString(object interface{}) string {
	return string(FormatJsonObject(object))
}

// GetRequestRealIp request real ip.
func GetRequestRealIp(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", fmt.Errorf("no valid ip found")
}

// GetRequestClientType 终端类型
func GetRequestClientType(r *http.Request) string {
	var userAgent = r.Header.Get("User-Agent")
	var cliType = ""
	if strings.Contains(userAgent, "Android") {
		cliType = "Android移动端"
		if strings.Contains(userAgent, "MicroMessenger") {
			cliType = "Android微信"
		}
	} else if strings.Contains(userAgent, "iPhone") {
		cliType = "iPhone移动客户端"
		if strings.Contains(userAgent, "MicroMessenger") {
			cliType = "iPhone微信"
		}
	} else if strings.Contains(userAgent, "iPad") {
		cliType = "iPad"
	} else if strings.Contains(userAgent, "Windows") {
		cliType = "Windows"
	} else if strings.Contains(userAgent, "Linux") {
		cliType = "Linux"
	} else {
		cliType = "UnKnow"
	}

	return cliType
}
