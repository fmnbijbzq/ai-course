.PHONY: wire clean build run env

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=ai-course
WIRE_GEN_FILE=internal/wire_gen.go

# Get GOPATH and set wire binary path
GOPATH=$(shell go env GOPATH)
WIRE=$(GOPATH)/bin/wire

# Wire tool installation and generation
wire:
	@echo "Installing wire..."
	@go install github.com/google/wire/cmd/wire@latest
	@echo "Generating wire_gen.go..."
	@cd internal/wire && $(WIRE)
	@echo "Wire generation completed."

# Clean wire generated files
clean-wire:
	@echo "Cleaning wire generated files..."
	@rm -f $(WIRE_GEN_FILE)
	@echo "Wire files cleaned."

# Build the application
build: wire
	@echo "Building..."
	@$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go
	@echo "Build completed."

# Clean build files
clean:
	@echo "Cleaning build files..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@echo "Clean completed."

# Run the application
run: build
	@echo "Running application..."
	@./$(BINARY_NAME)

# Show environment info
env:
	@echo "Go environment information:"
	@echo "GOPATH: $(GOPATH)"
	@echo "Wire binary: $(WIRE)"
	@$(GOCMD) env

# Help command
help:
	@echo "Available commands:"
	@echo "  make wire        - Generate wire_gen.go file"
	@echo "  make clean-wire  - Remove wire generated files"
	@echo "  make build      - Build the application (includes wire generation)"
	@echo "  make clean      - Clean build files"
	@echo "  make run        - Build and run the application"
	@echo "  make env        - Show Go environment information"
	@echo "  make help       - Show this help message"