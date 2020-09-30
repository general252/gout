package unet

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

/** http 请求
data, _ := HttpDo(
	"POST",
	"http://127.0.0.1/api/test/authorization",
	"{\"a\":1,\"b\":\"22\"}",
	map[string]string{
		"Content-Type":  "application/json; charset=utf-8",
		"Authorization": "abc|20C2FB99644A789BE386722EEFF37AE7D062866B",
		"Content-MD5":   "MzExNUJDMTkzMjIxMTJCMw==",
		"Date":          "Tue, 15 Nov 1994 08:12:31 GMT",
	})

fmt.Printf("%v\n", string(data))
*/
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
