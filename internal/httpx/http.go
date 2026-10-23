package httpx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Options configures transport and telemetry once for an upstream service.
type Options struct {
	BaseURL string
	Timeout time.Duration
	Headers map[string]string
	Service string
	Metrics Metrics
	// AcceptStatus defaults to accepting only HTTP 200.
	AcceptStatus func(int) bool
}

// Client owns HTTP execution, response limits, and request telemetry.
type Client struct {
	http    *http.Client
	options Options
}

func NewClient(options Options) *Client {
	options.Headers = copyHeaders(options.Headers)
	if options.AcceptStatus == nil {
		options.AcceptStatus = func(status int) bool { return status == http.StatusOK }
	}
	return &Client{
		http:    &http.Client{Timeout: options.Timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }},
		options: options,
	}
}

// copyHeaders gives the client its own header configuration so later changes to
// the caller's map cannot affect requests. Copying Options alone shares the map.
func copyHeaders(headers map[string]string) map[string]string {
	copied := make(map[string]string, len(headers))
	for name, value := range headers {
		copied[name] = value
	}
	return copied
}

func (c *Client) Close() { c.http.CloseIdleConnections() }

// Request describes an upstream operation. MetricPath overrides the endpoint
// label for dynamic paths, such as individual artwork downloads.
type Request struct {
	Method, Path, MetricPath string
	Body                     io.Reader
	Headers                  map[string]string
	Limit                    int64
}

// Fetch executes a bounded HTTP request without exposing URLs or response bodies in errors.
func (c *Client) Fetch(ctx context.Context, request Request) (data []byte, err error) {
	target := c.options.BaseURL
	if request.Path != "" {
		target = strings.TrimRight(target, "/") + request.Path
	}
	req, err := http.NewRequestWithContext(ctx, request.Method, target, request.Body)
	if err != nil {
		return nil, errors.New("invalid HTTP request")
	}
	for k, v := range c.options.Headers {
		req.Header.Set(k, v)
	}
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}
	metricPath := request.MetricPath
	if metricPath == "" {
		metricPath = req.URL.EscapedPath()
	}
	started := time.Now()
	statusCode := 0
	defer func() {
		duration := time.Since(started)
		if c.options.Metrics != nil {
			timing := RequestTiming{Service: c.options.Service, Method: req.Method, Path: metricPath, StatusCode: statusCode, Success: err == nil, Duration: duration}
			// Report timeouts and cancellations using the reporter's own timeout.
			if reportErr := c.options.Metrics.WriteHTTPRequestDuration(context.WithoutCancel(ctx), timing); reportErr != nil {
				log.Printf("report HTTP request duration: %v", reportErr)
			}
		}
	}()
	resp, err := c.http.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		// net/http errors contain URLs; never log them or remote response bodies.
		return nil, errors.New("HTTP network error or timeout")
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	if !c.options.AcceptStatus(resp.StatusCode) {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}
	data, err = io.ReadAll(io.LimitReader(resp.Body, request.Limit+1))
	if err != nil {
		return nil, errors.New("cannot read HTTP response")
	}
	if int64(len(data)) > request.Limit {
		return nil, errors.New("HTTP response exceeds size limit")
	}
	return data, nil
}
