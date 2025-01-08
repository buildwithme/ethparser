package blockfetch

import (
	"context"
	"sync"

	"github.com/buildwithme/ethparser/internal/rpcfetch"
	"github.com/buildwithme/ethparser/pkg/constants"
	"github.com/buildwithme/ethparser/pkg/env"
	"github.com/buildwithme/ethparser/pkg/logger"
	"github.com/buildwithme/ethparser/pkg/storage"
)

type BlockFetch interface {
	// Run starts the block fetching process.
	Run()
	// GetCurrentBlock returns the highest block we've successfully written to Storage.
	GetCurrentBlock() int
	// ProcessRange fetches and processes blocks from `start` to `end`.
	ProcessRange(ctx context.Context, start, end int) error
}

type blockFetcher struct {
	concurrency   int
	chunkSize     int
	maxRetries    int
	log           *logger.Logger
	storage       storage.Storage
	mu            sync.RWMutex
	lastProcessed int
	rpcFetcher    rpcfetch.Fetcher
}

// NewFetcher constructs a blockFetcher.
func NewFetcher(log *logger.Logger, sto storage.Storage, rpcFetcher rpcfetch.Fetcher) BlockFetch {
	concurrency := env.GetEnvInt(constants.ENV_CONCURRENCY, 1)
	chunkSize := env.GetEnvInt(constants.ENV_CHUNK_SIZE, 50)
	maxRetries := env.GetEnvInt(constants.ENV_MAX_RETRIES, 3)

	return &blockFetcher{
		log:         log,
		storage:     sto,
		rpcFetcher:  rpcFetcher,
		concurrency: concurrency,
		chunkSize:   chunkSize,
		maxRetries:  maxRetries,
	}
}
