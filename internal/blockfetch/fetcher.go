package blockfetch

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/buildwithme/ethparser/internal/rpcfetch"
	"github.com/buildwithme/ethparser/pkg/storage"
)

// GetCurrentBlock returns the highest block we've successfully written to Storage.
func (p *blockFetcher) GetCurrentBlock() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.lastProcessed
}

// ProcessRange fetches blocks [start..end], chunking them to reduce memory overhead.
func (p *blockFetcher) ProcessRange(ctx context.Context, start, end int) error {
	if start > end {
		return fmt.Errorf("start block (%d) > end block (%d)", start, end)
	}

	// Divide the range into chunks (e.g., 100 blocks per chunk).
	for chunkStart := start; chunkStart <= end; chunkStart += p.chunkSize {
		chunkEnd := chunkStart + p.chunkSize - 1
		if chunkEnd > end {
			chunkEnd = end
		}

		// Prepare a worker pool for this chunk
		blocks := make([]int, 0, chunkEnd-chunkStart+1)
		for b := chunkStart; b <= chunkEnd; b++ {
			blocks = append(blocks, b)
		}

		err := p.processChunk(ctx, blocks)
		if err != nil {
			// We might log and continue or abort. We'll abort here to keep it simple.
			p.log.Printf("Error processing chunk [%d..%d]: %v", chunkStart, chunkEnd, err)
			return err
		}

		// Update checkpoint
		p.mu.Lock()
		p.lastProcessed = chunkEnd
		p.mu.Unlock()

		p.log.Printf("Chunk [%d..%d] done; lastProcessed=%d",
			chunkStart, chunkEnd, p.lastProcessed)
	}

	return nil
}

// processChunk sets up a worker pool to fetch the blocks concurrently.
func (p *blockFetcher) processChunk(ctx context.Context, blocks []int) error {
	wp := NewWorkerPool(p.concurrency)
	results := wp.Run(ctx, blocks, p.fetchBlockWithRetry)

	// We'll gather results in memory, sort them by block, then store them.
	var successful []*rpcfetch.BlockResult
	for result := range results {
		if result.Err != nil {
			// Partial failure approach: log the error, skip that block.
			p.log.Printf("[WARN] Block %d failed after max retries: %v", result.BlockNumber, result.Err)
			continue
		}
		successful = append(successful, result)
	}

	// Sort ascending by block number
	for i := 0; i < len(successful)-1; i++ {
		for j := i + 1; j < len(successful); j++ {
			if successful[i].BlockNumber > successful[j].BlockNumber {
				successful[i], successful[j] = successful[j], successful[i]
			}
		}
	}

	// Insert in ascending order
	for _, s := range successful {
		var txs []storage.Transaction
		for _, t := range s.Transactions {
			tx := storage.Transaction{
				Hash:        t.Hash,
				From:        t.From,
				To:          t.To,
				BlockNumber: t.BlockNumber,
				Value:       t.Value,
			}

			txs = append(txs, tx)
		}

		if err := p.storage.StoreBlockTransactions(s.BlockNumber, txs); err != nil {
			return fmt.Errorf("store block %d error: %w", s.BlockNumber, err)
		}
	}
	return nil
}

// fetchBlockWithRetry wraps `fetchBlock` with exponential backoff retries.
func (p *blockFetcher) fetchBlockWithRetry(ctx context.Context, blockNum int) (*rpcfetch.BlockResult, error) {
	var lastErr error

	for attempt := 1; attempt <= p.maxRetries; attempt++ {
		result, err := p.rpcFetcher.FetchBlock(ctx, blockNum)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// Exponential backoff
		backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		p.log.Printf("[ERROR] block %d attempt %d failed: %v. Retrying in %v",
			blockNum, attempt, err, backoff)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
			continue
		}
	}

	return nil, fmt.Errorf("max retries reached: %w", lastErr)
}
