package metrics

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"divoom-gallery-client/internal/config"
	"divoom-gallery-client/internal/httpx"
)

const queueSize = 256

// requestKey uses the same endpoint labels as duration samples.
type requestKey struct {
	service, method, path string
}

// Client queues VictoriaMetrics samples for delivery by a single background worker.
type Client struct {
	http      *httpx.Client
	queue     chan string
	done      chan struct{}
	cancel    context.CancelFunc
	mu        sync.Mutex
	requests  map[requestKey]uint64
	closed    bool
	closeOnce sync.Once
}

func NewClient(settings config.Metrics) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{requests: make(map[requestKey]uint64), queue: make(chan string, queueSize), done: make(chan struct{}), cancel: cancel, http: httpx.NewClient(httpx.Options{
		BaseURL: settings.URL, Timeout: 5 * time.Second,
		Headers: map[string]string{
			"Content-Type":  "text/plain; version=0.0.4",
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(settings.Username+":"+settings.Password)),
		},
		AcceptStatus: func(status int) bool { return status >= 200 && status < 300 },
		// Telemetry is disabled for metrics delivery to avoid recursive reporting.
	})}
	go c.run(ctx)
	return c
}

// Close stops accepting samples and allows up to five seconds to flush the queue
// and request counts from the final partial minute.
// It is safe to call concurrently with reporting or another Close.
func (c *Client) Close() {
	c.closeOnce.Do(func() {
		c.mu.Lock()
		c.closed = true
		close(c.queue)
		c.mu.Unlock()
		timer := time.NewTimer(5 * time.Second)
		defer timer.Stop()
		select {
		case <-c.done:
		case <-timer.C:
			log.Print("metrics: shutdown flush timed out; discarding pending samples")
			c.cancel()
			<-c.done
		}
		c.cancel()
		c.http.Close()
	})
}

// WriteHTTPRequestDuration queues one duration sample and counts the attempt
// for the current minute, even if the duration queue is full.
// Reporting survives cancellation of the original request.
func (c *Client) WriteHTTPRequestDuration(_ context.Context, timing httpx.RequestTiming) error {
	outcome := "error"
	if timing.Success {
		outcome = "success"
	}
	body := fmt.Sprintf("# TYPE divoom_http_request_duration_seconds gauge\ndivoom_http_request_duration_seconds{service=%q,method=%q,path=%q,status_code=%q,outcome=%q} %.9f %d\n",
		timing.Service, timing.Method, timing.Path, fmt.Sprint(timing.StatusCode), outcome, timing.Duration.Seconds(), time.Now().UnixMilli())
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		c.requests[requestKey{timing.Service, timing.Method, timing.Path}]++
	}
	return c.enqueueLocked(body)
}

// enqueueLocked never waits for network delivery or queue space. Errors indicate
// that the sample was dropped; delivery failures are logged by the worker.
func (c *Client) enqueueLocked(body string) error {
	if c.closed {
		return errors.New("metrics: client closed")
	}
	select {
	case c.queue <- body:
		return nil
	default:
		return errors.New("metrics: queue full; sample dropped")
	}
}

func (c *Client) run(ctx context.Context) {
	defer close(c.done)
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.flushRequestCounts(ctx)
		case body, ok := <-c.queue:
			if ctx.Err() != nil {
				return
			}
			if !ok {
				c.flushRequestCounts(ctx)
				return
			}
			if err := c.write(ctx, body); err != nil {
				log.Printf("report metrics: %v", err)
			}
		}
	}
}

// flushRequestCounts snapshots and resets counts before delivery so concurrent
// reports belong to the next window. Known groups keep reporting even at zero.
func (c *Client) flushRequestCounts(ctx context.Context) {
	c.mu.Lock()
	counts := make(map[requestKey]uint64, len(c.requests))
	for key, count := range c.requests {
		counts[key] = count
		c.requests[key] = 0
	}
	c.mu.Unlock()
	if len(counts) == 0 {
		return
	}
	var body strings.Builder
	body.WriteString("# TYPE divoom_http_requests_per_minute gauge\n")
	timestamp := time.Now().UnixMilli()
	for key, count := range counts {
		fmt.Fprintf(&body, "divoom_http_requests_per_minute{service=%q,method=%q,path=%q} %d %d\n", key.service, key.method, key.path, count, timestamp)
	}
	if err := c.write(ctx, body.String()); err != nil {
		log.Printf("report request counts: %v", err)
	}
}

func (c *Client) write(ctx context.Context, body string) error {
	_, err := c.http.Fetch(ctx, httpx.Request{Method: http.MethodPost, Body: strings.NewReader(body), Limit: 64 * 1024})
	if err != nil {
		return fmt.Errorf("metrics: %w", err)
	}
	return nil
}
