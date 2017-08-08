package service

import (
	"github.com/nestorsokil/goto-url/util"
	"strconv"
)

type Request struct {
	url       string
	expire    int64
	customKey string
}

type RequestBuilder struct {
	request Request
}

func builder(conf *util.Configuration) *RequestBuilder {
	req := Request{expire: conf.ExpirationTimeHours}
	return &RequestBuilder{request: req}
}

func (rb *RequestBuilder) ForUrl(url string) *RequestBuilder {
	rb.request.url = url
	return rb
}

func (rb *RequestBuilder) WithCustomKey(key string) *RequestBuilder {
	rb.request.customKey = key
	return rb
}

func (rb *RequestBuilder) WithCustomExpirationTime(hours string) *RequestBuilder {
	if hours != "" {
		v, e := strconv.Atoi(hours)
		if e == nil {
			rb.request.expire = int64(v)
		}
	}
	return rb
}

func (rb *RequestBuilder) Build() Request {
	return rb.request
}
