// Package httpclient implements HTTP client.
package httpclient

import (
	"bytes"
	"context"
	"net/http"
	"time"
)

//go:generate mockgen -source=client.go -destination=../mocks/client_mocks.go -package=mocks

// InterfaceClient - for mock.
type InterfaceClient interface {
	Get(url string) (*http.Response, error)
	Post(url string, body []byte) (*http.Response, error)
}

// Client - simple web client.
type Client struct {
	client *http.Client
}

const (
	_defaultTimeout = 10 * time.Second
)

// NewClient - init new Client.
func NewClient(opts ...Option) *Client {
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

// Get - GET request with timeout.
func (s *Client) Get(url string) (*http.Response, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	return s.client.Do(req)
}

// Post - POST request with timeout.
func (s *Client) Post(url string, body []byte) (*http.Response, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return s.client.Do(req)
}
