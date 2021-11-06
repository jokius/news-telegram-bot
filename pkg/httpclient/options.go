package httpclient

import (
	"time"
)

// Option -.
type Option func(*Client)

// Timeout -.
func Timeout(timeout time.Duration) Option {
	return func(s *Client) {
		s.client.Timeout = timeout
	}
}
