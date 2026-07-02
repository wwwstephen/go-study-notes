package main

// Software often kicks off long-running, resource-intensive processes (often in goroutines). If the action that caused this gets cancelled or fails for some reason you need to stop these processes in a consistent way through your application.

// If you don't manage this your snappy Go application that you're so proud of could start having difficult to debug performance problems.

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("stopping work")
				return
			default:
				fmt.Println("working...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	time.Sleep(2 * time.Second)
	cancel() // <-- THIS is the missing piece

	time.Sleep(500 * time.Millisecond)
}
