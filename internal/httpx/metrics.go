package httpx

import (
	"context"
	"time"
)

// RequestTiming describes one HTTP attempt, including reading the response body.
// Path is a stable endpoint label, never a URL or an individual artwork path.
// Success describes the HTTP exchange, before application response validation.
type RequestTiming struct {
	Service, Method, Path string
	StatusCode            int // Zero when no HTTP response was received.
	Success               bool
	Duration              time.Duration
}

type Metrics interface {
	WriteHTTPRequestDuration(context.Context, RequestTiming) error
}
