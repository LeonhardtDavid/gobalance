package balancer

import (
	"context"
	"io"
	"net/http"
	"time"
)

type Client interface {
	execute(config ClientOptions) (*http.Response, error)
}

type HttpClient struct{}

type ClientOptions struct {
	timeout time.Duration
	ctx     context.Context
	method  string
	url     string
	body    io.Reader
	header  http.Header
}

func (c HttpClient) execute(options ClientOptions) (*http.Response, error) {
	client := http.Client{
		Timeout: options.timeout,
	}
	request, _ := http.NewRequest(options.method, options.url, options.body)
	request.WithContext(options.ctx)
	request.Header = options.header

	return client.Do(request)
}
