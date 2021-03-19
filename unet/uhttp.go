package unet

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// HttpDo
func HttpDo(method, url, data string, headers map[string]string) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				conn, err := net.DialTimeout(network, addr, time.Second*5) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				/*err = conn.SetDeadline(time.Now().Add(time.Second * 15)) //设置发送接受数据超时
				if err != nil {
					return nil, err
				}*/
				return conn, nil
			},
			//ResponseHeaderTimeout: time.Second * 2,
		},
	}

	req, err := http.NewRequest(
		method,
		url,
		strings.NewReader(data))
	if err != nil {
		// handle error
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// HttpRequestJson http请求, context-type: "application/json; charset=utf-8"
func HttpRequestJson(method string, url string, body []byte) ([]byte, error) {
	return HttpRequestWithContextType(method, url, body, "application/json; charset=utf-8", nil)
}

// HttpRequestMultipartFormData http请求, context-type: "multipart/form-data"
func HttpRequestMultipartFormData(method string, url string, body []byte) ([]byte, error) {
	return HttpRequestWithContextType(method, url, body, "multipart/form-data", nil)
}

// HttpRequestWithoutContextType http请求, no context-type
func HttpRequestWithoutContextType(method string, url string, body []byte) ([]byte, error) {
	return HttpRequestWithContextType(method, url, body, "", nil)
}

// HttpRequestWithContextType http请求
func HttpRequestWithContextType(method string, url string, body []byte, contextType string, headers map[string]string) ([]byte, error) {
	var p io.Reader
	if body != nil {
		p = bytes.NewBuffer(body)
	}
	req, err := http.NewRequest(method, url, p)
	if err != nil {
		return nil, err
	}

	if len(contextType) > 0 {
		req.Header.Add("Content-Type", contextType)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				conn, err := net.DialTimeout(network, addr, time.Second*5) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
		},
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func FormatJson(jsonData []byte) []byte {
	var obj interface{}
	if err := json.Unmarshal(jsonData, &obj); err != nil {
		return nil
	}

	if data, err := json.MarshalIndent(obj, "", "  "); err != nil {
		return nil
	} else {
		return data
	}
}

func FormatJsonObject(v interface{}) []byte {
	if data, err := json.MarshalIndent(v, "", "  "); err != nil {
		return nil
	} else {
		return data
	}
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
