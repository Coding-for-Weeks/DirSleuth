BINARY_NAME=$(if $(CUSTOM_BINARY_NAME),$(CUSTOM_BINARY_NAME),DirSleuth)
SRC=cmd/main.go

.PHONY: all build run clean test

all: build

build:
	@echo "Building the project..."
	go build -o $(BINARY_NAME) $(SRC)

run: build
	@echo "Running the project..."
	@if [ ! -f $(BINARY_NAME) ]; then \
	echo "Error: Binary file $(BINARY_NAME) not found. Please build the project first."; \
	exit 1; \
	fi
	./$(BINARY_NAME)

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