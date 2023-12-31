package main

import (
	"log"
	"os"
	"strconv"
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

	queue := make([]int, 0, numRoutines)
	for i := 0; i < numRoutines; i++ {
		queue = append(queue, 1)
	}
	head := 0

	for i := 0; i < numRoutines; i++ {
		item := queue[head]
		queue[head] = item + 1
		head++
	}
}
