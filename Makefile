# Makefile for resume-analyzer

# Binary name
BINARY_NAME=resume-analyzer

# Build directory
BUILD_DIR=bin

# Subfolder support - can be set via: make subfolder=myfolder convert-pdfs
SUBFOLDER ?=

# Define paths with subfolder support
INPUT_PDFS_DIR = input_pdfs$(if $(SUBFOLDER),/$(SUBFOLDER))
OUTPUT_TXTS_DIR = output_txts$(if $(SUBFOLDER),/$(SUBFOLDER))
OUTPUT_SUMMARIES_DIR = output_summaries$(if $(SUBFOLDER),/$(SUBFOLDER))
OUTPUT_CONSOLIDATED_DIR = output_consolidated$(if $(SUBFOLDER),/$(SUBFOLDER))

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
	@echo "Input: $(INPUT_PDFS_DIR), Output: $(OUTPUT_TXTS_DIR)"
	@mkdir -p $(OUTPUT_TXTS_DIR)
	./$(BUILD_DIR)/$(BINARY_NAME) convert-pdfs -i $(INPUT_PDFS_DIR) -o $(OUTPUT_TXTS_DIR)

# Run the summarize command
summarize: build
	@echo "Running summarize example..."
	@echo "Input: $(OUTPUT_TXTS_DIR), Output: $(OUTPUT_SUMMARIES_DIR)"
	@mkdir -p $(OUTPUT_SUMMARIES_DIR)
	./$(BUILD_DIR)/$(BINARY_NAME) summarize -i $(OUTPUT_TXTS_DIR) -o $(OUTPUT_SUMMARIES_DIR)

# Run the consolidate command
consolidate: build
	@echo "Running consolidate example..."
	@echo "Input: $(OUTPUT_SUMMARIES_DIR), Output: $(OUTPUT_CONSOLIDATED_DIR)"
	@mkdir -p $(OUTPUT_CONSOLIDATED_DIR)
	./$(BUILD_DIR)/$(BINARY_NAME) consolidate -i $(OUTPUT_SUMMARIES_DIR) -o $(OUTPUT_CONSOLIDATED_DIR)/consolidated_table_$(shell date +%Y%m%d_%H%M%S).csv

# Run the query command (example)
query: build
	@echo "Running query example..."
	@echo "Example: ./$(BUILD_DIR)/$(BINARY_NAME) query -p 'Who has the most experience with Python?' -i $(OUTPUT_TXTS_DIR)"
	@echo "Example: ./$(BUILD_DIR)/$(BINARY_NAME) query -p 'Compare the technical skills of all candidates' -i $(OUTPUT_TXTS_DIR)"
	@echo "Note: Responses are printed to console by default. Use -o filename to save to file."

# Master workflow: convert PDFs, summarize, and consolidate
all-steps: build
	@echo "=== Starting complete resume analysis workflow ==="
	@if [ -n "$(SUBFOLDER)" ]; then echo "Using subfolder: $(SUBFOLDER)"; fi
	@echo "Step 1: Converting PDFs to text..."
	@$(MAKE) convert-pdfs SUBFOLDER=$(SUBFOLDER)
	@echo "Step 2: Generating summaries..."
	@$(MAKE) summarize SUBFOLDER=$(SUBFOLDER)
	@echo "Step 3: Creating consolidated table..."
	@$(MAKE) consolidate SUBFOLDER=$(SUBFOLDER)
	@echo "=== Workflow complete! ==="

# Clean output directories
clean-outputs:
	@echo "Cleaning output directories..."
	@if [ -n "$(SUBFOLDER)" ]; then \
		echo "Cleaning subfolder: $(SUBFOLDER)"; \
		rm -rf $(OUTPUT_TXTS_DIR)/*; \
		rm -rf $(OUTPUT_SUMMARIES_DIR)/*; \
		rm -rf $(OUTPUT_CONSOLIDATED_DIR)/*; \
	else \
		rm -rf output_txts/*; \
		rm -rf output_summaries/*; \
		rm -rf output_consolidated/*; \
	fi
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
	@echo ""
	@echo "Subfolder usage:"
	@echo "  make SUBFOLDER=myfolder convert-pdfs    - Use input_pdfs/myfolder and output_txts/myfolder"
	@echo "  make SUBFOLDER=myfolder all-steps       - Run complete workflow with subfolder"
	@echo "  make SUBFOLDER=myfolder clean-outputs   - Clean only subfolder outputs" 