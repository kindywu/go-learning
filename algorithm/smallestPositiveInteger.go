package algorithm

import "slices"

const maxCount = 1000000

func SmallestPositiveInteger(A []int) int {
	hash := make(map[int]bool, len(A))
	max := A[0]
	for _, value := range A {
		if value > 0 {
			hash[value] = true
			if value > max {
				max = value
			}
		}
	}
	if len(hash) == 0 {
		return 1
	}
	// N is an integer within the range [1..100,000]
	for i := 1; i <= max+1; i++ {
		if _, ok := hash[i]; ok {
			continue
		} else {
			return i
		}
	}
	return -1
}

func SmallestPositiveInteger2(A []int) int {
	A = slices.DeleteFunc(A, func(e int) bool {
		return e < 0
	})
	if len(A) == 0 {
		return 1
	}
	// N is an integer within the range [1..100,000]
	for i := 1; i <= maxCount; i++ {
		if slices.Contains(A, i) {
			continue
		} else {
			return i
		}
	}
	return -1
}
