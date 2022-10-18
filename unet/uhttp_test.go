package unet

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"path/filepath"
	"time"
)

var (
	BaseURL = "http://127.0.0.1/api/test"
)

func ExampleHttpDo() {
	data, _ := HttpDo(
		"POST",
		fmt.Sprintf("%v/authorization", BaseURL),
		"{\"a\":1,\"b\":\"22\"}",
		map[string]string{
			"Content-Type":  "application/json; charset=utf-8",
			"Authorization": "abc|20C2FB99644A789BE386722EEFF37AE7D062866B",
			"Content-MD5":   "MzExNUJDMTkzMjIxMTJCMw==",
			"Date":          "Tue, 15 Nov 1994 08:12:31 GMT",
		})

	fmt.Printf("%v\n", string(data))
	// output:
	//
}

func ExampleHttpRequestJson() {
	url := fmt.Sprintf("%v/file/%v", BaseURL, "fileId")
	resp, err := HttpRequestJson(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(FormatJson(resp)))
	// output:
	//
}

// 上传文件
func ExampleHttpRequestWithContextType() {
	url := fmt.Sprintf("%v/file", BaseURL)
	path := `/your/file/path`

	body, contextType, err := getPostFileBody("uploadFile", path)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := HttpRequestWithContextType(http.MethodPost, url, body, contextType, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(FormatJson(resp)))

	// output:
	//

}

// getPostFileBody 获取上传文件的body, context-type
func getPostFileBody(uploadFileKey string, fileFullPath string) ([]byte, string, error) {
	fileData, err := ioutil.ReadFile(fileFullPath)
	if err != nil {
		return nil, "", err
	}

	var resp bytes.Buffer

	// 文件写入 body
	writer := multipart.NewWriter(&resp)
	defer writer.Close()

	part, err := writer.CreateFormFile(uploadFileKey, filepath.Base(fileFullPath))
	if err != nil {
		return nil, "", err
	}

	_, _ = part.Write(fileData)

	return resp.Bytes(), writer.FormDataContentType(), nil
}

func ExampleHttpRequestJsonWithHeader() {
	resp, err := HttpRequestJsonWithHeader(http.MethodPost, "url", nil, map[string]string{
		"Authorization": "abc",
	})
	_, _ = resp, err

	// output:
	//
}

func ExampleHttpRequestJsonWithClient() {
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

	for i := 0; i < 100; i++ {
		resp, err := HttpRequestJsonWithClient(http.MethodPost, "url", nil, nil, cli)
		if err != nil {
			return
		}
		log.Println(FormatJsonString(resp))
	}

	// output:

}

func ExampleHttpRequestCustom() {
	param := NewHttpRequestParam()
	body, err := HttpRequestCustom(param, func(req *http.Request, param *HttpRequestParam) {

		param.cli = &http.Client{
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
	})

	_, _ = body, err

	// output:
}
