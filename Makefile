# Makefile for go-module project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint
GOIMPORTS=goimports

# Version (fallback to dev when git metadata is unavailable)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

# Binary name
BINARY_NAME=gocraft

# Where go install places binaries
GOBIN_DIR := $(shell go env GOBIN)
GOPATH_BIN := $(shell go env GOPATH)/bin
INSTALL_DIR := $(if $(GOBIN_DIR),$(GOBIN_DIR),$(GOPATH_BIN))


# Build directory
BUILD_DIR=build

# Main package path
MAIN_PACKAGE=./cmd/gocraft

.PHONY: all build install uninstall run test test-coverage clean lint deps verify help goimports

all: test goimports fmt build

# Build the project
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags="-s -w -X github.com/nduyhai/gocraft/pkg/version.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run linter
lint:
	$(GOLINT) run

# Format code
fmt:
	$(GOCMD) fmt ./...

# Run goimports
goimports:
	@which $(GOIMPORTS) > /dev/null || go install golang.org/x/tools/cmd/goimports@latest
	$(GOIMPORTS) -w ./

# Verify dependencies
verify:
	$(GOMOD) verify

run:
	$(GOCMD) run $(MAIN_PACKAGE)

# Install the binary into GOPATH/bin or GOBIN
install:
	$(GOCMD) install -ldflags="-s -w -X github.com/nduyhai/gocraft/pkg/version.Version=$(VERSION)" $(MAIN_PACKAGE)

# Uninstall the binary from GOBIN or GOPATH/bin
uninstall:
	@if [ -z "$(INSTALL_DIR)" ]; then echo "Cannot determine install dir"; exit 1; fi
	@echo "Removing $(INSTALL_DIR)/$(BINARY_NAME)"
	@rm -f "$(INSTALL_DIR)/$(BINARY_NAME)"

# Show help
help:
	@echo "Make targets:"
	@echo "  all            - Run tests and build"
	@echo "  build          - Build the gocraft binary"
	@echo "  install        - Install the gocraft binary to GOBIN/GOPATH/bin"
	@echo "  uninstall      - Remove the installed gocraft binary from GOBIN/GOPATH/bin"
	@echo "  run            - Run the gocraft CLI"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  goimports      - Run goimports to format code and update imports"
	@echo "  verify         - Verify dependencies"
	@echo "  help           - Show this help"
