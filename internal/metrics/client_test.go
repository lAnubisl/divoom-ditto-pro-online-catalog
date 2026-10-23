package metrics

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"divoom-gallery-client/internal/config"
	"divoom-gallery-client/internal/httpx"
)

func TestReportingDoesNotWaitForDelivery(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() { close(started) })
		<-release
	}))
	defer server.Close()
	client := NewClient(config.Metrics{URL: server.URL})
	defer client.Close()
	defer close(release)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	returned := make(chan error, 1)
	go func() { returned <- client.WriteHistorySize(ctx, 42) }()
	select {
	case err := <-returned:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("reporting waited for network delivery")
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("cancelled caller prevented background delivery")
	}
	for i := 0; i < queueSize; i++ {
		if err := client.WriteHistorySize(ctx, i); err != nil {
			t.Fatalf("enqueue sample %d: %v", i, err)
		}
	}
	if err := client.WriteHistorySize(ctx, 999); err == nil {
		t.Fatal("full queue did not reject new sample")
	}
}

func TestCloseFlushesSamplesAndRejectsFurtherWrites(t *testing.T) {
	var mu sync.Mutex
	var bodies []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "user" || password != "secret" {
			t.Error("missing metrics authentication")
		}
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		defer mu.Unlock()
		bodies = append(bodies, string(body))
		// A delivery failure must not prevent the next queued sample being sent.
		if len(bodies) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))
	defer server.Close()
	client := NewClient(config.Metrics{URL: server.URL, Username: "user", Password: "secret"})
	if err := client.WriteHistorySize(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
	if err := client.WriteHTTPRequestDuration(context.Background(), httpx.RequestTiming{
		Service: "gallery", Method: "GET", Path: "/artwork", StatusCode: 200,
		Success: true, Duration: 250 * time.Millisecond,
	}); err != nil {
		t.Fatal(err)
	}
	client.Close()
	client.Close()
	mu.Lock()
	defer mu.Unlock()
	if len(bodies) != 2 {
		t.Fatalf("close delivered %d samples, want 2", len(bodies))
	}
	if !strings.Contains(bodies[0], "divoom_gallery_history_size 7 ") ||
		!strings.Contains(bodies[1], `service="gallery",method="GET",path="/artwork",status_code="200",outcome="success"} 0.250000000 `) {
		t.Fatalf("unexpected samples: %q", bodies)
	}
	if err := client.WriteHistorySize(context.Background(), 8); err == nil {
		t.Fatal("closed client accepted a sample")
	}
}

func TestCloseBoundsFlushTime(t *testing.T) {
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started <- struct{}{}
		<-release
	}))
	defer server.Close()
	defer close(release)
	client := NewClient(config.Metrics{URL: server.URL})
	defer client.Close()
	if err := client.WriteHistorySize(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("background delivery did not start")
	}
	if err := client.WriteHistorySize(context.Background(), 2); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() { client.Close(); close(done) }()
	select {
	case <-done:
	case <-time.After(7 * time.Second):
		t.Fatal("shutdown did not bound the total flush time")
	}
}

func TestConcurrentCloseAndReporting(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	client := NewClient(config.Metrics{URL: server.URL})
	var workers sync.WaitGroup
	for i := 0; i < 20; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for j := 0; j < 20; j++ {
				_ = client.WriteHistorySize(context.Background(), j)
			}
			client.Close()
		}()
	}
	workers.Wait()
}
