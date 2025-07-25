# TermUp - S3 compatible filesharing from terminal
# Build variables
BINARY_NAME=upl
MAIN_PATH=./cmd/upl
BUILD_DIR=build
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GOVERSION ?= $(shell go version | cut -d' ' -f3)

# Linker flags
LDFLAGS=-ldflags "\
	-X github.com/nizar0x1f/termup/pkg/version.Version=$(VERSION) \
	-X github.com/nizar0x1f/termup/pkg/version.Commit=$(COMMIT) \
	-X github.com/nizar0x1f/termup/pkg/version.Date=$(DATE) \
	-X github.com/nizar0x1f/termup/pkg/version.GoVersion=$(GOVERSION) \
	-s -w"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build for development (no version info)
.PHONY: build-dev
build-dev:
	@echo "Building $(BINARY_NAME) for development..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Install the binary
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME) $(VERSION)..."
	go install $(LDFLAGS) $(MAIN_PATH)

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe $(MAIN_PATH)

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Create a new release
.PHONY: release
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make release VERSION=v1.2.0"; \
		exit 1; \
	fi
	@./scripts/release.sh $(VERSION)

# Run go mod tidy
.PHONY: tidy
tidy:
	go mod tidy

# Show version information
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"
	@echo "Go Version: $(GOVERSION)"

# Development workflow
.PHONY: dev
dev: tidy fmt build-dev

# Release workflow
.PHONY: release
release: tidy fmt test build-all

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        Build the binary with version info"
	@echo "  build-dev    Build for development (no version info)"
	@echo "  build-all    Build for all platforms"
	@echo "  install      Install the binary with version info"
	@echo "  test         Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo "  clean        Clean build artifacts"
	@echo "  fmt          Format code"
	@echo "  lint         Lint code"
	@echo "  tidy         Run go mod tidy"
	@echo "  version      Show version information"
	@echo "  dev          Development workflow (tidy, fmt, build-dev)"
	@echo "  release      Release workflow (tidy, fmt, test, build-all)"
	@echo "  help         Show this help message"
