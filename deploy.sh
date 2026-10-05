#!/usr/bin/env bash
set -euo pipefail

# Deploys the site to server01: the generated public/ tree plus the Go server
# that serves it, run as a systemd service behind cloudflared.
#
# Modelled on onesheet's deploy.sh. Replaces the old two-script setup
# (deploy.sh + server_deploy.sh) that targeted the basement box, ryz-2.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Optional: MY_SITE_WEBMENTION_APP lives here if it is not in the shell.
if [[ -f .deploy.env ]]; then
  source .deploy.env
fi

BINARY="site_server"
REMOTE="root@server01"
REMOTE_PATH="/srv/mywebsite"
SERVICE_NAME="mywebsite"
SERVICE_USER="mywebsite"
SERVICE_GROUP="mywebsite"
SITE_URL="jonwear.com"
LOCAL_PUBLIC="./public"

# Loopback only; cloudflared routes the site here. Must match
# deploy/mywebsite.service. 8088-8092 are taken on server01 (Mercury Falling,
# newsjawn, onesheet, ctrlgov); the port check below catches any other clash.
HEALTH_PORT="8093"

UNIT_PATH="/etc/systemd/system/${SERVICE_NAME}.service"
UNIT_BACKUP="/etc/systemd/system/${SERVICE_NAME}.service.prev"

# public/ is synced with --delete, so the server mirrors whichever branch is
# checked out. Only main is a complete picture of the site.
BRANCH="$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo detached)"
if [[ "${BRANCH}" != "main" && "${ALLOW_ANY_BRANCH:-}" != "1" ]]; then
  echo "Refusing to deploy from '${BRANCH}'."
  echo "public/ syncs with --delete. Merge to main and deploy from there."
  echo "If you really mean it: ALLOW_ANY_BRANCH=1 ./deploy.sh"
  exit 1
fi

# === SSH multiplexing (type password once) ===
# ConnectTimeout so an unreachable host fails in seconds instead of hanging.
SSH_SOCK="/tmp/deploy-ssh-$$"
SSH=(ssh -o ControlPath="$SSH_SOCK")
echo "Opening SSH to ${REMOTE}…"
ssh -o ConnectTimeout=10 -o ControlMaster=yes -o ControlPersist=120 -o ControlPath="$SSH_SOCK" -fN "${REMOTE}"
trap 'ssh -o ControlPath="$SSH_SOCK" -O exit "${REMOTE}" 2>/dev/null' EXIT

echo "Building site…"
go run ./cmd/build

echo "Converting images to WebP…"
total_before=0
total_after=0
while IFS= read -r img; do
  [ -f "$img" ] || continue
  size_before=$(stat -f%z "$img")
  total_before=$((total_before + size_before))
  webp="${img%.*}.webp"
  cwebp -q 80 "$img" -o "$webp" -quiet && rm "$img"
  size_after=$(stat -f%z "$webp")
  total_after=$((total_after + size_after))
  # Show per-file stats
  name=$(basename "$webp")
  kb_before=$((size_before / 1024))
  kb_after=$((size_after / 1024))
  file_pct=$((( size_before - size_after ) * 100 / size_before))
  echo "  ${name}: ${kb_before}KB → ${kb_after}KB (${file_pct}%)"
done < <(find "${LOCAL_PUBLIC}/images" -type f \( -iname "*.png" -o -iname "*.jpg" -o -iname "*.jpeg" \) 2>/dev/null)
if [ $total_before -gt 0 ]; then
  saved=$((total_before - total_after))
  pct=$((saved * 100 / total_before))
  before_kb=$((total_before / 1024))
  after_kb=$((total_after / 1024))
  saved_kb=$((saved / 1024))
  echo "  Total: ${before_kb}KB → ${after_kb}KB (saved ${saved_kb}KB, ${pct}%)"
fi

# amd64 to match server01. Pure Go, no cgo.
echo "Building ${BINARY} for linux/amd64…"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "${BINARY}" ./cmd/serve

# Something else on the port would make the health check pass against the
# wrong service, so refuse unless the listener is our own unit.
echo "Checking port ${HEALTH_PORT}…"
"${SSH[@]}" "${REMOTE}" "
  if ss -ltnH 'sport = :${HEALTH_PORT}' | grep -q . && ! systemctl is-active --quiet ${SERVICE_NAME}; then
    echo '!! Port ${HEALTH_PORT} is in use by something other than ${SERVICE_NAME}:'
    ss -ltnpH 'sport = :${HEALTH_PORT}'
    exit 1
  fi"

