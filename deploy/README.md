# Woodpecker deployment

`.woodpecker.yaml` has two steps: `build` and `deploy`. Pushes, pull requests,
and manual runs build the application. Only successful builds for push or manual
events on `main` deploy. There is no test step. Artifacts share the workflow
workspace; no registry or artifact service is needed.

The existing SilverBullet service and hardware notes were inspected before this
implementation. They document a Raspberry Pi Model B Rev 2 and an ARMv6 build:
`CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6`. CI cross-compiles that target,
not the CI runner's architecture. The installer rejects incompatible server
architectures. No server connection was made while implementing this pipeline.

## Woodpecker setup

Activate the repository in Woodpecker and use `.woodpecker.yaml` as its pipeline
configuration. The pipeline uses current Woodpecker syntax, including workflow
`concurrency: 1`. Use a version supporting that field. It queues runs instead of
cancelling an installation. In repository **Settings → General**, also set the
pipeline concurrency limit to **1** if your version exposes that setting. Disable
automatic cancellation of previous runs for this repository. Keep sufficient step
timeout for package installation, up to 600 seconds waiting for the server lock,
and up to two 300-second restarts plus two 210-second health windows. Allow at
least 30 minutes for the deploy step. The server lock independently serializes
installations, including rollback.

In repository **Settings → Secrets**, add each name below with its exact spelling.
Paste its value into the value field and save it. For **every secret**, enable
only **push** and **manual** events; disable pull requests, tags, cron, deployment,
and all other events. Leave **Plugins / Images restrictions empty**: this is a
commands-based deployment, not a plugin. Only `deploy` references the secrets.
All seven values are required; there are no endpoint or account defaults.

| Secret | Value to save |
| --- | --- |
| `deploy_host` | SSH hostname or IP reachable from the CI runner; use an unbracketed IPv6 address if applicable. |
| `deploy_port` | Actual SSH port as a decimal integer. |
| `deploy_user` | Existing SSH login account with the CI public key authorized. |
| `deploy_service_user` | Existing application's service account; must exactly match the installed unit's `User` property. |
| `deploy_health_url` | Full HTTP(S) login/health URL reachable **from the target server**, returning exactly HTTP 200 without authentication or redirects. Use the existing local login endpoint documented in SilverBullet. Do not put credentials in this URL. |
| `deploy_ssh_key` | Complete multiline, passphrase-free OpenSSH private key, including header and footer, with real line breaks. |
| `deploy_known_hosts` | Verified OpenSSH known-hosts entry matching `deploy_host` and `deploy_port`. Nondefault ports require the `[hostname]:port` host field, using the actual secret values. |

Reuse the existing CI key if available. Its public key must be in the deployment
login account's `~/.ssh/authorized_keys`; Woodpecker holds its corresponding
private key. The known-hosts entry is the **server's public host key**, a different
key from the login key. Verify its fingerprint through a trusted console or
existing trusted inventory before saving it. Unverified `ssh-keyscan` output
alone is insufficient. No server-side copy of the CI private key is needed once
it is safely saved in Woodpecker. Do not generate or replace working keys merely
to enable this pipeline. Protect the SSH directory and authorized-keys file with
the usual owner-only write permissions.

The scripts disable shell tracing, require noninteractive public-key
authentication and strict host-key verification, and use temporary mode-0700
directories and mode-0600 key files. They suppress SSH and application diagnostics
that could contain sensitive values. Failure logs are intentionally generic;
inspect service logs privately on the server when needed. Remote settings travel
as archive files, never interpolated remote commands. Do not enable debug tracing
or echo the deployment environment.

## Server prerequisites and authority

This pipeline updates an **existing installation** only. Preserve these locations:

| Item | Existing path |
| --- | --- |
| Executable | `/opt/divoom-gallery-web/divoom-gallery-web` |
| Configuration | `/etc/divoom-gallery-web/divoom-gallery-web.env` |
| Unit | `/etc/systemd/system/divoom-gallery-web.service` |
| Persistent data | `/var/lib/divoom-gallery-web/` |

The pipeline updates the executable while preserving the existing server
configuration, installed service unit, accounts, and unrelated services.
The executable remains owned by numeric UID/GID 0 with mode 0755.
Existing environment-file and data permissions stay intact. The service account
comes solely from its repository secret. The application rejects all CLI
arguments and has no validation-only mode; therefore a separate configuration
check as the service account is not supported. CI never sources its systemd
environment file as shell code. Restarting the existing unit runs the candidate
with the existing service account, environment and sandbox, followed by health
verification. No second controller is launched.

Required tools on the server: Bash at `/bin/bash`, OpenSSH server, `sudo`,
systemd (`systemctl`), GNU coreutils (`install`, `mktemp`, `timeout`, `sha256sum`,
`mv`, `cp`, `cut`, `id`, `date`, `sleep`, `wc`, `cat`, `rm`), GNU tar, grep,
util-linux `flock`, curl with HTTP(S) support, and trusted CA certificates for
HTTPS. The SSH login shell must support POSIX syntax. The service unit and the
three existing deployment files must be available as regular files.

