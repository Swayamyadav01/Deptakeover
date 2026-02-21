# DepTakeover Makefile

# Build configuration
BINARY_NAME=deptakeover
VERSION=$(shell git describe --tags --always --dirty)
BUILD_DIR=build
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")

# Default target
.PHONY: all
all: clean test build

# Development targets
.PHONY: dev
dev: 
	@echo "🔧 Building development version..."
	go build -o $(BINARY_NAME) ./cmd/deptakeover

.PHONY: run
run: dev
	@echo "🚀 Running DepTakeover..."
	./$(BINARY_NAME)

.PHONY: test
test:
	@echo "🧪 Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: lint
lint:
	@echo "🔍 Running linters..."
	go vet ./...
	go fmt ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed, skipping advanced linting"; \
	fi

.PHONY: fmt
fmt:
	@echo "💫 Formatting code..."
	go fmt ./...

# Build targets
.PHONY: build
build: clean
	@echo "📦 Building for current platform..."
	mkdir -p $(BUILD_DIR)
	go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/deptakeover

.PHONY: build-all
build-all: clean
	@echo "🌍 Building for all platforms..."
	./build.sh

.PHONY: build-linux
build-linux: clean
	@echo "🐧 Building for Linux..."
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/deptakeover

.PHONY: build-windows
build-windows: clean
	@echo "🪟 Building for Windows..."
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/deptakeover

.PHONY: build-mac
build-mac: clean
	@echo "🍎 Building for macOS..."
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/deptakeover
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/deptakeover

# Release targets
.PHONY: release
release: test lint build-all
	@echo "🎉 Creating release..."
	./release.sh $(VERSION)

# Clean targets
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

.PHONY: clean-all
clean-all: clean
	@echo "🧹 Cleaning all generated files..."
	rm -rf .github_repos/
	rm -f *_report.json

# Dependencies
.PHONY: deps
deps:
	@echo "📥 Installing dependencies..."
	go mod download
	go mod tidy

.PHONY: deps-update
deps-update:
	@echo "⬆️  Updating dependencies..."
	go get -u ./...
	go mod tidy

# Development tools
.PHONY: install-tools
install-tools:
	@echo "🔧 Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/sonatypecommunity/nancy@latest

# Documentation
.PHONY: docs
docs:
	@echo "📚 Generating documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		echo "Starting godoc server at http://localhost:6060/pkg/deptakeover/"; \
		godoc -http=:6060; \
	else \
		echo "godoc not installed. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# Security
.PHONY: security
security:
	@echo "🔒 Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Quick test commands
.PHONY: test-npm
test-npm: dev
	@echo "🧪 Testing npm scanning..."
	./$(BINARY_NAME) npm lodash/lodash

.PHONY: test-pypi  
test-pypi: dev
	@echo "🧪 Testing PyPI scanning..."
	./$(BINARY_NAME) pypi django/django

.PHONY: test-composer
test-composer: dev
	@echo "🧪 Testing Composer scanning..."
	./$(BINARY_NAME) composer laravel/laravel

.PHONY: test-org
test-org: dev
	@echo "🧪 Testing organization scanning (small org)..."
	./$(BINARY_NAME) org-npm vercel

# Help
.PHONY: help
help:
	@echo "DepTakeover Build System"
	@echo "======================="
	@echo ""
	@echo "Development:"
	@echo "  dev           Build development version"
	@echo "  run           Build and run tool"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo "  lint          Run linters"
	@echo "  fmt           Format code"
	@echo ""
	@echo "Building:"
	@echo "  build         Build for current platform"
	@echo "  build-all     Build for all platforms"
	@echo "  build-linux   Build for Linux"
	@echo "  build-windows Build for Windows"
	@echo "  build-mac     Build for macOS"
	@echo ""
	@echo "Release:"
	@echo "  release       Create release build with checksums"
	@echo ""
	@echo "Maintenance:"
	@echo "  clean         Clean build artifacts"
	@echo "  clean-all     Clean all generated files"
	@echo "  deps          Install dependencies"
	@echo "  deps-update   Update dependencies"
	@echo ""
	@echo "Tools:"
	@echo "  install-tools Install development tools"
	@echo "  docs          Start documentation server"
	@echo "  security      Run security checks"
	@echo ""
	@echo "Quick Tests:"
	@echo "  test-npm      Test npm scanning"
	@echo "  test-pypi     Test PyPI scanning"  
	@echo "  test-composer Test Composer scanning"
	@echo "  test-org      Test organization scanning"
