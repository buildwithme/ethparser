package main

import (
	"github.com/buildwithme/ethparser/internal/blockfetch"
	"github.com/buildwithme/ethparser/internal/parser"
	"github.com/buildwithme/ethparser/internal/rpcfetch"
	"github.com/buildwithme/ethparser/pkg/env"
	"github.com/buildwithme/ethparser/pkg/logger"
	"github.com/buildwithme/ethparser/pkg/storage"
)

func main() {
	logger := logger.NewLogger()

	// Parse CLI flags & override .env if needed
	cf := ParseFlags()

	cf.ApplyEnvFile()

	if err := env.LoadDotEnv(); err != nil {
		logger.Fatalf("[FATAL] .env not loaded: %v", err)
	}

	// Apply any CLI flag overrides to the environment
	cf.ApplyConfig()

	storage := storage.NewMemoryStorage()
	rpcFetcher := rpcfetch.NewFetcher(logger)
	blockFetcher := blockfetch.NewFetcher(logger, storage, rpcFetcher)
	parser := parser.NewParser(logger, storage, blockFetcher)

	// Subscribe addresses
	subscribeEnvAddresses(parser, logger)

	// Fetch blocks
	blockFetcher.Run()
}
