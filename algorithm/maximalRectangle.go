package algorithm

import "fmt"

func MaximalRectangle(matrix [][]byte) int {
	if len(matrix) == 0 {
		return 0
	}

	m := len(matrix[0])
	n := len(matrix)
	heights := make([]int, m+1)
	maxArea := 0

	// 遍历每一层作为底层
	for i := 0; i < n; i++ {
		stack := make([]int, 0)
		stack = append(stack, -1) // 初始化栈，将-1放入栈底
		for j := 0; j <= m; j++ {
			// 计算j位置的高度
			if j == m || matrix[i][j] == 0 {
				heights[j] = 0
			} else {
				heights[j] = heights[j] + 1
			}

			// 单调栈维护宽度
			for len(stack) > 1 && heights[j] < heights[stack[len(stack)-1]] {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				width := j - stack[len(stack)-1] - 1
				area := width * heights[top]
				if area > maxArea {
					maxArea = area
				}
			}
			stack = append(stack, j)
		}
	}

	return maxArea
}

func MaximalRectangle2(matrix [][]byte) int {
	if len(matrix) == 0 {
		return 0
	}

	m := len(matrix[0])
	n := len(matrix)
	heights := make([]int, m)
	maxArea := 0

	// 遍历每一层作为底层
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			// 计算j位置的高度
			if j == m || matrix[i][j] == 0 {
				heights[j] = 0
			} else {
				heights[j] = heights[j] + 1
			}
		}

		maxArea = max(maxArea, LargestRectangleInHistogram_Scan(heights))

		fmt.Println("height", heights, "maxArea", maxArea)
	}

	return maxArea
}

func MaximalRectangle3(matrix [][]byte) int {
	if len(matrix) == 0 {
		return 0
	}

	m := len(matrix[0])
	n := len(matrix)
	heights := make([]int, m)
	maxArea := 0

	// 遍历每一层作为底层
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			// 计算j位置的高度
			if j == m || matrix[i][j] == 0 {
				heights[j] = 0
			} else {
				heights[j] = heights[j] + 1
			}
		}

		maxArea = max(maxArea, LargestRectangleInHistogram(heights))

		fmt.Println("height", heights, "maxArea", maxArea)
	}

	return maxArea
}
