#!/bin/bash

# Build script for deptakeover

echo "🔨 Building deptakeover..."

mkdir -p build

# Windows
echo "📦 Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o build/deptakeover-windows-amd64.exe ./cmd/deptakeover

# Linux
echo "📦 Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o build/deptakeover-linux-amd64 ./cmd/deptakeover

# macOS Intel
echo "📦 Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o build/deptakeover-macos-amd64 ./cmd/deptakeover

# macOS ARM
echo "📦 Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o build/deptakeover-macos-arm64 ./cmd/deptakeover

echo "✅ Build complete!"
echo ""
echo "Binaries available in ./build/"
echo "  - deptakeover-windows-amd64.exe"
echo "  - deptakeover-linux-amd64"
echo "  - deptakeover-macos-amd64"
echo "  - deptakeover-macos-arm64"
