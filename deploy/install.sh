#!/bin/bash
# Invoked through noninteractive sudo. This is an update, not a first installer.
set +x
set -euo pipefail
umask 077
fail() { printf '%s\n' "$1" >&2; exit 1; }
[[ $(id -u) == 0 ]] || fail 'Installer requires administrative authority'
[[ $# == 1 ]] || fail 'Staging directory is required'
binary=/opt/divoom-gallery-web/divoom-gallery-web
configuration=/etc/divoom-gallery-web/divoom-gallery-web.env
unit=/etc/systemd/system/divoom-gallery-web.service
service=divoom-gallery-web.service
backup_root=/var/backups/divoom-gallery-web
case $(uname -m) in
    armv6l|armv7l) ;;
    *) fail 'Server architecture does not match the ARMv6 artifact' ;;
esac

# The lock spans snapshots, installation, health verification, and rollback.
exec 9>/run/lock/divoom-gallery-web.lock
flock -w 600 9 || fail 'Installation lock timed out'
install -d -o 0 -g 0 -m 0700 "$backup_root"
work=$(mktemp -d "$backup_root/.stage.XXXXXXXX")
replacement=
changed=0
backup=
finish() {
    rc=$?
    trap - EXIT HUP INT TERM
    if (( changed )); then
        printf '%s\n' 'Deployment failed; restoring previous executable' >&2
        # Keep failure status even when recovery succeeds. Never expose app logs.
        if install -o 0 -g 0 -m 0755 "$backup/executable" "$replacement" 2>/dev/null &&
            mv -fT -- "$replacement" "$binary" 2>/dev/null &&
            timeout 300 systemctl restart "$service" >/dev/null 2>&1 && health; then
            printf '%s\n' 'Previous executable restored and health verified' >&2
        else
            printf '%s\n' 'Automatic recovery failed; manual recovery is required' >&2
        fi
        rc=1
    fi
    [[ -z $replacement ]] || rm -f -- "$replacement"
    rm -rf -- "$work"
    exit "$rc"
}
trap finish EXIT
trap 'exit 1' HUP INT TERM

# Snapshot uploaded inputs into administrator-owned storage before using them.
for name in candidate SHA256SUMS service_user health_url; do
    [[ -f $1/$name && ! -L $1/$name ]] || fail 'Missing or invalid uploaded input'
    install -o 0 -g 0 -m 0600 "$1/$name" "$work/$name"
done
cd "$work"
[[ $(wc -l < SHA256SUMS) == 1 ]] || fail 'Invalid artifact checksum manifest'
grep -Eq '^[0-9a-f]{64}  candidate$' SHA256SUMS || fail 'Invalid artifact checksum manifest'
sha256sum --check --status SHA256SUMS || fail 'Uploaded artifact checksum mismatch'
service_user=$(cat service_user)
health_url=$(cat health_url)
[[ $service_user =~ ^[a-zA-Z_][a-zA-Z0-9_-]*$ ]] || fail 'Invalid service account setting'
[[ $health_url == http://* || $health_url == https://* ]] || fail 'Invalid health URL setting'
[[ $health_url != *$'\n'* && $health_url != *$'\r'* ]] || fail 'Invalid health URL setting'
[[ $(systemctl show -p User --value "$service") == "$service_user" ]] || fail 'Service account does not match the installed unit'
id "$service_user" >/dev/null 2>&1 || fail 'Service account does not exist'
for path in "$binary" "$configuration" "$unit"; do
    [[ -f $path && ! -L $path ]] || fail 'Expected an existing regular deployment file'
done
health() {
    local attempt code
    for (( attempt=0; attempt<30; attempt++ )); do
        if systemctl is-active --quiet "$service"; then
            code=$(curl --disable --silent --output /dev/null --write-out '%{http_code}' \
                --noproxy '*' --proto '=http,https' --connect-timeout 2 --max-time 5 \
                --url "$health_url" 2>/dev/null) || code=
            [[ $code == 200 ]] && return 0
        fi
        sleep 2
    done
    return 1
}
backup=$(mktemp -d "$backup_root/backup-$(date -u +%Y%m%dT%H%M%SZ)-XXXXXXXX")
cp -p -- "$binary" "$backup/executable"
cp -p -- "$configuration" "$backup/configuration"
cp -p -- "$unit" "$backup/service-unit"
cp -- SHA256SUMS "$backup/candidate.sha256"
sha256sum "$backup/executable" "$backup/configuration" "$backup/service-unit" > "$backup/backup.sha256"
printf 'Backup saved: %s\n' "$backup"

# The application explicitly rejects CLI arguments and has no config-check mode.
# Do not source the systemd EnvironmentFile as shell code or launch a second controller.
printf '%s\n' 'No application configuration-check mode; verifying through service startup'
replacement=$(mktemp "${binary}.XXXXXXXX")
install -o 0 -g 0 -m 0755 candidate "$replacement"
[[ $(sha256sum "$replacement" | cut -d ' ' -f 1) == $(cut -d ' ' -f 1 SHA256SUMS) ]] || fail 'Installed candidate checksum mismatch'
changed=1
mv -fT -- "$replacement" "$binary"
timeout 300 systemctl restart "$service" >/dev/null 2>&1 || fail 'Service restart failed'
health || fail 'Candidate health verification failed'
changed=0
printf '%s\n' 'Deployment completed successfully'
