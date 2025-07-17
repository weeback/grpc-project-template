# Makefile for Application
.DEFAULT_GOAL := help
.PHONY: build clean test run build-docker help

# ====================================================================================
# Variables
# ====================================================================================
SHELL = /usr/bin/env bash

# Application settings
PACKAGE_NAME := "github/weeback/grpc-project-template/pkg"
MAIN_GO_FILE ?= cmd/v1/*.go

# Build artifacts
BIN_DIR := $(shell pwd)/bin
BINARY  ?= $(BIN_DIR)/app

# Build-time variables
# Use ?= to replace with command line argument (exp: make build VERSION=1.2.3)
VERSION      ?= $(shell git describe --tags --always --dirty)
COMMIT_HASH  := $(shell git rev-parse HEAD)
BRANCH       := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE   := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_USER   := $(shell whoami)@$(shell hostname)
REPO_URL     := $(shell git remote get-url origin 2>/dev/null || echo \
	"https://github.com/weeback/grpc-project-template.git")

# Go build flags
LDFLAGS := -ldflags="\
	-s -w \
	-X '$(PACKAGE_NAME).Version=$(VERSION)' \
	-X '$(PACKAGE_NAME).BuildCommit=$(COMMIT_HASH)' \
	-X '$(PACKAGE_NAME).BuildBranch=$(BRANCH)' \
	-X '$(PACKAGE_NAME).BuildDate=$(BUILD_DATE)' \
	-X '$(PACKAGE_NAME).BuildUser=$(BUILD_USER)' \
	-X '$(PACKAGE_NAME).RepoURL=$(REPO_URL)' \
"

# Build application binary
build:
	@echo "==> Building application..."
	@mkdir -p $(BIN_DIR)
	@export CGO_ENABLED=0; \
		go build -v $(LDFLAGS) -o $(BINARY) -trimpath $(MAIN_GO_FILE)
	@echo "==> Build successful: $(BINARY)"

# Run the application (sẽ build nếu cần)
run: build
	@echo "==> Running application..."
	@$(BINARY)

# Run all tests với race detector
test:
	@echo "==> Running tests..."
	@go test -v -race ./...

# Clean up build artifacts
clean:
	@echo "==> Cleaning up..."
	@rm -rf $(BIN_DIR)

# Build the Docker image (placeholder)
build-docker:
	@echo "Building Docker image... (not implemented yet)"

# Defined help command
help:
	@echo "Available commands:"
	@echo "  make build        : Build the application binary."
	@echo "  make run          : Build and run the application."
	@echo "  make test         : Run all unit tests with race detector."
	@echo "  make clean        : Remove build artifacts."
	@echo "  make build-docker : Build the Docker image (TBD)."
	@echo "  make help         : Show this help message."

init:
	@echo "==> Initializing Go module..."
	@rm -rf go.mod go.sum
	@echo 'module github/weeback/grpc-project-template' > go.mod
	@echo '' >> go.mod
	@echo 'go 1.24' >> go.mod
	@echo "==> Adding default replace directive..."
	@echo '' >> go.mod
	@echo '' >> go.mod
	@echo "==> Initialization completed."
	@go mod tidy; \
		go get -u=patch ./...; \
		go get -u all; \
		go mod tidy
	@echo "Successful!"

# ====================================================================================
# gRPC generate code
# ====================================================================================

# Path Generating gRPC reflection code
PROTO_PATH = proto
PROTO_GEN = pb
PROTO_PATH_HELLO = $(PROTO_PATH)

grpc-generate: grpc-force
	@echo "==> Generating gRPC code..."
	@protoc --proto_path=$(PROTO_PATH) \
		--go_out=$(PROTO_GEN) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GEN) --go-grpc_opt=paths=source_relative \
		proto/**/*.proto || { echo "Failed to generate gRPC code"; exit 1; }
		@echo "==> gRPC code generated successful"
		@tree $(PROTO_GEN) | echo "Done!"

grpc-force:
	@echo "==> Creating output directory..."
	@mkdir -p $(PROTO_GEN)
	@echo "==> Cleaning up old generated files..."
	@rm -rf $(PROTO_GEN)/*.go