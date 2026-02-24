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

	for _, arg := range argsWithoutProg {

		reqUrl, err := url.ParseRequestURI(arg)
		if err != nil {
			fmt.Printf("%-8s : bad-url\n",     "URL")
			fmt.Printf("%-8s : invalid URL\n", "Error")
			fmt.Println("--\n")
			continue
		}

		urlString := reqUrl.String()
		fmt.Printf("%-8s : %s\n", "URL", urlString)

		startTime := time.Now()
		resp, err := http.Get(urlString)

		if err != nil {
			var dnsErr *net.DNSError
			if errors.As(err, &dnsErr) {
				fmt.Printf("%-8s : DNS lookup failed\n", "Error")
			} else {
				fmt.Printf("%-8s : %s\n", "Error", err.Error())
			}
			fmt.Println("--\n")
			continue
		}

		endTime  := time.Now()
		duration := endTime.Sub(startTime).Milliseconds()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%-8s : %s\n",    "Error",    err.Error())
			fmt.Printf("%-8s : %s ms\n", "Duration", duration)
			fmt.Println("--\n")
			continue
		}

		fmt.Printf("%-8s : %s\n",       "Status",   resp.Status)
		fmt.Printf("%-8s : %d bytes\n", "Size",     len(body))
		fmt.Printf("%-8s : %d ms\n",    "Duration", duration)
		fmt.Println("--\n")

		resp.Body.Close()
	}
}
