package utils

func Max(nums []float64) float64 {
	if len(nums) == 0 {
		panic("empty slice")
	}

	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}

	return max
}

func Min(nums []float64) float64 {
	if len(nums) == 0 {
		panic("empty slice")
	}

	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}

	return min
}
