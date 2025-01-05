#!/bin/bash

set -e

# GitHub repo information
REPO="bigjk/clai"
BINARY_NAME="clai"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Convert architecture to Go format
case "$ARCH" in
    "x86_64")
        ARCH="amd64"
        ;;
    "aarch64" | "arm64")
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Handle OS-specific settings
case "$OS" in
    "darwin")
        INSTALL_DIR="/usr/local/bin"
        ;;
    "linux")
        INSTALL_DIR="/usr/local/bin"
        ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

# Create temp directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Detecting latest version..."
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "Failed to detect latest version"
    exit 1
fi

echo "Latest version: $LATEST_VERSION"
echo "Downloading CLAI for $OS $ARCH..."

# Download URL
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/clai-$OS-$ARCH.tar.gz"

# Download and extract
curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/clai.tar.gz"
tar -xzf "$TMP_DIR/clai.tar.gz" -C "$TMP_DIR"

# Ensure install directory exists
sudo mkdir -p "$INSTALL_DIR"

# Install binary
echo "Installing to $INSTALL_DIR..."
sudo mv "$TMP_DIR/clai" "$INSTALL_DIR/"
sudo chmod +x "$INSTALL_DIR/clai"

echo "Installation complete! CLAI version $LATEST_VERSION has been installed to $INSTALL_DIR/clai"
echo "You can now run 'clai' from anywhere."

# Verify installation
if command -v clai >/dev/null 2>&1; then
    echo "Verification: CLAI is successfully installed and accessible from PATH"
else
    echo "Warning: CLAI is installed but might not be accessible from PATH"
    echo "You might need to add $INSTALL_DIR to your PATH"
fi
