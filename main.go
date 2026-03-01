package main

import (
	"fmt"
	"os"
	"flag"
	"sync"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		panic("no urls provided")
	}

	concurrency := flag.Int("concurrency", 5, "no. of concurrent urls to process")
	timeout     := flag.Int("timeout",     5, "timeout in seconds")

	flag.Parse()

	fmt.Println("Flags:")
	fmt.Println("concurrency :", *concurrency)
	fmt.Println("timeout     :", *timeout)
	fmt.Println("---------------\n")

	var wg sync.WaitGroup

	urls    := make(chan string)
	results := make(chan Result)

	workerPoolSize := min(len(argsWithoutProg), *concurrency)
	fmt.Println("workerPoolSize", workerPoolSize)

	// Create n workers, n = allowed concurrency number
	for range workerPoolSize {
		wg.Add(1)
		go worker(&wg, urls, results)
	}

	// Assign jobs to the workers
	// Send jobs from a goroutine so main can receive results concurrently
	go func() {
		for _, arg := range argsWithoutProg {
			urls <- arg
		}

		// Signal end of jobs
		close(urls)
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

	fmt.Println("All URLs processed")
}

func worker(wg *sync.WaitGroup, urls <-chan string, results chan<- Result) {
	defer wg.Done()
	// Listen for jobs as long as there is at least one to process
	for url := range urls {
		result := ProcessUrl(url)
		results <- result
	}
}