The exact sudo invocation is:

```sh
sudo -n /bin/bash "$stage/install.sh" "$stage"
```

`stage` is a fresh absolute temporary directory created by the remote login
shell. The existing login account must be allowed to run `/bin/bash` as UID 0
with arbitrary script-path and staging-directory arguments, without a password,
TTY, or interactive prompts. An administrator must provision this deliberately
using the local sudo policy; nothing in this repository installs sudoers rules.
The command must work with `sudo -n`, and the sudo policy must not require a TTY.

**This grants the deployment key root authority.** The key can upload and execute
an arbitrary shell script with that permission, and installs an executable used
by a system service. A rule permitting only this path pattern does not confine
that authority. Repository writers, CI administrators, and anyone holding the
key must be trusted accordingly. A more restricted model would require a
separately administered, immutable server-side deployment helper and a different
sudo command; that helper is not silently provisioned here.

Required network access:

- The CI runner must reach the forge, the container registry for the Go and Debian
  images, Debian package repositories over their configured HTTP(S) transports,
  and the target SSH endpoint from the secrets. The deploy step installs its SSH
  client and CA certificates through apt. The Go module currently has no external
  dependencies; future modules may require Go proxy/checksum-service access.
- The target server must reach `deploy_health_url` directly (the check bypasses
  proxy environment variables), and retain its existing access to the device
  gateway, gallery, artwork hosts, and optional metrics endpoint. CI does not
  need access to the application's credentials or its HTTP listener.

## Installation, backups, and recovery

After upload, the installer locks `/run/lock/divoom-gallery-web.lock`, snapshots
the candidate into administrator-owned staging, and verifies its SHA-256 before
installation. It creates a unique directory under
`/var/backups/divoom-gallery-web/` (parent and backup directories mode 0700) with
the previous `executable`, `configuration`, and `service-unit`, preserving their
metadata, plus `backup.sha256` and `candidate.sha256`. Backups contain secrets;
never upload them as CI artifacts. The CI log reports the backup directory only.

The candidate is installed beside the current executable, checked again, then
atomically renamed over it on the same filesystem. Restart has a 300-second
timeout; health checks make at most 30 attempts with a five-second HTTP timeout
and two-second pauses. Success requires both active systemd state and HTTP 200.
Redirects are not followed. Configuration and unit files are only backed up,
never replaced.

Any failure from the replacement onward, including a handled termination signal,
attempts atomic restoration of the previous executable, restarts the unchanged
unit, and checks its health with the same bounds. The CI run remains failed after
successful rollback. Power loss, SIGKILL, or an unavailable server can prevent
automatic recovery; inspect the retained backup before retrying. Earlier failures
leave the installed executable untouched. No schema migration is performed and
mutable application data is not rolled back or included in these snapshots.

Backups are retained indefinitely; there is no automatic deletion. Administrators
should review disk usage and remove older backup directories manually under the
same installation lock, after verifying a newer successful backup and preserving
at least the latest known-good executable and matching configuration/unit.
Arrange independent consistent backups of persistent application data.

For manual recovery, open an administrative shell on the target server, set
`backup` to the chosen protected backup directory reported by CI, then run:

```sh
set -eu
: "${backup:?Set backup to the selected protected backup directory}"
exec 9>/run/lock/divoom-gallery-web.lock
flock -w 600 9
sha256sum --check "$backup/backup.sha256"
systemctl stop divoom-gallery-web.service
replacement=$(mktemp /opt/divoom-gallery-web/divoom-gallery-web.XXXXXXXX)
install -o 0 -g 0 -m 0755 "$backup/executable" "$replacement"
mv -fT -- "$replacement" /opt/divoom-gallery-web/divoom-gallery-web
systemctl start divoom-gallery-web.service
systemctl is-active --quiet divoom-gallery-web.service
```

Run those commands with error checking (`set -eu`) and keep the lock until
verification completes. Privately check the configured health URL with curl and
inspect service logs as needed. Routine executable recovery does not restore the
configuration or unit, because CI never modifies them. If separately damaged,
restore `configuration` and/or `service-unit` from that backup to their paths
above with `cp -p`, preserving ownership/mode. Validate a restored unit using
`systemd-analyze verify`, run `systemctl daemon-reload`, and restart and verify
the service again. Do not restore older credentials or mutable data blindly.

Woodpecker references: [workflow syntax](https://woodpecker-ci.org/docs/usage/workflow-syntax),
[concurrency](https://woodpecker-ci.org/docs/usage/workflows#concurrency), and
[secret event/plugin filters](https://woodpecker-ci.org/docs/usage/secrets).
