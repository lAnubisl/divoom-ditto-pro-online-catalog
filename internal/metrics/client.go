package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"divoom-gallery-client/internal/config"
	"divoom-gallery-client/internal/httpx"
)

// Client owns VictoriaMetrics request encoding, authentication, and transport.
type Client struct {
	settings config.Metrics
	http     *http.Client
}

func NewClient(settings config.Metrics) *Client {
	return &Client{settings: settings, http: httpx.NewClient(5 * time.Second)}
}

func (c *Client) Close() {
	c.http.CloseIdleConnections()
}

func (c *Client) WriteHistorySize(ctx context.Context, size int) error {
	// The Prometheus import API expects timestamps in milliseconds.
	body := fmt.Sprintf("# TYPE divoom_gallery_history_size gauge\ndivoom_gallery_history_size %d %d\n", size, time.Now().UnixMilli())
	return c.write(ctx, body)
}

// WriteDeviceAPISuccess emits a signal for a validated successful device request.
// A timestamp gauge allows consumers to check how recently communication succeeded.
func (c *Client) WriteDeviceAPISuccess(ctx context.Context, method, path string) error {
	now := time.Now()
	body := fmt.Sprintf("# TYPE divoom_device_api_last_success_timestamp_seconds gauge\ndivoom_device_api_last_success_timestamp_seconds{method=%q,path=%q} %.3f %d\n", method, path, float64(now.UnixMilli())/1000, now.UnixMilli())
	return c.write(ctx, body)
}

func (c *Client) write(ctx context.Context, body string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.settings.URL, strings.NewReader(body))
	if err != nil {
		return errors.New("invalid metrics HTTP request")
	}
	req.Header.Set("Content-Type", "text/plain; version=0.0.4")
	req.SetBasicAuth(c.settings.Username, c.settings.Password)
	resp, err := c.http.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return errors.New("metrics HTTP network error or timeout")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("metrics HTTP status %d", resp.StatusCode)
	}
	return nil
}
