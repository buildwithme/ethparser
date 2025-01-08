APP_NAME_CLI := ethcli

.PHONY: help all build run-cli clean

help:
	@echo "Usage:"
	@echo "  make build         # Build the binary"
	@echo "  make run           # Run the binary (with .env config)"
	@echo "  ./bin/$(APP_NAME) --help   # Print CLI usage"
	@echo
	@echo "Examples:"
	@echo "  ./bin/$(APP_NAME) --concurrency 8 --addresses 0x123,0x456"

all: build

build:
	@echo "===> Building CLI..."
	go build -o bin/$(APP_NAME_CLI) ./cmd/ethcli

run-cli: build
	@echo "===> Running CLI..."
	./bin/$(APP_NAME_CLI)

clean:
	@echo "===> Cleaning binary"
	rm -rf bin