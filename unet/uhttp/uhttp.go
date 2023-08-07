package uhttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/general252/gout/uoption"
)

type Method string

const (
	MethodGet    Method = http.MethodGet
	MethodHead   Method = http.MethodHead
	MethodPost   Method = http.MethodPost
	MethodPut    Method = http.MethodPut
	MethodPatch  Method = http.MethodPatch
	MethodDelete Method = http.MethodDelete
)

func HttpRequestObject[T any](url string, options ...uoption.Option[*requestOption]) (T, error) {
	return newObjectHttpRequest[T]().getObject(url, options...)
}

func HttpRequestBytes(url string, options ...uoption.Option[*requestOption]) ([]byte, error) {
	return newObjectHttpRequest[int]().getBytes(url, options...)
}

type requestOption struct {
	ctx             context.Context
	method          Method
	contextType     string
	body            io.Reader
	headers         map[string]string
	client          *http.Client
	handle          func(cli *http.Client, req *http.Request) (*http.Response, error)
	handleReply     func(reply *http.Response) error
	handleReplyBody func(body []byte) error
}

type objectHttpRequest[T any] struct {
	requestOption
}

func newObjectHttpRequest[T any]() *objectHttpRequest[T] {
	return &objectHttpRequest[T]{
		requestOption{
			method:      MethodGet,
			contextType: "application/json; charset=utf-8",
			body:        nil,
			headers:     map[string]string{},
			client: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			},
			handle: func(cli *http.Client, req *http.Request) (*http.Response, error) {
				return cli.Do(req)
			},
			handleReply: func(response *http.Response) error {
				return nil
			},
			handleReplyBody: func(body []byte) error {
				return nil
			},
		},
	}
}

func (tis *objectHttpRequest[T]) getBytes(url string, options ...uoption.Option[*requestOption]) ([]byte, error) {

	for _, option := range options {
		option.Apply(&tis.requestOption)
	}

	if tis.ctx == nil {
		tis.ctx = context.TODO()
	}

	request, err := http.NewRequestWithContext(tis.ctx, string(tis.method), url, tis.body)
	if err != nil {
		return nil, err
	}

	for k, v := range tis.headers {
		request.Header.Set(k, v)
	}

	resp, err := tis.handle(tis.client, request)
	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
	}

	if err = tis.handleReply(resp); err != nil {
		return nil, err
	}

	// https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
	//Note
	//	1xx: Informational - Request received, continuing process
	//	2xx: Success - The action was successfully received, understood, and accepted
	//	3xx: Redirection - Further action must be taken in order to complete the request
	//	4xx: Client Error - The request contains bad syntax or cannot be fulfilled
	//	5xx: Server Error - The server failed to fulfill an apparently valid request
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%v", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = tis.handleReplyBody(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (tis *objectHttpRequest[T]) getObject(url string, options ...uoption.Option[*requestOption]) (value T, err error) {

	data, err := tis.getBytes(url, options...)
	if err != nil {
		return value, err
	}

	if err = json.Unmarshal(data, &value); err != nil {
		return value, err
	}

	return value, nil
}

func WithMethod(m Method) uoption.Option[*requestOption] {
	return uoption.NewFuncOption[*requestOption](func(o *requestOption) {
		o.method = m
	})
}

func WithBody(body io.Reader) uoption.Option[*requestOption] {
	return uoption.NewFuncOption[*requestOption](func(o *requestOption) {
		o.body = body
	})
}

func WithContext(ctx context.Context) uoption.Option[*requestOption] {
	return uoption.NewFuncOption[*requestOption](func(o *requestOption) {
		o.ctx = ctx
	})
}

func WithBufferBody(body []byte) uoption.Option[*requestOption] {
	return uoption.NewFuncOption[*requestOption](func(o *requestOption) {
		o.body = bytes.NewBuffer(body)
	})
}

func WithObjectBody(object any) uoption.Option[*requestOption] {
	return uoption.NewFuncOption[*requestOption](func(o *requestOption) {
		if data, err := json.Marshal(object); err == nil {
			o.body = bytes.NewBuffer(data)
		}
	})
}

func WithContextType(contextType string) uoption.Option[*requestOption] {
	return uoption.NewFuncOption[*requestOption](func(o *requestOption) {
		o.contextType = contextType
	})
}

func WithHeaders(headers map[string]string) uoption.Option[*requestOption] {
	return uoption.NewFuncOption[*requestOption](func(o *requestOption) {
		for k, v := range headers {
			o.headers[k] = v
		}
	})
}

func WithHeaderAuthorization(token string) uoption.Option[*requestOption] {
	return uoption.NewFuncOption[*requestOption](func(o *requestOption) {
		o.headers["Authorization"] = token
	})
}

func WithClient(client *http.Client) uoption.Option[*requestOption] {
	return uoption.NewFuncOption[*requestOption](func(o *requestOption) {
		o.client = client
	})
}

func WithDo(handle func(cli *http.Client, req *http.Request) (*http.Response, error)) uoption.Option[*requestOption] {
	return uoption.NewFuncOption[*requestOption](func(o *requestOption) {
		o.handle = handle
	})
}

func WithHandleReply(handle func(reply *http.Response) error) uoption.Option[*requestOption] {
	return uoption.NewFuncOption[*requestOption](func(o *requestOption) {
		o.handleReply = handle
	})
}

func WithHandleReplyBody(handle func(body []byte) error) uoption.Option[*requestOption] {
	return uoption.NewFuncOption[*requestOption](func(o *requestOption) {
		o.handleReplyBody = handle
	})
}
