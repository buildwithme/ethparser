package main

import (
	"context"
	"os"
	"strings"

	"github.com/buildwithme/ethparser/internal/blockfetch"
	"github.com/buildwithme/ethparser/internal/parser"
	"github.com/buildwithme/ethparser/internal/rpcfetch"
	"github.com/buildwithme/ethparser/pkg/env"
	"github.com/buildwithme/ethparser/pkg/logger"
	"github.com/buildwithme/ethparser/pkg/storage"
)

func main() {
	logger := logger.NewLogger()

	if err := env.LoadDotEnv("../../.env"); err != nil {
		logger.Fatalf("[WARN] .env not loaded: %v", err)
	}

	storage := storage.NewMemoryStorage()
	rpcFetcher := rpcfetch.NewFetcher(logger)
	blockFetcher := blockfetch.NewFetcher(logger, storage, rpcFetcher)
	parser := parser.NewParser(logger, storage, blockFetcher)

	// Subscribe addresses
	subscribeEnvAddresses(parser, logger)

	// Fetch blocks
	fetchBlocks(logger, blockFetcher, rpcFetcher, parser)
}

// subscribeEnvAddresses pulls addresses from ADDRESSES in .env
func subscribeEnvAddresses(parser parser.Parser, log *logger.Logger) {
	addrs := os.Getenv("ADDRESSES")
	if addrs == "" {
		return
	}

	for _, a := range strings.Split(addrs, ",") {
		a = strings.TrimSpace(a)
		if a != "" {
			parser.Subscribe(a)
			log.Printf("[INFO] Subscribed env address: %s", a)
		}
	}
}

func fetchBlocks(log *logger.Logger, blockFetcher blockfetch.BlockFetch, rpcFetcher rpcfetch.Fetcher, parser parser.Parser) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	latest, err := rpcFetcher.GetLatestBlock(ctx)
	if err != nil {
		log.Fatalf("Cannot fetch latest block: %v", err)
	}

	startFlag := latest - 10
	endFlag := latest

	err = blockFetcher.ProcessRange(ctx, startFlag, endFlag)
	if err != nil {
		log.Fatalf("ProcessRange error: %v", err)
	}

	log.Printf("Done. Last processed = %d", parser.GetCurrentBlock())
}
