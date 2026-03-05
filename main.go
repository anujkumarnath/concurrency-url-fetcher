package main

import (
	"fmt"
	"flag"
	"os/signal"
	"syscall"
	"context"
	"os"
	"time"
	"runtime"
)

func main() {
	fmt.Println("#goroutines: ", runtime.NumGoroutine())
	globalTimeout := flag.Int("globalTimeout", 10, "global timeout for the program")
	concurrency   := flag.Int("concurrency",    5, "no. of concurrent urls to process")
	timeout       := flag.Int("timeout",        5, "timeout in seconds")
	debug         := flag.Bool("debug",     false, "enable debug logs")

	flag.Parse()

	fmt.Println("Flags:")
	fmt.Println(" - globalTimeout  :", *globalTimeout)
	fmt.Println(" - concurrency    :", *concurrency)
	fmt.Println(" - timeout        :", *timeout)
	fmt.Println(" - debug          :", *debug)
	fmt.Println("--------------------------\n")

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

	done := make(chan struct{})

	go func () {
		StartApp(ctx, *concurrency, *timeout, urlsAsArgs)
		cancel()
		close(done)
	}()

	<-ctx.Done()
	fmt.Println("Shutting down gracefully...")
	<-done

	if *debug {
		time.Sleep(2 * time.Second)
		fmt.Println("#goroutines: ", runtime.NumGoroutine())
		buf := make([]byte, 1<<16)
		n := runtime.Stack(buf, true)
		fmt.Printf("%s", buf[:n])
	}
}
