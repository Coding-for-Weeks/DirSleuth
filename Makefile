# Makefile for GoBuster Clone

# Variables
# Allows the binary name to be customized using CUSTOM_BINARY_NAME; defaults to 'DirSleuth' if not specified
BINARY_NAME=$(if $(CUSTOM_BINARY_NAME),$(CUSTOM_BINARY_NAME),DirSleuth)
SRC=main.go
CONFIG?=config.json

.PHONY: all build run clean test

# Default target
all: build

# Build the project
build:
	@echo "Building the project..."
	go build -o $(BINARY_NAME) $(SRC)

# Run the project with configurable config file
run: build
	@echo "Running the project..."
	@if [ ! -f $(BINARY_NAME) ]; then \
		echo "Error: Binary file $(BINARY_NAME) not found. Please build the project first."; \
		exit 1; \
	fi
	@if [ ! -f $(CONFIG) ]; then \
		echo "Error: Configuration file $(CONFIG) not found."; \
		exit 1; \
	fi
	./$(BINARY_NAME) -config $(CONFIG)

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@if [ -f $(BINARY_NAME) ]; then \
		rm -f $(BINARY_NAME); \
	fi
	@if [ -d temp ]; then \
		rm -rf temp; \
	fi
	@if [ -d logs ]; then \
		rm -rf logs; \
	fi
	@if [ -d cache ]; then \
		rm -rf cache; \
	fi

# Run tests
test:
	@echo "Running tests..."
	go test ./... $(TEST_FLAGS)

# Run all including tests
testall: test build

# Example for running specific tests or enabling verbose output
# Usage:
# make test TEST_FLAGS="-run TestSpecificFunction -v"
# make test TEST_FLAGS="-bench ."  # Run benchmark tests
# make test TEST_FLAGS="-count=1"  # Disable test result caching
# make test TEST_FLAGS="-cover"  # Generate test coverage report
