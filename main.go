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

	for _, arg := range argsWithoutProg {
		reqUrl, err := url.ParseRequestURI(arg)
		if err != nil {
			fmt.Printf("%-8s : bad-url\n",     "URL")
			fmt.Printf("%-8s : invalid URL\n", "Error")
			fmt.Println("--\n")
			continue
		}

		urlString := reqUrl.String()

		wg.Add(1)

		go func() {
			defer wg.Done()
			result := processUrl(urlString)

			fmt.Printf("%-8s : %s\n", "URL", result.URL)

			if result.Error == "" {
				fmt.Printf("%-8s : %s\n",       "Status",   result.Status)
				fmt.Printf("%-8s : %d bytes\n", "Size",     result.Size)
				fmt.Printf("%-8s : %d ms\n",    "Duration", result.Duration)
			} else {
				fmt.Printf("%-8s : %s\n", "Error", result.Error)
			}

			fmt.Println("--\n")
		}()

	}

	wg.Wait()
}
