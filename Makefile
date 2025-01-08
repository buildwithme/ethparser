APP_NAME_CLI := ethcli
APP_NAME_SERVER := ethserver

.PHONY: help all build run-cli run-server clean

help:
	@echo "Usage:"
	@echo "  make build            # Build both CLI and server binaries"
	@echo "  make run-cli          # Build and run the CLI (with .env config)"
	@echo "  make run-server       # Build and run the HTTP server (with .env config)"
	@echo
	@echo "Examples:"
	@echo "  ./bin/$(APP_NAME_CLI) --help"
	@echo "  ./bin/$(APP_NAME_CLI) --concurrency=8 --addresses=0x123,0x456"
	@echo "  ./bin/$(APP_NAME_SERVER) -port=:9090"
	@echo

all: build

build:
	@echo "===> Building CLI..."
	go build -o bin/$(APP_NAME_CLI) ./cmd/ethcli
	@echo "===> Building Server..."
	go build -o bin/$(APP_NAME_SERVER) ./cmd/ethserver

run-cli: build
	@echo "===> Running CLI..."
	./bin/$(APP_NAME_CLI)

run-server: build
	@echo "===> Running Server..."
	./bin/$(APP_NAME_SERVER)

clean:
	@echo "===> Cleaning binaries..."
	rm -rf bin