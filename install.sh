#!/bin/sh
# Miner installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/miner/main/install.sh | sh

set -e

REPO="YOUR_USERNAME/miner"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest release version
echo "Fetching latest release..."
VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "Error: Could not fetch latest version"
    exit 1
fi

echo "Latest version: $VERSION"

# Construct download URL
if [ "$OS" = "windows" ]; then
    FILENAME="miner-${OS}-${ARCH}.zip"
else
    FILENAME="miner-${OS}-${ARCH}.tar.gz"
fi

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

echo "Downloading $FILENAME..."
TMP_DIR="$(mktemp -d)"
cd "$TMP_DIR"

if ! curl -fsSL -o "$FILENAME" "$DOWNLOAD_URL"; then
    echo "Error: Failed to download from $DOWNLOAD_URL"
    exit 1
fi

# Extract
echo "Extracting..."
if [ "$OS" = "windows" ]; then
    unzip -q "$FILENAME"
    BINARY="miner-${OS}-${ARCH}.exe"
else
    tar xzf "$FILENAME"
    BINARY="miner-${OS}-${ARCH}"
fi

# Install
echo "Installing to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY" "$INSTALL_DIR/miner"
    chmod +x "$INSTALL_DIR/miner"
else
    echo "Requesting elevated privileges to install to $INSTALL_DIR..."
    sudo mv "$BINARY" "$INSTALL_DIR/miner"
    sudo chmod +x "$INSTALL_DIR/miner"
fi

# Cleanup
cd - > /dev/null
rm -rf "$TMP_DIR"

echo ""
echo "✓ Miner $VERSION installed successfully!"
echo ""
echo "Next steps:"
echo "  1. Run: sudo miner install"
echo "  2. Start: miner"
echo "  3. Access Adminer at: http://miner.local:88"
echo ""
