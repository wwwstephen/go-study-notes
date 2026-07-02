package arraysandslices

func Sum(numbers []int) int {
	r := 0
	for _, number := range numbers {
		r += number
	}
	return r
}

func SumAll(numbers1 ...[]int) []int {
	nums := []int{}

	for _, list := range numbers1 {
		value := 0
		for _, number := range list {
			value += number
		}
		nums = append(nums, value)
	}

	return nums
}

func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		length := len(numbers)
		if length == 0 {
			return append(sums, 0)
		}
		tail := numbers[length-1:]
		sums = append(sums, Sum(tail))
	}

	return sums
}
