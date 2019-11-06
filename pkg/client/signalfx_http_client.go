package client

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

const (
	contentType = "application/json"
	timeout     = 16 * time.Second
	userAgent   = "signalview"
)

// NewConfiguredClient returns a configured client and allows the caller to
// update any client settings they desire
func NewConfiguredClient(opts ...func(*http.Client)) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			TLSHandshakeTimeout: 8 * time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			MaxIdleConns:    6,
			IdleConnTimeout: timeout,
		},
		Timeout: timeout,
	}
	// Apply any configuration options to the client
	for _, op := range opts {
		op(client)
	}
	return client
}

// NewCachedRequest caches the token and applies them to each newly created request.
func NewCachedRequest(token string) func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	// No-Op field, used just to cache token within the request generation
	return func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, err
		}
		req.Header.Add("X-SF-TOKEN", token)
		req.Header.Add("Content-Type", contentType)
		// Override the default language user agent to be clearly defined
		req.Header.Add("User-Agent", userAgent)
		return req, nil
	}
}
