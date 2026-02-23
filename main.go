package main

import (
	"fmt"
	"os"
	"flag"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		panic("no urls provided")
	}

	fmt.Println("URLs:")
	for _, arg := range argsWithoutProg {
		fmt.Println(arg)
	}
	fmt.Println("-----------------------")

	concurrency := flag.Int("concurrency", 5, "no. of concurrent urls to process")
	timeout     := flag.Int("timeout",     5, "timeout in seconds")

	flag.Parse()

	fmt.Println("Flags:")
	fmt.Println("concurrency:", *concurrency)
	fmt.Println("timeout:", *timeout)
	fmt.Println("-----------------------")
}
