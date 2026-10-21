package device

import (
	"bytes"
	"context"
	"divoom-gallery-client/internal/config"
	"divoom-gallery-client/internal/httpx"
	medialimits "divoom-gallery-client/internal/media"
	"divoom-gallery-client/internal/settings"
	"encoding/json"
	"errors"
	"image/gif"
	"log"
	"net/http"
	"time"
)

const (
	requestTimeout  = 210 * time.Second
	maxResponseSize = 64 * 1024 // 64 KiB.
)

// Client sends media and clock requests to a configured device HTTP API.
type Client struct {
	http         *http.Client
	baseURL, key string
	gate         chan struct{}
	metrics      Metrics
}

// New creates a client with a deadline longer than the server operation deadline.
func New(d config.DeviceAPI, metrics Metrics) *Client {
	return &Client{http: httpx.NewClient(requestTimeout), baseURL: d.BaseURL, key: d.APIKey, gate: make(chan struct{}, 1), metrics: metrics}
}

func (d *Client) request(ctx context.Context, method, path string, body []byte, contentType, expected string) error {
	headers := map[string]string{"X-Api-Key": d.key}
	if contentType != "" {
		headers["Content-Type"] = contentType
	}
	data, err := httpx.Fetch(ctx, d.http, method, d.baseURL+path, bytes.NewReader(body), headers, maxResponseSize)
	if err != nil {
		return err
	}
	var result apiResponse
	if json.Unmarshal(data, &result) != nil || result.Status != expected {
		return errors.New("unexpected device API success response")
	}
	if d.metrics != nil {
		if err := d.metrics.WriteDeviceAPISuccess(ctx, method, path); err != nil {
			log.Printf("report device API success: %v", err)
		}
	}
	return nil
}

// Health checks API reachability and authentication, not physical device health.
func (d *Client) Health(ctx context.Context) error {
	return d.request(ctx, "GET", "/api/health", nil, "", "ok")
}

// SendAnimation sends a one-frame GIF as an image, otherwise as an animation.
func (d *Client) SendAnimation(ctx context.Context, media []byte) error {
	if len(media) == 0 || len(media) > medialimits.MaxAnimationSize {
		return errors.New("invalid animation upload size")
	}
	decoded, err := gif.DecodeAll(bytes.NewReader(media))
	if err != nil {
		return errors.New("invalid GIF upload")
	}
	path := "/api/animations"
	if len(decoded.Image) == 1 {
		path = "/api/images"
	}
	return d.request(ctx, "POST", path, media, "image/gif", "sent")
}

// Lock serializes complete playback cycles with manual device actions.
func (d *Client) Lock(ctx context.Context) error {
	select {
	case d.gate <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (d *Client) Unlock() { <-d.gate }

// ShowClock sends the supplied preferences; the client owns no settings state.
func (d *Client) ShowClock(ctx context.Context, clock settings.Clock) error {
	if err := clock.Validate(); err != nil {
		return err
	}
	body, err := json.Marshal(clock)
	if err != nil {
		return err
	}
	return d.request(ctx, "POST", "/api/display/clock", body, "application/json", "sent")
}
