# Justfile for prompt-mcp project

# Default recipe that runs build
default: build

# Build the project
build:
    go build -o bin/prompt-mcp ./cmd/server

# Run the server
run: build
    ./bin/prompt-mcp --prompts-dir ./prompts

# Run tests
test:
    go test ./...

# Run tests with verbose output
test-verbose:
    go test -v ./...

# Format the code
fmt:
    go fmt ./...

# Lint the code (requires golangci-lint)
lint:
    golangci-lint run

# Tidy dependencies
tidy:
    go mod tidy

# Clean build artifacts
clean:
    rm -rf bin/

# Install dev dependencies
install-dev:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run development server with file watching (requires air)
dev:
    air

# Initialize air config
init-air:
    air init

# Show help
help:
    @just --list
