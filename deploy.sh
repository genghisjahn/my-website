#!/usr/bin/env bash
set -euo pipefail

REMOTE_USER="genghisjahn"
REMOTE_HOST="ryz-2"
SSH_PORT="22"
REMOTE_DIR="/home/genghisjahn/jonwear.com"
LOCAL_PUBLIC="./public"

# Reusable SSH options
CTL="/tmp/ssh_mux_%h_%p_%r"
SSH_BASE=(ssh -p "$SSH_PORT" -o ControlPath="$CTL")
SSH_MASTER=(ssh -p "$SSH_PORT" -o ControlMaster=auto -o ControlPersist=10m -o ControlPath="$CTL")

echo "Building…"
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

echo "Opening master SSH (one password prompt)…"
"${SSH_MASTER[@]}" -N -f "${REMOTE_USER}@${REMOTE_HOST}"

echo "Ensure remote dir…"
"${SSH_BASE[@]}" "${REMOTE_USER}@${REMOTE_HOST}" "mkdir -p '${REMOTE_DIR}.tmp'"

echo "Rsync…"
RSYNC_SSH="ssh -p ${SSH_PORT} -o ControlPath=${CTL}"
rsync -azP --delete -e "$RSYNC_SSH" "${LOCAL_PUBLIC}/" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}.tmp/"

echo "Activate…"
"${SSH_BASE[@]}" "${REMOTE_USER}@${REMOTE_HOST}" "
  set -e;
  if [ -d '${REMOTE_DIR}' ]; then rm -rf '${REMOTE_DIR}.bak' && mv '${REMOTE_DIR}' '${REMOTE_DIR}.bak'; fi
  mv '${REMOTE_DIR}.tmp' '${REMOTE_DIR}'
"

echo "Close master SSH…"
"${SSH_BASE[@]}" -O exit "${REMOTE_USER}@${REMOTE_HOST}"

echo "Sending webmentions…"
if [ -n "$MY_SITE_WEBMENTION_APP" ]; then
  curl -s -X POST "https://webmention.app/check?token=${MY_SITE_WEBMENTION_APP}&url=https://jonwear.com/feed.xml" > /dev/null &
  curl -s -X POST "https://webmention.app/check?token=${MY_SITE_WEBMENTION_APP}&url=https://jonwear.com/notes/feed.xml" > /dev/null &
  wait
  echo "  Webmentions sent (posts + notes)."
else
  echo "  Skipped (MY_SITE_WEBMENTION_APP not set)"
fi

echo "Done."
