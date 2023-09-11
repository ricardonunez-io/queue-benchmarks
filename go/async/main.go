package main

import (
	"log"
	"os"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Please provide a number of goroutines to run.")
	}
	numRoutines, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln("Failed to parse number of goroutines to run provided in argument")
	}
	if numRoutines < 100 {
		numRoutines = 100
	}

	queue := make(chan int, numRoutines)
	for i := 0; i < numRoutines; i++ {
		queue <- 1
	}

	var wg sync.WaitGroup

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			item, ok := <-queue
			if ok {
				queue <- item + 1
			}
		}()
	}

	wg.Wait()
	close(queue)
}
