package uproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
}

type ProxyOption struct {
	Director       func(*http.Request)
	ModifyResponse func(*http.Response) error
	ErrorHandler   func(http.ResponseWriter, *http.Request, error)
}

// SimpleReverseProxy targetURL: "http://127.0.0.1:5567"
//	SimpleReverseProxy("http://127.0.0.1:5567", writer, request, &ProxyOption{
//		Director: func(request *http.Request) {
//			request.URL.Path = "/target/path/" + uid
//		},
//	})
func SimpleReverseProxy(targetURL string, writer http.ResponseWriter, request *http.Request, option *ProxyOption) {
	target, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	original := ProxyOption{
		Director:       proxy.Director,
		ModifyResponse: proxy.ModifyResponse,
		ErrorHandler:   proxy.ErrorHandler,
	}

	// 修改请求头
	proxy.Director = func(request *http.Request) {
		if original.Director != nil {
			original.Director(request)
		}

		if option != nil && option.Director != nil {
			option.Director(request)
		}
	}

	proxy.ModifyResponse = func(response *http.Response) error {
		if original.ModifyResponse != nil {
			if err := original.ModifyResponse(response); err != nil {
				return err
			}
		}

		if option != nil && option.ModifyResponse != nil {
			return option.ModifyResponse(response)
		}

		return nil
	}

	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		if original.ErrorHandler != nil {
			original.ErrorHandler(writer, request, err)
		}

		if option != nil && option.ErrorHandler != nil {
			option.ErrorHandler(writer, request, err)
		}
	}

	proxy.ServeHTTP(writer, request)
}
