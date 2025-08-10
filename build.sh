#!/bin/bash

# Ensure the 'build' folder exists
mkdir -p build

# Build for the current platform (local)
go build -ldflags="-s -w" -trimpath -o build/goi .

# Build for Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/goi-linux .

# Build for macOS (64-bit)
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/goi-macos .

# Build for Windows (64-bit) (Note: must specify .exe for Windows)
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/goi-win.exe .
