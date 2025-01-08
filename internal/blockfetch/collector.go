package blockfetch

import (
	"context"
	"time"
)

const BLOCK_SYNC_TIMEOUT = 2 * time.Second

func (p *blockFetcher) Run() {
	// Create a context we can cancel if needed (for graceful shutdown)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch the current chain tip
	latest, err := p.rpcFetcher.GetLatestBlock(ctx)
	if err != nil {
		p.log.Fatalf("Cannot fetch latest block: %v", err)
	}

	// Define how many blocks behind the tip to start
	startFlag := latest - 10
	if startFlag < 0 {
		startFlag = 0
	}

	// Process the initial range from (tip - 10) up to the tip
	err = p.ProcessRange(ctx, startFlag, latest)
	if err != nil {
		p.log.Fatalf("ProcessRange error: %v", err)
	}

	p.log.Printf("Initial catch-up done. Last processed = %d", p.GetCurrentBlock())

	// Loop indefinitely, always trying to catch up to the latest block
	for {
		// If the context is canceled (e.g., SIGTERM), exit gracefully
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Fetch the updated chain tip
		currentTip, err := p.rpcFetcher.GetLatestBlock(ctx)
		if err != nil {
			p.log.Printf("Failed to fetch latest block: %v", err)
			// Sleep briefly and retry
			time.Sleep(BLOCK_SYNC_TIMEOUT)
			continue
		}

		// Get the last block we've processed
		lastProcessed := p.GetCurrentBlock()

		// If we're behind, process the range
		if lastProcessed < currentTip {
			err = p.ProcessRange(ctx, lastProcessed+1, currentTip)
			if err != nil {
				p.log.Printf("ProcessRange error (blocks %d..%d): %v", lastProcessed+1, currentTip, err)
			} else {
				p.log.Printf("Updated. Last processed = %d", p.GetCurrentBlock())
			}
		}

		time.Sleep(BLOCK_SYNC_TIMEOUT)
	}
}
