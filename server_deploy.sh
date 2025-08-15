

#!/usr/bin/env bash
set -euo pipefail

OUTPUT_BINARY="./site_server"

# Configuration
REMOTE_USER="genghisjahn"
REMOTE_HOST="ryz-2"
REMOTE_DIR="/home/${REMOTE_USER}/jonwear.com"

echo "Stopping any running instance on $REMOTE_HOST..."
ssh "${REMOTE_USER}@${REMOTE_HOST}" "
  pkill -f ~/web_server/site_server || echo 'No running process found.'
"

echo "Building site_server for linux/amd64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$OUTPUT_BINARY" ./cmd/serve/main.go

echo "Binary built at $OUTPUT_BINARY"

# Copy the binary to the remote server

echo "Copying binary to ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}..."
scp "$OUTPUT_BINARY" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/"
rm "$OUTPUT_BINARY"

echo "Starting remote server..."
ssh "${REMOTE_USER}@${REMOTE_HOST}" "
  nohup ~/web_server/site_server \
    -public \"$REMOTE_DIR\" \
    -css \"$REMOTE_DIR/css\" \
    -images \"$REMOTE_DIR/images\" \
    -addr \"0.0.0.0:8088\" \
    >> ~/web_server/site_server.log 2>&1 < /dev/null &
  disown
"