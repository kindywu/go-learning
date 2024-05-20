package main

import (
	"fmt"
)

// Solution 函数用于找出最大网络等级
func Solution(A []int, B []int, N int) int {
	// 初始化一个映射，用于存储每个城市的连接道路数量
	roadCount := make(map[int]int)
	// 遍历所有道路，统计每个城市的连接道路数量
	for _, cityA := range A {
		roadCount[cityA]++
	}
	for _, cityB := range B {
		roadCount[cityB]++
	}

	connected := make(map[string]bool, len(A)*2)
	for i := 0; i < len(A); i++ {
		connected[fmt.Sprintf("%d-%d", A[i], B[i])] = true
		connected[fmt.Sprintf("%d-%d", B[i], A[i])] = true
	}

	// 初始化最大网络等级为0
	maxNetworkRank := 0

	// 遍历所有城市，找出连接道路最多的两个城市
	for cityA, rankA := range roadCount {
		for cityB, rankB := range roadCount {
			// 避免自连接的情况，确保城市之间有连接
			if cityA != cityB && (connected[fmt.Sprintf("%d-%d", cityA, cityB)]) {
				// 计算当前两个城市的网络等级，并更新最大网络等级
				maxNetworkRank = max(maxNetworkRank, rankA+rankB)
			}
		}
	}
	return maxNetworkRank - 1 // 去掉两个城市之间的冗余线
}

func main() {
	// 测试用例1
	A1 := []int{1, 2, 3, 3}
	B1 := []int{2, 3, 1, 4}
	N1 := 4
	fmt.Println(Solution(A1, B1, N1)) // 应该输出 4

	// 测试用例2
	A2 := []int{1, 2, 4, 5}
	B2 := []int{2, 3, 5, 6}
	N2 := 6
	fmt.Println(Solution(A2, B2, N2)) // 应该输出 2

	// 测试用例2
	A3 := []int{1, 2, 4, 5, 2}
	B3 := []int{2, 3, 5, 6, 5}
	N3 := 6
	fmt.Println(Solution(A3, B3, N3)) // 应该输出 2
}
