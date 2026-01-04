#!/bin/bash

set -e

if [ ! -f /etc/debian_version ]; then
    echo "Error: This script is only compatible with Ubuntu/Debian."
    exit 1
fi

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root (use sudo)"
  exit 1
fi

INSTALL_DIR="/opt/gotimedate"
SERVICE_NAME="gotimedate"
BINARY_NAME="gotimedate"
USER_NAME="gotimedate"
REPO_URL="https://github.com/shabilullah/gotimedate.git"

echo "--- gotimedate Service Management ---"

echo "[1/5] Checking dependencies..."
export DEBIAN_FRONTEND=noninteractive
MISSING_PACKAGES=()

if ! command -v git >/dev/null 2>&1; then MISSING_PACKAGES+=("git"); fi
if ! command -v go >/dev/null 2>&1; then MISSING_PACKAGES+=("golang-go"); fi

if [ ${#MISSING_PACKAGES[@]} -gt 0 ]; then
    echo "Installing missing packages: ${MISSING_PACKAGES[*]}..."
    apt-get update
    apt-get install -y "${MISSING_PACKAGES[@]}"
else
    echo "All dependencies already installed."
fi

if [ -d "$INSTALL_DIR/.git" ]; then
    echo "[2/5] Update detected. Updating source code..."
    
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        echo "Stopping service for update..."
        systemctl stop "$SERVICE_NAME"
    fi

    cd "$INSTALL_DIR"
    git config --global --add safe.directory "$INSTALL_DIR" || true
    
    git -c safe.directory="$INSTALL_DIR" fetch --all
    git -c safe.directory="$INSTALL_DIR" reset --hard origin/master
    echo "Source code updated to latest master."
else
    echo "[2/5] New installation. Cloning repository..."
    if [ -d "$INSTALL_DIR" ]; then
        echo "Target directory exists but is not a git repo. Cleaning up..."
        rm -rf "$INSTALL_DIR"
    fi
    git clone "$REPO_URL" "$INSTALL_DIR"
    cd "$INSTALL_DIR"
    git config --global --add safe.directory "$INSTALL_DIR" || true
fi

echo "[3/5] Building binary natively..."
mkdir -p build
go build -v -o "build/$BINARY_NAME" main.go
chmod +x "build/$BINARY_NAME"

if ! id -u "$USER_NAME" >/dev/null 2>&1; then
    echo "[4/5] Creating system user '$USER_NAME'..."
    useradd -r -s /bin/false "$USER_NAME"
fi

chown -Rv "$USER_NAME":"$USER_NAME" "$INSTALL_DIR"

echo "[5/5] Configuring systemd service..."
if [ ! -f "/etc/systemd/system/$SERVICE_NAME.service" ]; then
    cat <<EOF > /etc/systemd/system/$SERVICE_NAME.service
[Unit]
Description=Go TimeDate API Service
After=network.target

[Service]
Type=simple
User=$USER_NAME
Group=$USER_NAME
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/build/$BINARY_NAME
Restart=always
RestartSec=5
ReadWritePaths=$INSTALL_DIR

[Install]
WantedBy=multi-user.target
EOF
    systemctl daemon-reload
    systemctl enable "$SERVICE_NAME"
fi

echo "Restarting service..."
systemctl restart "$SERVICE_NAME"

echo "--- Operation complete! ---"
echo "Service status:"
systemctl status "$SERVICE_NAME" --no-pager

echo ""
echo "Access the API at: http://localhost:8080"
echo "Web Interface: http://localhost:8080/ws-test"
echo "Documentation: http://localhost:8080/swagger/index.html"
echo ""
echo "Configuration: $INSTALL_DIR/config.env"
echo "Logs: $INSTALL_DIR/server.log"
