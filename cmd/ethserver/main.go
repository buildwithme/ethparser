package main

import (
	"fmt"
	"net/http"

	"github.com/buildwithme/ethparser/internal/blockfetch"
	"github.com/buildwithme/ethparser/internal/httphandlers"
	"github.com/buildwithme/ethparser/internal/parser"
	"github.com/buildwithme/ethparser/internal/rpcfetch"
	"github.com/buildwithme/ethparser/internal/storage"
	"github.com/buildwithme/ethparser/pkg/constants"
	"github.com/buildwithme/ethparser/pkg/env"
	"github.com/buildwithme/ethparser/pkg/logger"
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

	// Fetch blocks
	go blockFetcher.Run()

	// Register HTTP handlers
	handlers := httphandlers.New(parser)
	handlers.RegisterHandlers()

	endpoint := fmt.Sprintf(":%s", env.GetEnvString(constants.ENV_PORT, "8080"))

	logger.Printf("[INFO]: HTTP server starting on %s", endpoint)

	err := http.ListenAndServe(endpoint, nil)
	if err != nil && err != http.ErrServerClosed {
		logger.Fatalf("server error: %v", err)
	}
}
