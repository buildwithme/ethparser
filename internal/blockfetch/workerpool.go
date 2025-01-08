package blockfetch

import (
	"context"
	"sync"

	"github.com/buildwithme/ethparser/internal/rpcfetch"
)

// WorkFunc is the signature for the job function each worker runs.
type WorkFunc func(ctx context.Context, blockNum int) (*rpcfetch.BlockResult, error)

// WorkerPool manages concurrency for processing blocks.
type WorkerPool struct {
	concurrency int
}

// NewWorkerPool with a given concurrency level.
func NewWorkerPool(concurrency int) *WorkerPool {
	if concurrency <= 0 {
		concurrency = 1
	}
	return &WorkerPool{concurrency: concurrency}
}

// Run spins up `concurrency` workers that consume from `blocks` and call `fn`.
// Returns a channel of results to be read by the caller.
func (wp *WorkerPool) Run(ctx context.Context, blocks []int, fn WorkFunc) <-chan *rpcfetch.BlockResult {
	results := make(chan *rpcfetch.BlockResult, len(blocks))

	go func() {
		defer close(results)

		var wg sync.WaitGroup
		sem := make(chan struct{}, wp.concurrency)

		for _, b := range blocks {
			wg.Add(1)
			sem <- struct{}{} // blocks until a slot is available

			go func(blockNum int) {
				defer wg.Done()
				defer func() { <-sem }() // release slot

				r, err := fn(ctx, blockNum)
				if err != nil {
					results <- &rpcfetch.BlockResult{BlockNumber: blockNum, Err: err}
					return
				}

				results <- r // success
			}(b)
		}

		wg.Wait()
	}()

	return results
}
