package unet

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
	"time"
)

// HttpDo http do
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

// HttpRequestJsonWithHeader http请求, context-type: "application/json; charset=utf-8"
func HttpRequestJsonWithHeader(method string, url string, body []byte, headers map[string]string) ([]byte, error) {
	return HttpRequestWithContextType(method, url, body, "application/json; charset=utf-8", headers)
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

// HttpRequestJsonWithClient http请求, context-type: "application/json; charset=utf-8"
func HttpRequestJsonWithClient(method string, url string, body []byte, headers map[string]string, cli *http.Client) ([]byte, error) {
	var p io.Reader
	if body != nil {
		p = bytes.NewBuffer(body)
	}
	req, err := http.NewRequest(method, url, p)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

type HttpRequestParam struct {
	URL     string
	Method  string // GET/POST
	Body    []byte
	Headers map[string]string
	cli     *http.Client
}

func NewHttpRequestParam() *HttpRequestParam {
	return &HttpRequestParam{
		URL:    "https://www.baidu.com/s?ie=utf-8&wd=hello",
		Method: http.MethodPost,
		Body:   nil,
		Headers: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
		cli: http.DefaultClient,
	}
}

// HttpRequestCustom http请求, context-type: "application/json; charset=utf-8"
func HttpRequestCustom(param *HttpRequestParam, fn func(req *http.Request, param *HttpRequestParam)) ([]byte, error) {
	if param == nil {
		return nil, fmt.Errorf("invate param")
	}

	var p io.Reader
	if param.Body != nil {
		p = bytes.NewBuffer(param.Body)
	}
	req, err := http.NewRequest(param.Method, param.URL, p)
	if err != nil {
		return nil, err
	}

	for key, value := range param.Headers {
		req.Header.Set(key, value)
	}

	if fn != nil {
		fn(req, param)
	}

	cli := param.cli
	if cli == nil {
		cli = http.DefaultClient
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// NewMultipartFormData multipart/form-data
func NewMultipartFormData(uri string, params map[string]string, fieldName, filename string, fileData *bytes.Buffer) (*http.Request, error) {
	var (
		err  error
		body = &bytes.Buffer{}

		writer *multipart.Writer
		part   io.Writer

		request *http.Request
	)

	// multipart
	writer = multipart.NewWriter(body)
	for k, v := range params {
		_ = writer.WriteField(k, v)
	}

	if part, err = writer.CreateFormFile(fieldName, filename); err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, fileData); err != nil {
		return nil, err
	}

	_ = writer.Close()

	if request, err = http.NewRequest(http.MethodPost, uri, body); err != nil {
		return nil, err
	} else {
		request.Header.Set("Content-Type", writer.FormDataContentType())
	}

	return request, nil
}
