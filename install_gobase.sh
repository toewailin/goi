#!/bin/bash

# GoBase installation script

# Set the URL for the GoBase release binary
GOBASE_URL="https://github.com/toewailin/gobase/releases/download/v1.0.0-alpha/gobase-osx"

# Set the download directory (default to $HOME/Downloads)
DOWNLOAD_DIR="$HOME/Downloads"
GOBASE_BINARY="$DOWNLOAD_DIR/gobase-osx"

# Download GoBase binary
echo "Downloading GoBase from $GOBASE_URL..."
curl -L -o "$GOBASE_BINARY" "$GOBASE_URL"

# Check if the download was successful
if [[ ! -f "$GOBASE_BINARY" ]]; then
    echo "Failed to download the GoBase binary."
    exit 1
fi

# Make the binary executable
echo "Making GoBase executable..."
chmod +x "$GOBASE_BINARY"

# Move the binary to /usr/local/bin
echo "Installing GoBase to /usr/local/bin..."
sudo mv "$GOBASE_BINARY" /usr/local/bin/gobase

# Verify installation
if command -v gobase &>/dev/null; then
    echo "GoBase installed successfully!"
    gobase --version
else
    echo "GoBase installation failed."
    exit 1
fi
