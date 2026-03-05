package main

import (
	"sync"
	"context"
)

func StartApp(ctx context.Context, concurrency int, timeout int, urlsAsArgs []string) {
	urlProcessor := NewUrlProcessor(timeout)

	// To avoid leaking goroutine from open http client
	defer urlProcessor.Clean()

	var wg sync.WaitGroup

	urls    := make(chan string)
	results := make(chan Result)

	// Using this instead of map[string]bool,
	// because bool takes 1 byte but empty struct takes 0 bytes
	seen := make(map[string]struct{})
	var uniqueUrls []string

	for _, arg := range urlsAsArgs {
		if _, present := seen[arg]; present {
			continue
		}
		uniqueUrls = append(uniqueUrls, arg)
		// struct{} defines the type,
		// and the last {} assigns an empty value to it
		seen[arg] = struct{}{}
	}

	workerPoolSize := min(len(uniqueUrls), concurrency)

	// Create n workers, n = allowed concurrency number
	for range workerPoolSize {
		wg.Add(1)
		go worker(ctx, &wg, urls, results, urlProcessor)
	}

	// Assign jobs to the workers
	// Send jobs from a goroutine so main can receive results concurrently
	go func() {
		// Signal end of jobs
		defer close(urls)

		for _, arg := range uniqueUrls {
			select {
			case <-ctx.Done():
				return
			default:
				urls <- arg
			}
		}

	}()

	// Close results once all workers are done, unblocking the range loop in main
	go func() {
		wg.Wait()
		close(results)
	}()

	// Read all results until results channel is closed
	for result := range results {
		PrintResult(result)
	}

}

func worker(
	ctx     context.Context,
	wg      *sync.WaitGroup,
	urls    <-chan string,
	results chan<- Result,
	up      *UrlProcessor,
) {
	defer wg.Done()
	// Listen for jobs as long as there is at least one to process
	for url := range urls {
		result := up.ProcessUrl(ctx, url)
		results <- result
	}
}
