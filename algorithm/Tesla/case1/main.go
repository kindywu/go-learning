package main

import (
	"fmt"
)

// 返回计算数组中和为0的连续子数组中长度最长的数量
func MaxSumZeroSliceLength(A []int) int {
	if len(A) == 0 {
		return -1
	}

	if len(A) >= 100000 {
		return -1
	}

	sum := make([]int, len(A)+1) // 累积和历史列表
	maxLen := 0                  // 最大长度
	currSum := 0                 // 当前累积和

	sum[0] = 0
	for i, num := range A {
		currSum += num     // 更新当前累积和
		sum[i+1] = currSum // 新增累积和历史列表

		// 反向查相同的累积和
		for j := i; j >= 0; j-- {
			// 发现相等的历史累积和
			if sum[j] == currSum {
				maxLen = max(maxLen, i-j+1)
				break
			}
		}
	}

	fmt.Println(sum)
	return maxLen
}

func main() {
	fmt.Println(MaxSumZeroSliceLength([]int{-3, -4, 1, 1, 1, 1, 1, 2, 2}))
	// [0 -3 -7 -6 -5 -4 -3 -2 0 2] => 0~0=8
	fmt.Println(MaxSumZeroSliceLength([]int{3, -3, -4, 1, 1, 1, 1, 1, 2, 2}))
	// [0 3 0 -4 -3 -2 -1 0 1 3 5] => 3~3=8
	fmt.Println(MaxSumZeroSliceLength([]int{2, -2, 3, 0, 4, -7}))
	// [0 2 0 3 3 7 0] => 0~0=4
}
