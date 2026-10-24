# Architecture

This document describes the current implementation in this repository. The
application is a single Go executable containing an HTTP server, embedded web
assets, and a background autoplay scheduler. It stores state in local files and
controls the display through a separately deployed HTTP gateway.

The diagrams use Mermaid and render directly in GitHub Markdown or a Markdown
viewer with Mermaid support. The component diagrams use the following styles:

| Element | Style |
| --- | --- |
| Internal components | Blue rounded rectangles |
| Browser, external services, and device | Gray rectangles |
| Local file storage | Amber cylinders |
| Calls and data access | Solid arrows `-->` |
| Device commands | Thick arrows `==>` |
| Notifications, telemetry, and startup wiring | Dotted arrows `-.->` |
| Connector routing | Straight segments with gently rounded turns (`curve: rounded`) |

Responses are implicit in the component diagrams. The sequence diagram uses
matching blue participants and amber notes, with solid request arrows and dashed
response arrows.

## System overview

```mermaid
%%{init: {"flowchart": {"curve": "rounded"}}}%%
flowchart LR
    browser["Browser<br/>HTML, HTMX, JavaScript"]:::external

    subgraph host["Application host"]
        app("Go executable<br/>HTTP server + background autoplay<br/>default: 127.0.0.1:8090"):::service
        disk[("Local files<br/>DATA_DIR")]:::storage
    end

    api["Divoom gallery API<br/>appin.divoom-gz.com"]:::external
    files["Divoom artwork server<br/>fin.divoom-gz.com"]:::external
    gateway["Separate device HTTP gateway<br/>DIVOOM_API_BASE_URL"]:::external
    display["Divoom Ditto Pro<br/>16 x 16 display"]:::external
    metrics["Optional metrics endpoint<br/>VictoriaMetrics Prometheus import API"]:::external

    browser -->|"HTTP: pages, fragments, previews, commands"| app
    app -->|"Read/write preferences, history, skipped IDs, GIF cache"| disk
    app -->|"HTTPS: login, category and album metadata"| api
    app -->|"HTTPS: download original artwork"| files
    app ==>|"HTTP(S) + X-Api-Key: media and clock commands"| gateway
    gateway ==>|"Device transport managed by gateway"| display
    app -.->|"Queued HTTP(S) POST + Basic auth"| metrics

    classDef service fill:#dbeafe,stroke:#2563eb,color:#172554
    classDef external fill:#f1f5f9,stroke:#64748b,color:#0f172a
    classDef storage fill:#fef3c7,stroke:#d97706,color:#451a03
```

Gallery credentials and the device API key stay on the application server. The
browser communicates with this application, which fetches gallery data and
forwards media or clock commands. The gateway owns the physical device
connection; its transport implementation is outside this repository.

## Components inside the Go process

