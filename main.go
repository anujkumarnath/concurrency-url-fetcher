package main

import (
	"fmt"
	"flag"
	"os/signal"
	"syscall"
	"context"
	"os"
	"time"
)

func main() {
	globalTimeout := flag.Int("globalTimeout", 10, "global timeout for the program")
	concurrency   := flag.Int("concurrency",    5, "no. of concurrent urls to process")
	timeout       := flag.Int("timeout",        5, "timeout in seconds")

	flag.Parse()

	fmt.Println("Flags:")
	fmt.Println(" - globalTimeout  :", *globalTimeout)
	fmt.Println(" - concurrency    :", *concurrency)
	fmt.Println(" - timeout        :", *timeout)
	fmt.Println("----------------------\n")

	// Positional args only
	urlsAsArgs := flag.Args()
	if len(urlsAsArgs) == 0 {
		panic("no urls provided")
	}

	topCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(*globalTimeout) * time.Second,
	)
	defer cancel()

	ctx, stop := signal.NotifyContext(
		topCtx,
		os.Interrupt,
		syscall.SIGTERM,
	);
	defer stop()

	go StartApp(ctx, *concurrency, *timeout, urlsAsArgs)

	<-ctx.Done()
	fmt.Println("Shutting down gracefully...")
}
