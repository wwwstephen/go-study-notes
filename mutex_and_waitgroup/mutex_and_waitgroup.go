package main

import (
	"sync"
)

func doWork() {
	// pretend work
}

func main() {
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		doWork()
	}()

	go func() {
		defer wg.Done()
		doWork()
	}()

	go func() {
		defer wg.Done()
		doWork()
	}()

	wg.Wait()
}
