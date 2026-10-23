# Divoom Ditto Pro Online Catalog Client

A web application for browsing the Divoom community gallery and displaying
16x16 pixel art on a Divoom Ditto Pro. Explore categories and albums, send artwork
to your display, or let automatic rotation alternate animations with your clock.

![Divoom gallery web interface showing animated pixel art previews](docs/demo.gif)

## Features

- Browse gallery categories and Discover albums with animated previews.
- Send still images and animations to your display from the browser.
- Rotate artwork automatically from a chosen category or album, with configurable
  intervals and display duration.
- Restore the clock between automatic animations and customize its style and color.
- Keep autoplay preferences, clock appearance, playback history, and downloaded
  artwork across restarts.
- Access a password-protected interface served by a single Go executable with
  embedded web assets.

Automatic rotation runs in the background while the application is running,
including when no browser is connected. Manual artwork and clock commands pause
rotation; resume it through the web interface when ready.

Gallery lists load metadata without downloading artwork for format checks.
Artwork is downloaded and decoded when needed for a preview or playback, and
successful conversions are cached on disk. Unsupported encodings discovered
during loading are recorded as gallery IDs in `DATA_DIR/skipped.json`, a JSON
array such as `[4203677,4223443]`. Rotation excludes these IDs across categories,
albums, restarts, and playback-history resets, and immediately tries another item
when it discovers an unsupported encoding. Preview requests also reuse this list
to avoid repeated downloads. Network failures and malformed files with recognized
encodings are not added. Stop the app and remove an ID (or clear the array) to
retry it after decoder support changes.

## Requirements

