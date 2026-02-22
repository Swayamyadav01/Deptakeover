#!/bin/bash

# DepTakeover Release Script
# This script builds cross-platform binaries for release

set -e

VERSION=${1:-$(git describe --tags --always --dirty)}
BUILD_DIR="build"
BINARY_NAME="deptakeover"

echo "🚀 Building DepTakeover v$VERSION"

# Clean build directory
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

# Build matrix
platforms=(
    "windows/amd64"
    "windows/arm64" 
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

for platform in "${platforms[@]}"; do
    os=${platform%/*}
    arch=${platform#*/}
    output_name=$BINARY_NAME-$os-$arch
    
    if [ "$os" = "windows" ]; then
        output_name+='.exe'
    fi
    
    echo "📦 Building for $os/$arch..."
    
    env GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
        -ldflags="-w -s -X main.version=$VERSION" \
        -o $BUILD_DIR/$output_name \
        ./cmd/deptakeover
    
    if [ $? -eq 0 ]; then
        echo "✅ Built $output_name"
    else
        echo "❌ Failed to build $output_name"
        exit 1
    fi
done

# Generate checksums
echo "🔐 Generating checksums..."
cd $BUILD_DIR
sha256sum * > checksums.txt

echo ""
echo "📊 Build Summary:"
echo "=================="
ls -lh
echo ""
echo "🎉 Release build complete!"
echo "📁 Binaries available in: $BUILD_DIR/"
echo "🔐 Checksums: $BUILD_DIR/checksums.txt"