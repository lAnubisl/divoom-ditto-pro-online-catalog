#!/bin/bash
# Values stay in files or argument arrays; none become remote shell source.
set +x
set -euo pipefail
umask 077
fail() { printf '%s\n' "$1" >&2; exit 1; }
for name in DEPLOY_HOST DEPLOY_PORT DEPLOY_USER DEPLOY_SERVICE_USER DEPLOY_HEALTH_URL DEPLOY_SSH_KEY DEPLOY_KNOWN_HOSTS; do
    [[ -n ${!name:-} ]] || fail "Missing required deployment setting: $name"
done
[[ $DEPLOY_HOST =~ ^[a-zA-Z0-9:][a-zA-Z0-9.:%_-]*$ ]] || fail 'Invalid SSH host setting'
[[ $DEPLOY_PORT =~ ^[0-9]{1,5}$ ]] || fail 'Invalid SSH port setting'
(( 10#$DEPLOY_PORT >= 1 && 10#$DEPLOY_PORT <= 65535 )) || fail 'Invalid SSH port setting'
[[ $DEPLOY_USER =~ ^[a-zA-Z_][a-zA-Z0-9_-]*$ ]] || fail 'Invalid SSH login setting'
[[ $DEPLOY_SERVICE_USER =~ ^[a-zA-Z_][a-zA-Z0-9_-]*$ ]] || fail 'Invalid service account setting'
[[ $DEPLOY_HEALTH_URL == http://* || $DEPLOY_HEALTH_URL == https://* ]] || fail 'Invalid health URL setting'
[[ $DEPLOY_HEALTH_URL != *$'\n'* && $DEPLOY_HEALTH_URL != *$'\r'* ]] || fail 'Invalid health URL setting'
stage=$(mktemp -d)
trap 'rm -rf -- "$stage"' EXIT
printf '%s\n' "$DEPLOY_SSH_KEY" > "$stage/key"
printf '%s\n' "$DEPLOY_KNOWN_HOSTS" > "$stage/known_hosts"
ssh-keygen -y -P '' -f "$stage/key" >/dev/null 2>&1 || fail 'SSH key must be valid and passphrase-free'
cp dist/candidate dist/SHA256SUMS "$stage/"
cp deploy/install.sh "$stage/install.sh"
printf '%s' "$DEPLOY_SERVICE_USER" > "$stage/service_user"
printf '%s' "$DEPLOY_HEALTH_URL" > "$stage/health_url"
unset DEPLOY_SSH_KEY DEPLOY_KNOWN_HOSTS DEPLOY_SERVICE_USER DEPLOY_HEALTH_URL
ssh_options=(
    -F /dev/null -T -p "$DEPLOY_PORT" -l "$DEPLOY_USER"
    -i "$stage/key" -o IdentitiesOnly=yes -o IdentityAgent=none
    -o BatchMode=yes -o PasswordAuthentication=no -o KbdInteractiveAuthentication=no
    -o StrictHostKeyChecking=yes -o "UserKnownHostsFile=$stage/known_hosts"
    -o GlobalKnownHostsFile=/dev/null -o UpdateHostKeys=no
    -o ConnectTimeout=15 -o ConnectionAttempts=1
    -o ServerAliveInterval=15 -o ServerAliveCountMax=3 -o LogLevel=ERROR
)
# The remote command is constant. The archive carries all remote settings as data.
if ! tar -C "$stage" -cf - candidate SHA256SUMS install.sh service_user health_url |
    ssh "${ssh_options[@]}" "$DEPLOY_HOST" '
        set -eu
        umask 077
        stage=$(mktemp -d)
        trap '\''rm -rf -- "$stage"'\'' EXIT
        tar -xf - -C "$stage"
        sudo -n /bin/bash "$stage/install.sh" "$stage"
    ' 2>/dev/null; then
    fail 'Deployment failed; inspect protected server backups and service state before retrying'
fi