# Idempotent: a system account with no home and no login shell.
echo "Ensuring service account and remote path exist…"
"${SSH[@]}" "${REMOTE}" "set -e
  getent group ${SERVICE_GROUP} >/dev/null || groupadd --system ${SERVICE_GROUP}
  id -u ${SERVICE_USER} >/dev/null 2>&1 || useradd --system --gid ${SERVICE_GROUP} \
      --no-create-home --home-dir /nonexistent --shell /usr/sbin/nologin \
      --comment 'mywebsite static server' ${SERVICE_USER}
  install -d -o root -g ${SERVICE_GROUP} -m 750 ${REMOTE_PATH}
  install -d -o root -g ${SERVICE_GROUP} -m 750 ${REMOTE_PATH}/public"

# Copy to .new and rename: rename over a running executable is fine on Linux,
# writing into it is not (ETXTBSY).
echo "Copying binary…"
scp -o ControlPath="$SSH_SOCK" "${BINARY}" "${REMOTE}:${REMOTE_PATH}/${BINARY}.new"
rm "${BINARY}"
"${SSH[@]}" "${REMOTE}" "mv ${REMOTE_PATH}/${BINARY}.new ${REMOTE_PATH}/${BINARY}"

echo "Syncing public/…"
rsync -az --delete --exclude '.DS_Store' \
  -e "ssh -o ControlPath=${SSH_SOCK}" \
  "${LOCAL_PUBLIC}/" "${REMOTE}:${REMOTE_PATH}/public/"

# rsync -a carries this Mac's uid, which means nothing on the server, so the
# ownership pass runs every deploy.
echo "Setting ownership and permissions…"
"${SSH[@]}" "${REMOTE}" "set -e
  chown root:${SERVICE_GROUP} ${REMOTE_PATH}/${BINARY}
  chmod 750 ${REMOTE_PATH}/${BINARY}
  chown -R root:${SERVICE_GROUP} ${REMOTE_PATH}/public
  chmod -R u+rwX,g+rX,o-rwx ${REMOTE_PATH}/public"

echo "Installing systemd unit…"
"${SSH[@]}" "${REMOTE}" "test -f ${UNIT_PATH} && cp -a ${UNIT_PATH} ${UNIT_BACKUP} || true"
scp -o ControlPath="$SSH_SOCK" "deploy/${SERVICE_NAME}.service" "${REMOTE}:${UNIT_PATH}"
"${SSH[@]}" "${REMOTE}" "systemctl daemon-reload && systemctl enable ${SERVICE_NAME} && systemctl restart ${SERVICE_NAME} || true"

# Retried, because a single curl straight after restart races systemd.
echo "Waiting for the listener…"
if "${SSH[@]}" "${REMOTE}" "
  for i in \$(seq 1 10); do
    curl -fsS http://127.0.0.1:${HEALTH_PORT}/ >/dev/null 2>&1 && exit 0
    sleep 1
  done
  exit 1"; then
  echo "  <- ${SERVICE_NAME} (${HEALTH_PORT}) ok"
else
  echo ""
  echo "!! HEALTH CHECK FAILED -- rolling back to the previous unit."
  "${SSH[@]}" "${REMOTE}" "journalctl -u ${SERVICE_NAME} -n 40 --no-pager" || true
  if "${SSH[@]}" "${REMOTE}" "test -f ${UNIT_BACKUP}"; then
    "${SSH[@]}" "${REMOTE}" "set -e
      cp -a ${UNIT_BACKUP} ${UNIT_PATH}
      systemctl daemon-reload
      systemctl restart ${SERVICE_NAME}" || true
    echo "!! Rolled back. State now:"
    "${SSH[@]}" "${REMOTE}" "systemctl is-active ${SERVICE_NAME}" || true
  else
    echo "!! No previous unit to restore (${UNIT_BACKUP} absent). Needs hands."
  fi
  exit 1
fi

echo "Sending webmentions…"
if [ -n "${MY_SITE_WEBMENTION_APP:-}" ]; then
  curl -s -X POST "https://webmention.app/check?token=${MY_SITE_WEBMENTION_APP}&url=https://${SITE_URL}/feed.xml" > /dev/null &
  curl -s -X POST "https://webmention.app/check?token=${MY_SITE_WEBMENTION_APP}&url=https://${SITE_URL}/notes/feed.xml" > /dev/null &
  wait
  echo "  Webmentions sent (posts + notes)."
else
  echo "  Skipped (MY_SITE_WEBMENTION_APP not set)"
fi

echo "Done. https://${SITE_URL}/"