```mermaid
%%{init: {"flowchart": {"curve": "rounded"}}}%%
flowchart TB
    boot("cmd/server + internal/config<br/>Environment configuration and startup wiring"):::service

    subgraph ui["Web interface"]
        handlers("internal/handlers<br/>Routing, login, session cookie, CSRF<br/>Pages, previews, commands, status"):::service
        assets("web<br/>Embedded templates, CSS, HTMX, JavaScript"):::service
    end

    subgraph control["Playback coordination"]
        playback("internal/playback.Controller<br/>Manual commands and preference changes"):::service
        runner("internal/rotation.Runner<br/>Schedule, artwork selection, clock restoration"):::service
    end

    subgraph content["Gallery and media"]
        source("internal/source.Adapter<br/>Autoplay pages and shared animation loading"):::service
        gallery("internal/gallery.Client<br/>Metadata, downloads, login and session renewal"):::service
        decoder("internal/animation<br/>Decode Divoom artwork into 16 x 16 GIF"):::service
    end

    subgraph storage["Shared state and local persistence"]
        settings[("internal/settings<br/>settings.json")]:::storage
        history("internal/history.Store<br/>Ordered slice, Contains, Record<br/>Retention, persistence, size metrics"):::service
        historyFile[("DATA_DIR/history.json<br/>Ordered filename array")]:::storage
        skipped[("internal/skipped<br/>skipped.json")]:::storage
        media[("internal/media<br/>media/ID.gif")]:::storage
    end

    device("internal/device.Client<br/>Media and clock API<br/>Shared display lock"):::service
    httpx("internal/httpx.Client<br/>Shared HTTP implementation<br/>Bounded responses and request timing"):::service
    metrics("internal/metrics.Client<br/>256-sample queue + delivery worker<br/>Request-count map + minute timer"):::service
    upstream["Divoom gallery API and artwork server"]:::external
    gateway["External device gateway"]:::external
    sink["External metrics endpoint"]:::external

    boot -.->|"Constructs dependencies and starts HTTP server + runner"| handlers
    boot -.-> runner
    boot -.->|"Construct store; load history once"| history
    handlers -->|"Render and serve"| assets
    handlers -->|"Browse metadata"| gallery
    handlers -->|"GIF previews"| source
    handlers -->|"Send, clock, autoplay changes"| playback
    handlers -->|"Read preferences"| settings
    handlers -->|"Read live status"| runner
    playback -->|"Persist changes"| settings
    playback -.->|"Notify changes; cancel on manual takeover"| runner
    playback -->|"Load GIF"| source
    playback ==>|"Lock and issue manual commands"| device
    runner -->|"Read schedule and clock preferences"| settings
    runner -->|"Load GIF"| source
    runner ==>|"Lock, send, wait, restore clock"| device
    runner -->|"rotation.History interface: Contains / Record"| history
    history -->|"Load at startup; atomically save completed sends"| historyFile
    runner -->|"Fetch configured source pages"| source
    history -.->|"History-size samples after saves"| metrics
    source -->|"Read source preferences"| settings
    source -->|"Fetch pages and original artwork"| gallery
    source -->|"Read/write converted GIFs"| media
    source -->|"Read/write unsupported IDs"| skipped
    source -->|"Convert cache misses"| decoder
    gallery -->|"Gallery HTTP requests"| httpx
    device ==>|"Device HTTP requests"| httpx
    httpx -->|"Gallery requests"| upstream
    httpx ==>|"Device requests"| gateway
    httpx -.->|"HTTP attempts: duration samples + grouped counts"| metrics
    metrics -.->|"Delivery with telemetry disabled"| httpx
    httpx -.->|"Metrics POST when enabled"| sink

    classDef service fill:#dbeafe,stroke:#2563eb,color:#172554
    classDef external fill:#f1f5f9,stroke:#64748b,color:#0f172a
    classDef storage fill:#fef3c7,stroke:#d97706,color:#451a03
```

These are logical components within one process, not separate services. Gallery,
device, and metrics clients use separate configured instances of the shared HTTP
helper. Metrics delivery disables its own telemetry to prevent recursion.
Playback and rotation consume interfaces backed by the same source adapter and
device client, so previews and sends reuse the same media cache, and manual and
automatic device operations share one lock. The runner uses the `rotation.History`
interface; `history.Store` owns the ordered slice, file persistence, retention, and
history-size metrics. Membership checks scan the slice.

## Main runtime flows

### Browse and preview

1. A signed-in browser requests category or album metadata through the handlers
   and gallery client. Listing metadata does not download artwork for decoding.
2. A preview request calls the source adapter. Known unsupported IDs are rejected
   before loading. A valid cached GIF is returned immediately.
3. On a cache miss, the gallery client downloads the original artwork, the decoder
   converts it to GIF, and the media store saves it. The handler returns the GIF.
   Unsupported encodings are recorded in `skipped.json` to avoid repeated work.

### Automatic playback

```mermaid
%%{init: {"theme": "base", "themeVariables": {"actorBkg": "#dbeafe", "actorBorder": "#2563eb", "actorTextColor": "#172554", "noteBkgColor": "#fef3c7", "noteBorderColor": "#d97706", "noteTextColor": "#451a03"}}}%%
sequenceDiagram
    participant R as Rotation runner
    participant A as Source adapter
    participant H as History store
    participant F as history.json
    participant M as Metrics client
    participant D as Device client / gateway

    H->>F: Load ordered filenames once during construction
    F-->>H: Saved history, or empty if missing
    Note over H: Keep an ordered slice in memory
    Note over R,D: Runner reads saved timing and source; no browser is required
    R->>D: Acquire shared display lock
    R->>A: Fetch configured category/album pages
    A-->>R: Eligible metadata, excluding known unsupported IDs
    loop Until first unplayed candidate or source exhausted
        R->>H: Contains(filename)
        H-->>R: Membership from slice scan
    end
    Note over R: If all candidates were played, select first candidate
    R->>A: Load cached GIF or download and decode
    A-->>R: GIF bytes
    R->>D: POST /api/images or /api/animations
    D-->>R: Send completed
    R->>H: Record(filename) with separate timeout
    Note over H: Prepare updated slice; reset if filename already present
    Note over H: Append filename; keep newest 1,000 entries
    H->>F: Atomically save ordered filenames
    F-->>H: Save succeeded
    H->>H: Replace active slice
    opt Metrics enabled
        H->>M: WriteHistorySize(retained count)
        Note over H,M: Metrics errors are logged; record remains successful
    end
    H-->>R: Record succeeded
    R->>R: Wait configured display duration
    R->>D: POST /api/display/clock
    D-->>R: Clock command completed
    R->>D: Release shared display lock
    Note over R: Wait until next scheduled cycle
```

