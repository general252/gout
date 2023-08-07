package uhttp

import (
	"log"
	"net/http"
)

func ExampleHttpRequest() {
	type JsonTimestamp struct {
		Api  string   `json:"api"`
		V    string   `json:"v"`
		Ret  []string `json:"ret"`
		Data struct {
			T string `json:"t"`
		} `json:"data"`
	}

	object, err := HttpRequestObject[*JsonTimestamp]("http://api.m.taobao.com/rest/api3.do?api=mtop.common.getTimestamp",
		WithMethod(MethodGet),
		WithContextType("application/json; charset=utf-8"),
		WithDo(func(cli *http.Client, req *http.Request) (*http.Response, error) {
			resp, err := cli.Do(req)
			return resp, err
		}),
		WithHandleReply(func(reply *http.Response) error {
			log.Println(reply.Status)
			log.Println(reply.Header)
			log.Println(reply.Cookies())
			return nil
		}),
	)
	if err != nil {
		return
	} else {
		_ = object.Api
	}

	// output:
	//
}
