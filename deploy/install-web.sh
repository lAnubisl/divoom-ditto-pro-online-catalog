#!/bin/sh
# Run as root with a directory containing the binary and this systemd unit.
set -eu
stage=${1:?Usage: install-web.sh /path/to/staging-directory}
if [ "$(id -u)" -ne 0 ]; then
    echo "Run this installer as root" >&2
    exit 1
fi
case "$(uname -m)" in
    aarch64) binary=divoom-gallery-web-linux-arm64 ;;
    armv6l|armv7l) binary=divoom-gallery-web-linux-armv6 ;;
    *) echo "Unsupported server architecture" >&2; exit 1 ;;
esac
test -f "$stage/$binary"
test -f "$stage/divoom-gallery-web.service"
env_file=/etc/divoom-gallery-web/divoom-gallery-web.env
if [ ! -f "$env_file" ]; then
    echo "Create $env_file from .env.example and configure it before installing" >&2
    exit 1
fi
if ! getent group divoom-gallery-web >/dev/null; then
    groupadd --system divoom-gallery-web
fi
if ! getent passwd divoom-gallery-web >/dev/null; then
    useradd --system --gid divoom-gallery-web --home-dir /var/lib/divoom-gallery-web --no-create-home --shell /bin/false divoom-gallery-web
fi
if systemctl cat divoom-gallery-web.service >/dev/null 2>&1; then
    systemctl stop divoom-gallery-web.service
fi
install -d -o root -g root -m 0755 /opt/divoom-gallery-web /etc/divoom-gallery-web
chown root:divoom-gallery-web "$env_file"
chmod 0640 "$env_file"
install -d -o divoom-gallery-web -g divoom-gallery-web -m 0750 /var/lib/divoom-gallery-web /var/lib/divoom-gallery-web/media
install -o root -g root -m 0755 "$stage/$binary" /opt/divoom-gallery-web/divoom-gallery-web
install -o root -g root -m 0644 "$stage/divoom-gallery-web.service" /etc/systemd/system/divoom-gallery-web.service
systemd-analyze verify /etc/systemd/system/divoom-gallery-web.service
systemctl daemon-reload
systemctl enable --now divoom-gallery-web.service
systemctl status --no-pager divoom-gallery-web.service
echo "Open http://<server-hostname>:<DIVOOM_WEB_PORT>/ (default port 8090) and sign in with DIVOOM_WEB_PASSWORD from $env_file"
