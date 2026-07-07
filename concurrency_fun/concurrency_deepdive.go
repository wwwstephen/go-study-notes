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
