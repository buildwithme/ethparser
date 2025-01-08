APP_NAME_CLI := ethcli

.PHONY: all build run-cli clean

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
	