The sequence shows a successful cycle. The runner chooses the first eligible
item absent from recent history, restarting history when the source is exhausted.
History is loaded once when the history store is constructed at application startup.
The store keeps an ordered slice in memory and uses `slices.Contains` for membership
checks. Adding a new filename beyond the 1,000-entry limit evicts the oldest entries.
Recording an already played filename resets history to that filename. The runner
uses a `Contains`/`Record` interface and records only completed sends. The store
updates memory and reports the retained size only after a successful save;
external edits to the history file require an application restart.
Loading or sending failures leave history unchanged; unsupported encodings cause
an immediate attempt with another item. A failed history save is reported as a
cycle error and may allow the sent artwork to repeat. History saves after a
completed send use a separate timeout so a concurrent pause does not cancel them.
A clock restoration failure stops autoplay and exposes an error in the UI.

### Manual commands

A send or clock command passes through the playback controller's command gate.
The controller persists the autoplay pause and notifies the runner to cancel
automatic work, including pending clock restoration during manual takeover. It
then acquires the shared display lock and sends the command. Autoplay remains
paused until resumed. A successful gateway response confirms send completion;
it does not verify the physical display's rendered output.

## Persistence and deployment

| Location | Owner | Contents |
| --- | --- | --- |
| `DATA_DIR/settings.json` | `internal/settings` | Autoplay source, timing, enabled state, clock appearance |
| `DATA_DIR/history.json` | `internal/history` | Successfully sent automatic artwork filenames; failures leave history unchanged |
| `DATA_DIR/skipped.json` | `internal/skipped`, populated by source adapter | Gallery IDs with unsupported artwork encodings |
| `DATA_DIR/media/<ID>.gif` | `internal/media`, populated by source adapter | Converted GIFs retained without automatic eviction |

The default data directory is `var/lib/divoom_gallery_web/`, relative to the
working directory. The supplied systemd service uses
`/var/lib/divoom-gallery-web` and starts the executable from
`/opt/divoom-gallery-web`. It reads environment configuration from
`/etc/divoom-gallery-web/divoom-gallery-web.env` and restarts on failure.
Templates and static assets are embedded at build time; no separate frontend
server, database, or message broker is required. Run one controller per account,
device, and media directory.

Metrics are enabled only when `DIVOOM_METRICS_URL` is set. History saves and
gallery/device HTTP attempts enqueue samples for one background worker. Its queue
holds 256 pending samples; full queues drop new samples, and delivery failures
are logged without failing the originating operation. Shutdown allows up to five
seconds to flush the queue and final partial request-count window. HTTP duration
and history samples send immediately. The metrics client also counts HTTP attempts
in a map keyed by service, method, and path. Every minute it swaps the map and
sends one batch containing `divoom_http_requests_per_minute` samples for all active
groups; new attempts accumulate in the next window. Counts include failed attempts
and duration samples dropped by a full queue. Empty windows send no count samples;
failed count deliveries are logged without retries.

## Implementation references

- [Startup and dependency wiring](../cmd/server/main.go)
- [HTTP server and scheduler lifecycle](../cmd/server/web.go)
- [Routes and web security](../internal/handlers/server.go)
- [Manual command coordination](../internal/playback/controller.go)
- [Automatic playback and selection](../internal/rotation/runner.go)
- [History interface](../internal/rotation/history.go) and [history store](../internal/history/store.go)
- [Gallery/media adapter](../internal/source/adapter.go) and [gallery API client](../internal/gallery/client.go)
- [Device gateway client](../internal/device/client.go) and [integration contract](http-client-agent-guide.md)
- [Shared HTTP helper](../internal/httpx/http.go) and [metrics worker](../internal/metrics/client.go)
- [systemd service](../deploy/divoom-gallery-web.service)
