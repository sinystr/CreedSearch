package crawler

import (
	"context"
	"sync"
)

// CrawlPages is a function which will be called to
// crawl pages concurrently.
// Input:
// pages: the slice with page adresses
// done:  the channel to trigger the end of task processing and return
// Output: the channel with results
func CrawlPages(ctx context.Context, pages []string) chan CrawledPage {

	// Create a worker for each page
	workers := make([]chan CrawledPage, 0, len(pages))

	for _, page := range pages {
		resultChannel := newWorker(ctx, page)
		workers = append(workers, resultChannel)
	}

	// Merge results from all workers
	out := merge(ctx, workers)
	return out
}

func newWorker(ctx context.Context, page string) chan CrawledPage {
	out := make(chan CrawledPage)
	go func() {
		defer close(out)

		select {
		case <-ctx.Done():
			// Received a signal to abandon further processing
			return
		case out <- CrawlPage(page):
			// Got some result
		}
	}()

	return out
}

func merge(ctx context.Context, workers []chan CrawledPage) chan CrawledPage {
	// Merged channel with results
	out := make(chan CrawledPage)

	// Synchronization over channels: do not close "out" before all tasks are completed
	var wg sync.WaitGroup

	// Define function which waits the result from worker channel
	// and sends this result to the merged channel.
	// Then it decreases the counter of running tasks via wg.Done().
	output := func(c <-chan CrawledPage) {
		defer wg.Done()
		for result := range c {
			select {
			case <-ctx.Done():
				// Received a signal to abandon further processing
				return
			case out <- result:
				// some message or nothing
			}
		}
	}

	wg.Add(len(workers))
	for _, workerChannel := range workers {
		go output(workerChannel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
