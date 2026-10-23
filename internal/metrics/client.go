package metrics

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"divoom-gallery-client/internal/config"
	"divoom-gallery-client/internal/httpx"
)

// Client encodes VictoriaMetrics samples and delivers them through the shared HTTP client.
type Client struct {
	http *httpx.Client
}

func NewClient(settings config.Metrics) *Client {
	return &Client{http: httpx.NewClient(httpx.Options{
		BaseURL: settings.URL, Timeout: 5 * time.Second,
		Headers: map[string]string{
			"Content-Type":  "text/plain; version=0.0.4",
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(settings.Username+":"+settings.Password)),
		},
		AcceptStatus: func(status int) bool { return status >= 200 && status < 300 },
		// Telemetry is disabled for metrics delivery to avoid recursive reporting.
	})}
}

func (c *Client) Close() {
	c.http.Close()
}

func (c *Client) WriteHistorySize(ctx context.Context, size int) error {
	// The Prometheus import API expects timestamps in milliseconds.
	body := fmt.Sprintf("# TYPE divoom_gallery_history_size gauge\ndivoom_gallery_history_size %d %d\n", size, time.Now().UnixMilli())
	return c.write(ctx, body)
}

// WriteHTTPRequestDuration pushes one duration sample per completed HTTP attempt.
func (c *Client) WriteHTTPRequestDuration(ctx context.Context, timing httpx.RequestTiming) error {
	outcome := "error"
	if timing.Success {
		outcome = "success"
	}
	body := fmt.Sprintf("# TYPE divoom_http_request_duration_seconds gauge\ndivoom_http_request_duration_seconds{service=%q,method=%q,path=%q,status_code=%q,outcome=%q} %.9f %d\n",
		timing.Service, timing.Method, timing.Path, fmt.Sprint(timing.StatusCode), outcome, timing.Duration.Seconds(), time.Now().UnixMilli())
	return c.write(ctx, body)
}

func (c *Client) write(ctx context.Context, body string) error {
	_, err := c.http.Fetch(ctx, httpx.Request{Method: http.MethodPost, Body: strings.NewReader(body), Limit: 64 * 1024})
	if err != nil {
		return fmt.Errorf("metrics: %w", err)
	}
	return nil
}
