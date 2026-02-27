package main

import (
	"fmt"
	"os"
	"flag"
	"errors"
	"net"
	"net/http"
	"net/url"
	"io"
	"time"
	"sync"
)

type Result struct {
	URL      string
	Status   string
	Size     int
	Duration int64
	Error    string
}

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

	/* Create n workers, n = allowed concurrency number */
	for range *concurrency {
		wg.Add(1)
		go worker(&wg, urls, results)
	}

	/* Assign jobs to the workers */
	for _, arg := range argsWithoutProg {
		reqUrl, err := url.ParseRequestURI(arg)
		if err != nil {
			fmt.Printf("%-8s : bad-url\n",     "URL")
			fmt.Printf("%-8s : invalid URL\n", "Error")
			fmt.Println("--\n")
			continue
		}

		urlString := reqUrl.String()
		urls <- urlString
	}

	/* Signal end of jobs */
	close(urls)

	go func() {
		/* Wait for all jobs to complete (all workes called wg.Done())*/
		wg.Wait()
		/* This allows main to move ahead and exit */
		close(results)
	}()

	/* Read all results till results channel is closed (no more result to process) */
	/* This also unblocks results channel */
	 for result := range results {
		printResult(result)
	}

	fmt.Println("All URLs processed")
}

func worker(wg *sync.WaitGroup, urls <-chan string, results chan<- Result) {
	defer wg.Done()
	/* Listen for jobs as long as there is at least one to process */
	for url := range urls {
		result := processUrl(url)
		/* Blocks on results channel */
		results <- result
	}
}

func processUrl(urlString string) Result {
	var result Result

	result.URL = urlString

	startTime := time.Now()
	resp, err := http.Get(urlString)

	if err != nil {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) {
			result.Error = "DNS lookup failed"
		} else {
			result.Error = err.Error()
		}
		return result
	}

	defer resp.Body.Close()

	endTime  := time.Now()
	duration := endTime.Sub(startTime).Milliseconds()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = err.Error()
		result.Duration = duration
		return result
	}

	result.Status   = resp.Status
	result.Size     = len(body)
	result.Duration = duration

	return result
}

func printResult(result Result) {
	fmt.Printf("%-8s : %s\n", "URL", result.URL)

	if result.Error == "" {
		fmt.Printf("%-8s : %s\n",       "Status",   result.Status)
		fmt.Printf("%-8s : %d bytes\n", "Size",     result.Size)
		fmt.Printf("%-8s : %d ms\n",    "Duration", result.Duration)
	} else {
		fmt.Printf("%-8s : %s\n", "Error", result.Error)
	}

	fmt.Println("--\n")
}
