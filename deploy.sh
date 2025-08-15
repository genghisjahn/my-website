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

echo "Done."