- Go 1.23 or later to build or run from source.
- A Divoom community account for gallery access.
- A running [Divoom device HTTP gateway](https://github.com/lAnubisl/divoom-ditto-pro-client)
  connected to your display, with its URL and API key.

The application fetches artwork from Divoom's hosted gallery and controls the
display through the separate gateway. Run one controller per account and device.

## Getting started

Set the following environment variables in your shell. These PowerShell examples
use placeholders; replace them with your own values:

```powershell
$env:DIVOOM_API_BASE_URL = 'http://YOUR-SERVER:8080'
$env:DIVOOM_API_KEY = 'YOUR-API-KEY'
$env:DIVOOM_GALLERY_EMAIL = 'your@email.example'
$env:DIVOOM_GALLERY_PASSWORD = 'your-gallery-password'
$env:DIVOOM_WEB_PASSWORD = 'your-private-web-password'

go run ./cmd/server
```

Open [http://127.0.0.1:8090/](http://127.0.0.1:8090/) and sign in with
`DIVOOM_WEB_PASSWORD`. The device URL should exclude the `/api` suffix.
[.env.example](.env.example) lists the available settings; the application reads
process environment variables and does not load `.env` files automatically.

On first startup, autoplay starts with the **NEW** category, a five-minute
interval, and a ten-second display duration. Use the web interface to choose a
category or album, adjust timing, or pause playback. Saved preferences take
priority on subsequent starts.

### Optional configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `DIVOOM_WEB_HOST` | `127.0.0.1` | Address the web server listens on |
| `DIVOOM_WEB_PORT` | `8090` | Web server port |
| `DATA_DIR` | `var/lib/divoom_gallery_web/` | Directory for settings, history, and cached artwork |

The default data directory is relative to the working directory. Downloaded
artwork is cached without automatic eviction. For access from another machine,
set `DIVOOM_WEB_HOST=0.0.0.0` and use the server's address. Use a trusted LAN or an
HTTPS reverse proxy for remote access.

Playback history is loaded once at startup and kept in memory. Successful
automatic sends save changes to disk; restart the app to pick up external edits
to `history.json`.

## Metrics

Set `DIVOOM_METRICS_URL` to enable metrics reporting, with
`DIVOOM_METRICS_USERNAME` and `DIVOOM_METRICS_PASSWORD` for HTTP Basic
authentication. There is no default URL. An unset or blank URL disables reporting
and skips credential validation, even when credentials are set.

The service pushes the `divoom_gallery_history_size` gauge (the number of entries
in the persisted history array) from the history store after each successful
history save following a completed automatic send, including source restarts
and replacements that keep the same size. The history store queues a sample after
saving; startup does not send a history sample.

Every outbound gallery and device HTTP attempt also pushes the
`divoom_http_request_duration_seconds` gauge. Labels identify `service`
(`gallery` or `device`), `method`, `path`, `status_code`, and `outcome`
(`success` or `error`). Gallery login, session renewal retries, category and album
requests, and artwork downloads are included, as are all device endpoints.
Artwork downloads share the `/artwork` path label to avoid a separate series for
every file. Labels contain no credentials, query strings, or full URLs.

The shared `internal/httpx.Client` executes HTTP requests and owns their timing
and telemetry reporting. Gallery and device clients configure it once and pass
operation details; they handle API payloads and response validation. Metrics
delivery uses the same helper with telemetry disabled to prevent recursion.

Durations use seconds and cover HTTP execution through reading and closing the
bounded response body. They exclude local work, waiting for the gallery or device
lock, and metrics delivery. Each attempt sends a sample, including HTTP errors,
network failures, timeouts, and cancellations. A `status_code` of `0` means no
HTTP response was received. The outcome describes the HTTP exchange; application
JSON validation happens afterward. This gauge records individual execution
times, rather than cumulative histogram buckets. Metrics delivery uses its own
five-second timeout even when the original request was cancelled.

Samples are timestamped when recorded and queued for a single background worker,
so history saves, gallery requests, and device commands do not wait for metrics
delivery. The queue holds up to 256 pending samples; when full, new samples are
dropped and the failure is logged. Delivery failures are also logged, without
retries or failing the original operation. Shutdown allows up to five seconds to
flush pending samples before cancelling delivery and discarding the remainder.
HTTP attempts also accumulate in a map grouped by `service`, `method`, and
`path`. Once per minute, the worker swaps the map and sends all groups in one
request as the `divoom_http_requests_per_minute` gauge. For example,
`{service="device",method="GET",path="/animation"} 5` means five attempts during
the preceding reporting window, including failures. Counts use the same normalized
paths as durations and include attempts whose duration samples were dropped by a
full queue. Empty windows send no count samples. Counts reset each window; failed
count deliveries are logged without retries. Shutdown also flushes the final
partial window within the existing five-second limit. Duration and history
samples continue to send immediately through the 256-entry queue.

Set the full endpoint URL, for example
`http://<SERVER-IP-ADDRESS>:8428/api/v1/import/prometheus`, using VictoriaMetrics'
[Prometheus import API](https://docs.victoriametrics.com/victoriametrics/url-examples/#apiv1importprometheus).
Use the port, scheme, and path appropriate for your deployment or authentication
proxy. Basic authentication over HTTP is unencrypted; use it on a trusted LAN or
configure HTTPS.
Requests time out after five seconds and do not follow redirects.
For systemd deployments, add these variables to
`/etc/divoom-gallery-web/divoom-gallery-web.env` and restart the service.

## Build

```sh
go vet ./...
go build -o bin/divoom-gallery-web ./cmd/server
```

On Windows, use `bin/divoom-gallery-web.exe` as the output filename. For Linux
service installation and updates, see the [deployment guide](deploy/README.md).

## Documentation

- [Architecture and component diagrams](docs/architecture.md)
- [Device gateway integration](docs/http-client-agent-guide.md)
- [Clock controls](docs/clock-settings.md)
- [Divoom hosted API reference](docs/divoom-api/README.md)
- [Artwork formats and decoding](docs/divoom-api/divoom-artwork-formats.md)
- [Linux service deployment](deploy/README.md)

## Acknowledgements

Gallery protocol and artwork decoding references include
[APIxoo](https://github.com/redphx/apixoo). The web interface uses
[HTMX](https://htmx.org/). See [third-party notices](THIRD_PARTY_NOTICES.md) for
attribution and license details.
