package main

import "fmt"

// empty slice no length/capacity
var declare_nums []int

// nums := make([]int) but of course we have to assign in main

var declare_map map[string]int

func main() {
	//slice literal
	nums := []int{}

	another_map := map[string]int{
		"one": 1,
	}

	aother_literal := map[string]string{
		"water": "wet",
	}

	for k, v := range aother_literal {
		fmt.Println(k + v)
	}

	// slice with make
	myslice := make([]int, 10, 10)
	myslice = append(myslice, 1, 2, 2)

	myMap2 := make(map[int]int)
}
