#!/usr/bin/env bash
#
# install.sh — SuperMac installer
# Usage: curl -fsSL https://raw.githubusercontent.com/CosmoLabs-org/SuperMac/master/install.sh | bash
#
set -euo pipefail

REPO="CosmoLabs-org/SuperMac"
BINARY="mac"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}[supermac]${NC} $*"; }
warn()  { echo -e "${YELLOW}[supermac]${NC} $*"; }
error() { echo -e "${RED}[supermac]${NC} $*" >&2; exit 1; }

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    arm64) ARCH="arm64" ;;
    x86_64) ARCH="amd64" ;;
    *) error "Unsupported architecture: $ARCH" ;;
esac

# Get latest release tag
get_latest_version() {
    if command -v gh &>/dev/null; then
        gh release view --repo "$REPO" --json tagName --jq '.tagName' 2>/dev/null && return
    fi
    curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" 2>/dev/null \
        | grep '"tag_name"' \
        | head -1 \
        | sed -E 's/.*"([^"]+)".*/\1/'
}

VERSION="${INSTALL_VERSION:-$(get_latest_version)}"
if [ -z "$VERSION" ]; then
    error "Could not determine latest version. Set INSTALL_VERSION manually."
fi

DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/mac-darwin-$ARCH.tar.gz"

info "Installing SuperMac $VERSION ($ARCH)..."
info "Downloading from $DOWNLOAD_URL"

# Download and extract
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMPDIR/mac.tar.gz"; then
    error "Download failed. Check that version $VERSION exists."
fi

# Verify checksum if available
CHECKSUM_URL="https://github.com/$REPO/releases/download/$VERSION/checksums.txt"
if curl -fsSL "$CHECKSUM_URL" -o "$TMPDIR/checksums.txt" 2>/dev/null; then
    info "Verifying checksum..."
    cd "$TMPDIR"
    if command -v shasum &>/dev/null; then
        EXPECTED=$(grep "mac-darwin-$ARCH.tar.gz" checksums.txt | awk '{print $1}')
        ACTUAL=$(shasum -a 256 mac.tar.gz | awk '{print $1}')
        if [ "$EXPECTED" != "$ACTUAL" ]; then
            error "Checksum mismatch! Expected $EXPECTED, got $ACTUAL"
        fi
        info "Checksum verified."
    fi
fi

# Extract
tar xzf "$TMPDIR/mac.tar.gz" -C "$TMPDIR"
if [ ! -f "$TMPDIR/mac-$ARCH" ]; then
    error "Expected binary mac-$ARCH not found in archive"
fi

# Install
mkdir -p "$INSTALL_DIR" 2>/dev/null || true
if [ ! -w "$INSTALL_DIR" ]; then
    warn "$INSTALL_DIR not writable, using sudo..."
    sudo cp "$TMPDIR/mac-$ARCH" "$INSTALL_DIR/$BINARY"
    sudo chmod +x "$INSTALL_DIR/$BINARY"
else
    cp "$TMPDIR/mac-$ARCH" "$INSTALL_DIR/$BINARY"
    chmod +x "$INSTALL_DIR/$BINARY"
fi

# Verify
if command -v "$BINARY" &>/dev/null; then
    INSTALLED=$("$BINARY" --version 2>/dev/null || echo "$VERSION")
    info "Successfully installed SuperMac $INSTALLED to $INSTALL_DIR/$BINARY"
    info "Run 'mac help' to get started."
else
    warn "Binary installed to $INSTALL_DIR/$BINARY but not found in PATH."
    warn "Add $INSTALL_DIR to your PATH or move the binary."
fi
