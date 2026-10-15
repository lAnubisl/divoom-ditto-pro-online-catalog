# Linux systemd Service: File Layout Cheat Sheet

## Scenario

A background service runs as a compiled executable, exposes a web interface, and is managed by `systemd`. It uses:

- An **administrator-managed environment file** (the service never modifies it).
- A **JSON settings file** that both the administrator and the service can modify.
- **Persistent local files**, such as images, videos, and other media.

## Recommended directory layout

```text
/
├── opt/myservice/
│   ├── myservice                 # Executable
├── etc/myservice/
│   └── myservice.env             # Admin-managed environment variables
├── var/lib/myservice/
│   ├── settings.json             # Service-managed JSON settings
│   └── media/                    # Persistent media and local files
└── etc/systemd/system/
    └── myservice.service         # systemd unit
```

| Path | Purpose | Writer |
| --- | --- | --- |
| `/opt/myservice/` | Executable, libraries, bundled web assets | Deployment process |
| `/etc/myservice/myservice.env` | Environment variables and administrative configuration | Administrator |
| `/var/lib/myservice/settings.json` | Persistent JSON settings editable through the web UI | Service and administrator |
| `/var/lib/myservice/media/` | Persistent images, videos, and other application files | Service |
| `/etc/systemd/system/myservice.service` | Service startup configuration | Administrator |
| `journald` | Service logs | systemd |

**Rule of thumb:** Keep deployed application files read-only to the service. Keep administrator-owned configuration in `/etc`, and mutable application state in `/var/lib`.

## Ownership and permissions

Run the process as a dedicated system user, `myservice`, rather than as `root`.

| Path | Ownership | Suggested mode |
| --- | --- | --- |
| `/opt/myservice/` | `root:root` | Directories `755`; executable `755` |
| `/etc/myservice/` | `root:root` | `755` or `750` as appropriate |
| `/etc/myservice/myservice.env` | `root:myservice` | `640` |
| `/var/lib/myservice/` | `myservice:myservice` | `750` |
| `/var/lib/myservice/settings.json` | `myservice:myservice` | `640` |
| `/var/lib/myservice/media/` | `myservice:myservice` | `750` |
| `/etc/systemd/system/myservice.service` | `root:root` | `644` |

> If an administrator edits `settings.json`, preserve the file's ownership so that the service can continue updating it. The service needs write permission to the containing directory if it performs atomic file replacement.

## Example systemd unit

Save as `/etc/systemd/system/myservice.service`:

```ini
[Unit]
Description=MyService Web Application
Wants=network-online.target
After=network-online.target

[Service]
Type=exec
User=myservice
Group=myservice

WorkingDirectory=/opt/myservice
ExecStart=/opt/myservice/myservice
EnvironmentFile=/etc/myservice/myservice.env

StateDirectory=myservice
StateDirectoryMode=0750

Restart=on-failure
RestartSec=5

NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true

[Install]
WantedBy=multi-user.target
```

`StateDirectory=myservice` lets systemd manage `/var/lib/myservice` and makes that directory writable to the service even with `ProtectSystem=strict`. Create the `media/` subdirectory and initial `settings.json` during provisioning or on application startup.

## Example environment file

Save as `/etc/myservice/myservice.env`:

```ini
PORT=8080
DATA_DIR=/var/lib/myservice
```

The application reads DATA_DIR and stores settings.json, history.json and media/ beneath it. An `EnvironmentFile` is a systemd environment-variable file, **not** a shell script. Avoid putting secrets in environment variables if a stronger mechanism such as systemd credentials is available.

## Operational guidance

- **Upgrades:** Replace files under `/opt/myservice/` without overwriting `/etc/myservice/` or `/var/lib/myservice/`.
- **Settings updates:** Write JSON to a temporary file in the same directory, then atomically rename it over `settings.json`.
- **Web security:** Never expose `/var/lib/myservice/` as a static web root. Serve media through explicitly controlled routes.
- **Backups:** Back up `/etc/myservice/` and `/var/lib/myservice/`; deployable binaries can usually be recreated.
- **Logging:** Prefer stdout/stderr captured by `journald` instead of maintaining a separate log directory.
- **Network:** Bind to an appropriate interface and use a reverse proxy/TLS where needed. A non-root service typically cannot bind directly to privileged ports below 1024.

## Useful commands

```bash
# Reload units after changing the service file
sudo systemctl daemon-reload

# Start the service and enable it at boot
sudo systemctl enable --now myservice

# Check status
systemctl status myservice

# Follow logs
journalctl -u myservice -f

# Restart after deploying new binaries/configuration
sudo systemctl restart myservice
```

## Quick reference

**Executable and bundled UI:** `/opt/myservice/`  
**Administrator-controlled environment:** `/etc/myservice/myservice.env`  
**Service-editable JSON:** `/var/lib/myservice/settings.json`  
**Persistent media:** `/var/lib/myservice/media/`  
**systemd unit:** `/etc/systemd/system/myservice.service`  
**Logs:** `journalctl -u myservice`
