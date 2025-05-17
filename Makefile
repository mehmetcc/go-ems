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

# Clean build artifacts
clean:
	@echo "Cleaning..."
	go clean
	rm -rf bin/

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Make all will build and run
all: build run

.PHONY: build run clean test all