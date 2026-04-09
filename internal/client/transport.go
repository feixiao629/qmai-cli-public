package client

import (
	"fmt"
	"net/http"
	"time"
)

// RetryTransport retries failed requests
type RetryTransport struct {
	Base       http.RoundTripper
	MaxRetries int
	Delay      time.Duration
}

func NewRetryTransport(base http.RoundTripper, maxRetries int) *RetryTransport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &RetryTransport{
		Base:       base,
		MaxRetries: maxRetries,
		Delay:      time.Second,
	}
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i <= t.MaxRetries; i++ {
		resp, err = t.Base.RoundTrip(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		if i < t.MaxRetries {
			time.Sleep(t.Delay * time.Duration(i+1))
		}
	}
	if err != nil {
		return nil, fmt.Errorf("request failed after %d retries: %w", t.MaxRetries, err)
	}
	return resp, nil
}

