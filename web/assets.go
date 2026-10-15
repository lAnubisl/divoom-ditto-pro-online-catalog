// Package web contains the web assets embedded in the server executable.
package web

import "embed"

// Assets includes templates and static files, independent of the working directory.
//
//go:embed templates/*.html static
var Assets embed.FS
