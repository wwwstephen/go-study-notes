package concurrency

package main

import "fmt"

func main() {
    // Create a channel that can only send and receive bool values.
    ch := make(chan bool)

    // Start a new goroutine.
    go func() {
        // Send the boolean value true into the channel.
        // This will block until another goroutine receives the value
        // because this is an unbuffered channel.
        ch <- true
    }()

    // Receive a value from the channel.
    // The goroutine above is paused until this receive happens.
    // The value that comes out is still a bool.
    value := <-ch

    // Print the type of the received value.
    // Output: bool
    fmt.Printf("%T\n", value)
}

// chan bool is a built-in Go type. It is not the value being sent; it is the communication mechanism that transports values.

// In this example:

// ch <- true
// ch is the channel (the pipe/mailbox).
// bool is the type of value that the channel can carry.
// true is the actual value being sent.


// Primitive	Purpose
// Goroutine	A concurrency primitive for executing code concurrently.
// Channel	A communication and synchronization primitive for passing data between goroutines.
// WaitGroup	A synchronization primitive for waiting until a set of goroutines has completed.

//select keyword. It's like switch for channels.
select {
case value := <-ch1:
    fmt.Println("Received:", value)

case value := <-ch2:
    fmt.Println("Received:", value)
}

//Mutex. Lock count if many goroutines are going to increase it 
var (
    mu    sync.Mutex
    count int
)

go func() {
    mu.Lock()
    count++
    mu.Unlock()
}()


// Mutex

// Use when:

// multiple goroutines share a data structure
// you need to protect maps, slices, or structs
// Atomic

// Use when:

// you're updating one variable
// counters
// flags
// statistics
// pointers

package main

import (
    "fmt"
    "time"
)

func main() {
    channels := make([]chan bool, 3)

    for i := range channels {
        channels[i] = make(chan bool)

        go func(id int, ch chan bool) {
            time.Sleep(time.Duration(id) * time.Second)

            fmt.Println("Worker", id, "finished")

            ch <- true
        }(i, channels[i])
    }

    for _, ch := range channels {
        <-ch
    }

    fmt.Println("Everyone is done!")
}