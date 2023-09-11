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

	queue := make([]int, 0, 1000+numRoutines)
	for i := 0; i < 1000; i++ {
		queue = append(queue, 1)
	}
	head := 0

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mutex.Lock()

			item := queue[head]
			head++
			queue = append(queue, item)

			mutex.Unlock()
		}()
	}

	wg.Wait()

	log.Println(head)
	log.Println(len(queue))
}
