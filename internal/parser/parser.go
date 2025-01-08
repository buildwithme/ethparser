package parser

import (
	"github.com/buildwithme/ethparser/internal/blockfetch"
	"github.com/buildwithme/ethparser/internal/storage"
	"github.com/buildwithme/ethparser/pkg/logger"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []storage.Transaction
}

type ethParser struct {
	log          *logger.Logger
	storage      storage.Storage
	blockFetcher blockfetch.BlockFetch
}

// NewParser constructs an ethParser that fetches blocks from `endpoint`.
func NewParser(log *logger.Logger, sto storage.Storage, blockFetcher blockfetch.BlockFetch) Parser {
	return &ethParser{
		log:          log,
		storage:      sto,
		blockFetcher: blockFetcher,
	}
}

// GetCurrentBlock returns the highest block we've successfully written to Storage.
func (p *ethParser) GetCurrentBlock() int {
	return p.blockFetcher.GetCurrentBlock()
}

// Subscribe proxies to the storage layer.
func (p *ethParser) Subscribe(address string) bool {
	return p.storage.SubscribeAddress(address)
}

// GetTransactions gets transactions for a specific address.
func (p *ethParser) GetTransactions(address string) []storage.Transaction {
	return p.storage.GetTransactions(address)
}
