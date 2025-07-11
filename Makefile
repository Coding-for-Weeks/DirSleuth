# Optional: allow overriding the binary name from CLI
BINARY_NAME ?= DirSleuth
SRC=cmd/dirsleuth/main.go

.PHONY: all build run clean test lint help

all: build

build:
	@echo "🔨 Building the project..."
	go build -o $(BINARY_NAME) $(SRC)

run: build
	@echo "🚀 Running the project..."
	@if [ ! -f $(BINARY_NAME) ]; then \
		echo "❌ Error: Binary file $(BINARY_NAME) not found. Please build the project first."; \
		exit 1; \
	fi
	./$(BINARY_NAME) -d example.com -w wordlist.txt -t 20 -timeout 30 -user-agent "DirSleuth/2.0" -status "200,301,403" -output json

clean:
	@echo "🧹 Cleaning up..."
	@if [ -f $(BINARY_NAME) ]; then \
		rm -f $(BINARY_NAME); \
	fi
	@rm -rf temp logs cache

test:
	@echo "🧪 Running tests..."
	go test ./... $(TEST_FLAGS)

lint:
	@echo "🔎 Running linter..."
	@golangci-lint run || echo "Install golangci-lint for linting support."

help:
	@echo "Available targets:"
	@echo "  build    Build the binary"
	@echo "  run      Build and run with example flags"
	@echo "  test     Run Go tests"
	@echo "  lint     Run Go linter (golangci-lint)"
	@echo "  clean    Remove binary and temp files"
