package main

import "fmt"

// Status is a custom type based on int.
// We use this to create "enum-like" values.
type Status int

// iota automatically generates incrementing values starting from 0
const (
	Pending  Status = iota // 0
	Active                 // 1
	Inactive               // 2
	Deleted                // 3
)

func main() {
	// Declare a variable of type Status
	// Assign it one of the defined constants
	var s Status = Active

	// Printing s will show its underlying numeric value (1)
	fmt.Println(s)
}
