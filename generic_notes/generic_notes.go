package main

// All non-test .go files in the same directory must declare the same package name

import "fmt"

type Number interface {
	int | float64
}

func Sum[T Number](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

func main() {
	ints := []int{1, 2, 3}
	floats := []float64{1.5, 2.5, 3.0}

	fmt.Println(Sum(ints))   // 6
	fmt.Println(Sum(floats)) // 7.0
}
