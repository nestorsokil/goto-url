package service

import (
	"strconv"

	"github.com/nestorsokil/goto-url/util"
)

// Request is a struct that carries request info
type Request struct {
	url       string
	expire    int64
	customKey string
}

// RequestBuilder is a convenient builder that can be used to set params on a request
type RequestBuilder struct {
	request Request
}

func builder(conf *util.ApplicationConfig) *RequestBuilder {
	req := Request{expire: conf.ExpirationTimeHours}
	return &RequestBuilder{request: req}
}

// ForURL sets the URL for the record to hold
func (rb *RequestBuilder) ForURL(url string) *RequestBuilder {
	rb.request.url = url
	return rb
}

// WithCustomKey sets a custom suffix for the short URL
func (rb *RequestBuilder) WithCustomKey(key string) *RequestBuilder {
	rb.request.customKey = key
	return rb
}

// WithCustomExpirationTime overrides the default expiration time for a record
func (rb *RequestBuilder) WithCustomExpirationTime(hours string) *RequestBuilder {
	if hours != "" {
		v, e := strconv.Atoi(hours)
		if e == nil {
			rb.request.expire = int64(v)
		}
	}
	return rb
}

// Build returns a record (terminal op)
func (rb *RequestBuilder) Build() Request {
	return rb.request
}
