.PHONY: all build dev test clean install deps help run

# Build variables
BINARY_NAME=matrixpulse
OUTPUT_DIR=.
GO=go
MAIN_PATH=cmd/matrixpulse/main.go

# Version info
VERSION?=1.0.0
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)"

all: clean deps test build

## build: Build production binary
build:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	$(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✓ Build complete: $(OUTPUT_DIR)/$(BINARY_NAME)"

## dev: Build development binary with race detector
dev:
	@echo "Building development version..."
	$(GO) build -race -o $(OUTPUT_DIR)/$(BINARY_NAME)-dev $(MAIN_PATH)
	@echo "✓ Development build complete"

## test: Run all tests
test:
	@echo "Running tests..."
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out
	@echo "✓ Tests complete"

## bench: Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

## deps: Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod verify
	@echo "✓ Dependencies ready"

## install: Install binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(MAIN_PATH)
	@echo "✓ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(OUTPUT_DIR)/$(BINARY_NAME)
	@rm -f $(OUTPUT_DIR)/$(BINARY_NAME)-dev
	@rm -f $(OUTPUT_DIR)/$(BINARY_NAME).exe
	@rm -f $(OUTPUT_DIR)/$(BINARY_NAME)-dev.exe
	@rm -f coverage.out
	@rm -f *.log
	@echo "✓ Clean complete"

## run: Build and run
run: build
	@echo "Starting $(BINARY_NAME)..."
	@$(OUTPUT_DIR)/$(BINARY_NAME)

## lint: Run linters
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...
	@echo "✓ Lint complete"

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "✓ Format complete"

## tidy: Tidy go.mod
tidy:
	@echo "Tidying go.mod..."
	$(GO) mod tidy
	@echo "✓ Tidy complete"

## help: Show this help
help:
	@echo "MatrixPulse Build System (Desktop Application)"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'