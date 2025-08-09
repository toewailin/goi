#!/bin/bash

# GoI installation script

# Set the URL for the GoI release binary
GOI_URL="https://github.com/toewailin/goi/releases/download/v1.0.2/goi-osx"

# Set the download directory (default to $HOME/Downloads)
DOWNLOAD_DIR="$HOME/Downloads"
GOI_BINARY="$DOWNLOAD_DIR/goi-osx"

# Download GoI binary
echo "Downloading GoI from $GOI_URL..."
curl -L -o "$GOI_BINARY" "$GOI_URL"

# Check if the download was successful
if [[ ! -f "$GOI_BINARY" ]]; then
    echo "Failed to download the GoI binary."
    exit 1
fi

# Make the binary executable
echo "Making GoI executable..."
chmod +x "$GOI_BINARY"

# Move the binary to /usr/local/bin
echo "Installing GoI to /usr/local/bin..."
sudo mv "$GOI_BINARY" /usr/local/bin/goi

# Verify installation
if command -v goi &>/dev/null; then
    echo "GoI installed successfully!"
    goi version
else
    echo "GoI installation failed."
    exit 1
fi
