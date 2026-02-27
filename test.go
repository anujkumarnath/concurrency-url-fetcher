package main

import (
	"fmt"
	"sync"
	"time"
) 

func work(wg *sync.WaitGroup, id int, jobs <-chan int) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(1 * time.Second)
	}
}

const WORKERS int = 5
const JOBS int = 10

func main() {
	jobs := make(chan int)
	var wg sync.WaitGroup
	startTime := time.Now()

	for i := range WORKERS {
		wg.Add(1)
		go work(&wg, i+1, jobs)
	}

	for i := range JOBS {
		jobs <- i+1
	}

	close(jobs)

	wg.Wait()
	fmt.Println("All jobs done!")
	timeTaken := time.Now().Sub(startTime).Seconds()
	fmt.Printf("Total time: %.2fs\n", timeTaken)
}
