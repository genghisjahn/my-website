#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Load config file if it exists, env vars override
if [[ -f "$SCRIPT_DIR/.deploy.env" ]]; then
  source "$SCRIPT_DIR/.deploy.env"
fi

# Required config (env vars override config file)
REMOTE_USER="${DEPLOY_USER:?Set DEPLOY_USER in .deploy.env or environment}"
REMOTE_HOST="${DEPLOY_HOST:?Set DEPLOY_HOST in .deploy.env or environment}"
SSH_PORT="${DEPLOY_PORT:-22}"
REMOTE_DIR="${DEPLOY_DIR:?Set DEPLOY_DIR in .deploy.env or environment}"

OUTPUT_BINARY="./site_server"

echo "Establishing SSH control master connection..."
ssh -p "$SSH_PORT" -M -S /tmp/ssh_mux_$REMOTE_HOST -fnNT "${REMOTE_USER}@${REMOTE_HOST}"

echo "Stopping any running instance on $REMOTE_HOST..."
ssh -p "$SSH_PORT" -S /tmp/ssh_mux_$REMOTE_HOST "${REMOTE_USER}@${REMOTE_HOST}" "
  pkill -f ~/web_server/site_server || echo 'No running process found.'
"

echo "Building site_server for linux/amd64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$OUTPUT_BINARY" ./cmd/serve/main.go

echo "Binary built at $OUTPUT_BINARY"

# Copy the binary to the remote server

echo "Copying binary to ${REMOTE_USER}@${REMOTE_HOST}:~/web_server/..."
scp -P "$SSH_PORT" -o ControlPath=/tmp/ssh_mux_$REMOTE_HOST "$OUTPUT_BINARY" "${REMOTE_USER}@${REMOTE_HOST}:~/web_server/"
rm "$OUTPUT_BINARY"

echo "Starting remote server..."
ssh -p "$SSH_PORT" -S /tmp/ssh_mux_$REMOTE_HOST "${REMOTE_USER}@${REMOTE_HOST}" "
  nohup ~/web_server/site_server \
    -public \"$REMOTE_DIR\" \
    -css \"$REMOTE_DIR/css\" \
    -images \"$REMOTE_DIR/images\" \
    -addr \"0.0.0.0:8088\" \
    >> ~/web_server/site_server.log 2>&1 < /dev/null &
  disown
"

ssh -p "$SSH_PORT" -S /tmp/ssh_mux_$REMOTE_HOST -O exit "${REMOTE_USER}@${REMOTE_HOST}"

echo "Done."
