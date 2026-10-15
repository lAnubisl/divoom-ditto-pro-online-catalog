package httpx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NewClient creates an HTTP client that rejects redirects to protect credentials.
func NewClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
}

// Fetch executes a bounded HTTP request without exposing URLs or response bodies in errors.
func Fetch(ctx context.Context, client *http.Client, method, target string, body io.Reader, headers map[string]string, limit int64) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, target, body)
	if err != nil {
		return nil, errors.New("invalid HTTP request")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		// net/http errors contain URLs; never log them or remote response bodies.
		return nil, errors.New("HTTP network error or timeout")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, errors.New("cannot read HTTP response")
	}
	if int64(len(data)) > limit {
		return nil, errors.New("HTTP response exceeds size limit")
	}
	return data, nil
}
