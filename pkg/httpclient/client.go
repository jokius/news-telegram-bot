// Package httpclient implements HTTP client.
package httpclient

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultTimeout = 10 * time.Second
)

// Client -.
type Client struct {
	client *http.Client
}

// New -.
func New(handler http.Handler, opts ...Option) *Client {
	httpClient := &http.Client{
		Timeout: _defaultTimeout,
	}

	s := &Client{
		client: httpClient,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Client) Get(url string) (*http.Response, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	return s.client.Do(req)
}
