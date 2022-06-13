package unet

import (
	"net/http"
	"net/url"
	"strconv"
)

type uHttpRequest struct {
	request    *http.Request
	queryCache url.Values
}

// NewUHttpRequest 封装http url查询
func NewUHttpRequest(request *http.Request) *uHttpRequest {
	return &uHttpRequest{request: request}
}

func (c *uHttpRequest) getQueryCache() bool {
	if c.queryCache == nil {
		c.queryCache = c.request.URL.Query()
	}

	return c.queryCache != nil
}

func (c *uHttpRequest) GetString(key string) (string, bool) {
	if !c.getQueryCache() {
		return "", false
	}

	val := c.queryCache[key]
	if len(val) == 0 {
		return "", false
	}

	return val[0], true
}

func (c *uHttpRequest) GetInt64(key string) (int64, bool) {
	val, ok := c.GetString(key)
	if !ok {
		return 0, false
	}

	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, false
	}

	return v, true
}

func (c *uHttpRequest) GetInt(key string) (int, bool) {
	v, ok := c.GetInt64(key)

	return int(v), ok
}

func (c *uHttpRequest) GetBool(key string) (bool, bool) {
	val, ok := c.GetString(key)
	if !ok {
		return false, false
	}

	v, err := strconv.ParseBool(val)
	if err != nil {
		return false, false
	}

	return v, true
}

func (c *uHttpRequest) GetFloat(key string) (float64, bool) {
	val, ok := c.GetString(key)
	if !ok {
		return 0, false
	}

	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, false
	}

	return v, true
}
