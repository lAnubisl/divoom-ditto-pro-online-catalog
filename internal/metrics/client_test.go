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

func sampleTiming(milliseconds int) httpx.RequestTiming {
	return httpx.RequestTiming{Service: "gallery", Method: "GET", Path: "/artwork",
		StatusCode: 200, Success: true, Duration: time.Duration(milliseconds) * time.Millisecond}
}

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
	go func() { returned <- client.WriteHTTPRequestDuration(ctx, sampleTiming(42)) }()
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
		if err := client.WriteHTTPRequestDuration(ctx, sampleTiming(i)); err != nil {
			t.Fatalf("enqueue sample %d: %v", i, err)
		}
	}
	if err := client.WriteHTTPRequestDuration(ctx, sampleTiming(999)); err == nil {
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
	if err := client.WriteHTTPRequestDuration(context.Background(), sampleTiming(7)); err != nil {
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
	if len(bodies) != 3 {
		t.Fatalf("close delivered %d requests, want 3", len(bodies))
	}
	if !strings.Contains(bodies[0], "divoom_http_request_duration_seconds") ||
		!strings.Contains(bodies[1], `service="gallery",method="GET",path="/artwork",status_code="200",outcome="success"} 0.250000000 `) {
		t.Fatalf("unexpected samples: %q", bodies)
	}
	if !strings.Contains(bodies[2], `divoom_http_requests_per_minute{service="gallery",method="GET",path="/artwork"} 2 `) {
		t.Fatalf("missing shutdown request count: %q", bodies[2])
	}
	if err := client.WriteHTTPRequestDuration(context.Background(), sampleTiming(8)); err == nil {
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
	if err := client.WriteHTTPRequestDuration(context.Background(), sampleTiming(1)); err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("background delivery did not start")
	}
	if err := client.WriteHTTPRequestDuration(context.Background(), sampleTiming(2)); err != nil {
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
				_ = client.WriteHTTPRequestDuration(context.Background(), sampleTiming(j))
			}
			client.Close()
		}()
	}
	workers.Wait()
}

func TestRequestCountsGroupAndReset(t *testing.T) {
	received := make(chan string, 3)
	var bodies []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		received <- string(body)
	}))
	defer server.Close()
	// Exercise window processing directly without waiting a minute or running
	// the duration worker. Delivery still uses the real HTTP metrics transport.
	client := &Client{
		http:  httpx.NewClient(httpx.Options{BaseURL: server.URL, Timeout: time.Second}),
		queue: make(chan string, queueSize), requests: make(map[requestKey]uint64),
	}
	defer client.http.Close()
	report := func(service, method, path string, success bool) {
		t.Helper()
		if err := client.WriteHTTPRequestDuration(context.Background(), httpx.RequestTiming{
			Service: service, Method: method, Path: path, Success: success,
		}); err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 5; i++ {
		report("device", "GET", "/animation", i != 0)
	}
	report("device", "POST", "/animation", true)
	report("device", "GET", "/clock", true)
	report("gallery", "GET", "/animation", true)
	if len(received) != 0 {
		t.Fatal("counts sent before window processing")
	}
	client.flushRequestCounts(context.Background())
	select {
	case body := <-received:
		bodies = append(bodies, body)
	case <-time.After(time.Second):
		t.Fatal("count batch was not delivered")
	}
	if len(bodies) != 1 {
		t.Fatalf("sent %d batches, want 1", len(bodies))
	}
	for _, want := range []string{
		`service="device",method="GET",path="/animation"} 5 `,
		`service="device",method="POST",path="/animation"} 1 `,
		`service="device",method="GET",path="/clock"} 1 `,
		`service="gallery",method="GET",path="/animation"} 1 `,
	} {
		if !strings.Contains(bodies[0], want) {
			t.Fatalf("missing %q in %q", want, bodies[0])
		}
	}
	client.flushRequestCounts(context.Background())
	select {
	case body := <-received:
		for _, labels := range []string{
			`service="device",method="GET",path="/animation"`,
			`service="device",method="POST",path="/animation"`,
			`service="device",method="GET",path="/clock"`,
			`service="gallery",method="GET",path="/animation"`,
		} {
			if !strings.Contains(body, labels+"} 0 ") {
				t.Fatalf("missing zero count for %s in %q", labels, body)
			}
		}
	case <-time.After(time.Second):
		t.Fatal("zero count batch was not delivered")
	}
	report("device", "GET", "/animation", true)
	client.flushRequestCounts(context.Background())
	select {
	case body := <-received:
		bodies = append(bodies, body)
	case <-time.After(time.Second):
		t.Fatal("next count batch was not delivered")
	}
	if len(bodies) != 2 || !strings.Contains(bodies[1], `service="device",method="GET",path="/animation"} 1 `) {
		t.Fatalf("new window did not reset count: %q", bodies)
	}
}

func TestRequestCountsSurviveFullDurationQueue(t *testing.T) {
	client := &Client{queue: make(chan string, 1), requests: make(map[requestKey]uint64)}
	timing := httpx.RequestTiming{Service: "device", Method: "GET", Path: "/animation"}
	if err := client.WriteHTTPRequestDuration(context.Background(), timing); err != nil {
		t.Fatal(err)
	}
	if err := client.WriteHTTPRequestDuration(context.Background(), timing); err == nil {
		t.Fatal("expected full queue")
	}
	key := requestKey{"device", "GET", "/animation"}
	if client.requests[key] != 2 {
		t.Fatalf("count=%d want 2", client.requests[key])
	}
	client.closed = true
	if err := client.WriteHTTPRequestDuration(context.Background(), timing); err == nil {
		t.Fatal("closed client accepted report")
	}
	if client.requests[key] != 2 {
		t.Fatal("closed client changed count")
	}
}
