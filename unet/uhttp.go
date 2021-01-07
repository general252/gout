package unet

import (
	"bytes"
	"context"
	"encoding/json"
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
