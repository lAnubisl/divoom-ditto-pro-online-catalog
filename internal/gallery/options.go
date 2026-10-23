package gallery

import "divoom-gallery-client/internal/httpx"

type Options struct {
	Email, Password string
	Metrics         httpx.Metrics
	// Optional endpoint overrides support isolated HTTP tests.
	APIURL, FileURL string
}
