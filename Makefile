# Optional: allow overriding the binary name from CLI
BINARY_NAME ?= DirSleuth
SRC=cmd/dirsleuth/main.go

.PHONY: all build run clean test testall

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
	./$(BINARY_NAME)

clean:
	@echo "🧹 Cleaning up..."
	@if [ -f $(BINARY_NAME) ]; then \
		rm -f $(BINARY_NAME); \
	fi
	@rm -rf temp logs cache

test:
	@echo "🧪 Running tests..."
	go test ./... $(TEST_FLAGS)

testall: test build
