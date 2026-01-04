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
USER_NAME="gotimedate"

echo "--- Uninstalling gotimedate ---"

if systemctl is-active --quiet "$SERVICE_NAME"; then
    echo "Stopping service..."
    systemctl stop "$SERVICE_NAME"
fi

if systemctl is-enabled --quiet "$SERVICE_NAME" 2>/dev/null; then
    echo "Disabling service..."
    systemctl disable "$SERVICE_NAME"
fi

if [ -f "/etc/systemd/system/$SERVICE_NAME.service" ]; then
    echo "Removing service file..."
    rm -f "/etc/systemd/system/$SERVICE_NAME.service"
    systemctl daemon-reload
fi

if [ -d "$INSTALL_DIR" ]; then
    echo "Removing installation directory $INSTALL_DIR..."
    rm -rf "$INSTALL_DIR"
fi

if id "$USER_NAME" &>/dev/null; then
    echo "Removing system user $USER_NAME..."
    userdel -r "$USER_NAME" 2>/dev/null || userdel "$USER_NAME"
fi

echo "--- Uninstallation complete! ---"
