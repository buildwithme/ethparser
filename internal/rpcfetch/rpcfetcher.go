package rpcfetch

import (
	"context"

	"github.com/buildwithme/ethparser/pkg/env"
	"github.com/buildwithme/ethparser/pkg/logger"
)

type (
	// BlockResult captures the outcome of processing a single block.
	BlockResult struct {
		BlockNumber  int
		Transactions []*BlockTransaction
		Err          error
	}

	// BlockTransaction is just a minimal representation before mapping to storage.Transaction.
	BlockTransaction struct {
		Hash        string
		From        string
		To          string
		BlockNumber int
		Value       string
	}

	// Fetcher is the interface for fetching blocks from an Ethereum node.
	Fetcher interface {
		// GetLatestBlock returns the latest block number from the endpoint.
		GetLatestBlock(ctx context.Context) (int, error)
		// FetchBlock returns a block result from the endpoint.
		FetchBlock(ctx context.Context, blockNum int) (*BlockResult, error)
	}

	// ethFetcher is the implementation of Fetcher.
	ethFetcher struct {
		log      *logger.Logger
		endpoint string
	}
)

// NewFetcher constructs an ethFetcher that fetches blocks from `endpoint`.
func NewFetcher(log *logger.Logger) Fetcher {
	endpoint := env.GetEnvString("RPC_ENDPOINT", "https://cloudflare-eth.com")

	return &ethFetcher{
		log:      log,
		endpoint: endpoint,
	}
}
