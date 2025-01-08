# Ethereum Parser Project

This repository contains two main Go binaries:

- **CLI (ethcli)** – A command-line tool for parsing Ethereum blocks and fetching transactions.
- **HTTP Server (ethserver)** – An HTTP API providing endpoints to manage subscriptions and retrieve block/transaction data.

## Features

- **Environment-based Configuration**: Load defaults from a `.env` file (addresses, concurrency, etc.).
- **Concurrency**: Process blocks in parallel with a worker pool.
- **Pluggable Storage**: Use in-memory or switch to a database (e.g., PostgreSQL) to store transactions.
- **CLI & HTTP**: Choose either the command-line interface or a REST API for integration.

## Project Structure

```
.
├── cmd
│   ├── ethcli       // CLI entrypoint
│   │   └── main.go
│   └── ethserver    // HTTP server entrypoint
│       └── main.go
├── internal         // Internal packages (blockfetch, rpcfetch, httphandlers, storage, etc.)
├── pkg              // Reusable packages (env, logger, constants etc.)
├── .env             // Environment variables (optional)
├── Makefile         // Builds and runs the project
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- **Go** (1.22+ recommended)

### Installation

1. Clone this repo:

```bash
git clone https://github.com/buildwithme/ethparser.git
cd ethparser
```

2. Build the binaries using the provided Makefile:

```bash
make build
```

This will create the binaries `ethcli` and `ethserver` in the `bin` directory.

### Configuration

#### `.env` (optional)

Place environment variables here. For example:

```ini
# .env

# Comma-separated addresses to watch:
ADDRESSES=0x0000000000000000000000000000000000000000,0x0000000000000000000000000000000000000000

# RPC Endpoint for ethereum blockchain
RPC_ENDPOINT=https://cloudflare-eth.com

# Default concurrency for block fetching:
CONCURRENCY=4

# Chunk size (blocks to process per chunk):
CHUNK_SIZE=50

# Max number of retries for transient RPC errors:
MAX_RETRIES=3

# Default start block to watch:
DEFAULT_START_BLOCK=-1

# Default end block to watch:
DEFAULT_END_BLOCK=-1

# Port to run the server on:
PORT=3000
```

#### CLI Flags

CLI flags can override `.env`. For instance:

```bash
./bin/ethcli --concurrency=8 --addresses=0xDEADBEEF...
```

## Usage

### CLI (ethcli)

Run the CLI with default `.env` settings:

```bash
make run-cli
```

Override concurrency or addresses:

```bash
./bin/ethcli --concurrency=8 --addresses=0x123,0x456
```

Help:

```bash
./bin/ethcli --help
```

### HTTP Server (ethserver)

Run the server with default `.env` settings:

```bash
make run-server
```

Specify custom address:

```bash
./bin/ethserver -port 9090
```

#### Endpoints (examples):

- **POST /subscribe?address=0x1234** → Adds an address.
- **GET /transactions?address=0x1234** → Returns all transactions for that address.
- **GET /current-block** → Shows the last processed block.

## Cleaning Up

Remove compiled binaries:

```bash
make clean
```

## Contributing

1. Fork the repo and create a new branch for your feature/fix.
2. Commit changes with clear messages.
3. Pull request into the `main` branch.

## License

[MIT](https://choosealicense.com/licenses/mit/)
