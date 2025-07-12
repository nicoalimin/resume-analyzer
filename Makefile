# Makefile for resume-analyzer

# Binary name
BINARY_NAME=resume-analyzer

# Build directory
BUILD_DIR=bin

# Go build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean test run help

# Default target
all: clean build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run the application (example with flags)
run:
	@echo "Running $(BINARY_NAME)..."
	go run . --help

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Install golangci-lint if not present
install-lint:
	@echo "Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run the convert-pdfs command
convert-pdfs: build
	@echo "Running convert-pdfs example..."
	./$(BUILD_DIR)/$(BINARY_NAME) convert-pdfs -i input_pdfs -o output_txts

# Run the summarize command
summarize: build
	@echo "Running summarize example..."
	./$(BUILD_DIR)/$(BINARY_NAME) summarize -i output_txts -o output_summaries

# Run the consolidate command
consolidate: build
	@echo "Running consolidate example..."
	./$(BUILD_DIR)/$(BINARY_NAME) consolidate -i output_summaries -o output_consolidated/consolidated_table_$(shell date +%Y%m%d_%H%M%S).csv

# Run the query command (example)
query: build
	@echo "Running query example..."
	@echo "Example: ./$(BUILD_DIR)/$(BINARY_NAME) query -p 'Who has the most experience with Python?' -i output_txts"
	@echo "Example: ./$(BUILD_DIR)/$(BINARY_NAME) query -p 'Compare the technical skills of all candidates' -i output_txts"
	@echo "Note: Responses are printed to console by default. Use -o filename to save to file."

# Master workflow: convert PDFs, summarize, and consolidate
all-steps: build
	@echo "=== Starting complete resume analysis workflow ==="
	@echo "Step 1: Converting PDFs to text..."
	@$(MAKE) convert-pdfs
	@echo "Step 2: Generating summaries..."
	@$(MAKE) summarize
	@echo "Step 3: Creating consolidated table..."
	@$(MAKE) consolidate
	@echo "=== Workflow complete! ==="

# Clean output directories
clean-outputs:
	@echo "Cleaning output directories..."
	@rm -rf output_txts/*
	@rm -rf output_summaries/*
	@echo "Output directories cleaned."

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-all     - Build for multiple platforms (Linux, macOS, Windows)"
	@echo "  clean         - Clean build artifacts"
	@echo "  clean-outputs - Clean all output directories (output_consolidated, output_txts, output_summaries)"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  deps          - Install dependencies"
	@echo "  run           - Run the application"
	@echo "  convert-pdfs  - Run convert-pdfs example (input_pdfs -> output_txts)"
	@echo "  summarize     - Run summarize example (output_txts -> output_summaries)"
	@echo "  consolidate   - Run consolidate example (output_summaries -> consolidated_table_YYYYMMDD_HHMMSS.csv)"
	@echo "  query         - Run query example"
	@echo "  all-steps     - Run complete workflow: convert-pdfs -> summarize -> consolidate"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  install-lint  - Install golangci-lint"
	@echo "  help          - Show this help" 