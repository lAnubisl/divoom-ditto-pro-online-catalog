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

## Metrics

Set `DIVOOM_METRICS_URL` to enable metrics reporting, with
`DIVOOM_METRICS_USERNAME` and `DIVOOM_METRICS_PASSWORD` for HTTP Basic
authentication. There is no default URL. An unset or blank URL disables reporting
and skips credential validation, even when credentials are set.

The service pushes the `divoom_gallery_history_size` gauge (the number of entries
in the persisted history array) from the rotation selector after each successful
history save, including resets, rollbacks, and replacements that keep the same
size. The selector calls the metrics client directly after saving; startup does
not send a history sample.

After every successful device API request, the device client pushes
`divoom_device_api_last_success_timestamp_seconds`, labelled with `method` and
`path`. Its value is the Unix time in seconds of the successful communication.
This includes the startup health check, image and animation uploads, and clock
commands, for both automatic playback and manual actions. Only HTTP 200 responses
with the expected JSON status count as success; failed requests emit no signal.
The health check confirms API reachability and authentication, not physical
device health.

Each metrics request completes or times out before the calling operation
continues. Failures are logged without failing the history save or device
command. There is no queue, change tracking, or periodic reporting.

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
