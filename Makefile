# Binary name
BINARY_NAME=go-ems

# Build the application
build:
	@echo "Building..."
	go build -o bin/$(BINARY_NAME) ./cmd/main.go

# Run the application
run:
	@echo "Running..."
	go run ./cmd/main.go

# Clean build artifacts, coverage reports, and output files
clean:
	@echo "Cleaning..."
	go clean
	rm -rf bin/ coverage/ output/

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Test with coverage
coverage:
	@echo "Generating coverage report..."
	@mkdir -p coverage
	@go test -coverprofile=coverage/coverage.out ./...
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@go tool cover -func=coverage/coverage.out
	@echo "\nCoverage report generated at coverage/coverage.html"

# Make all will build and run
all: build run

.PHONY: build run clean test coverage all