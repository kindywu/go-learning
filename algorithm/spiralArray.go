package algorithm

import "fmt"

//go:test -timeout 30s -run ^TestSpiralArray_JustPrint$ main/algorithm -v

// SpiralArray creates a 2D array filled with numbers in a spiral order.
func SpiralArray(n int) [][]int {
	// Initialize the n x n array with zeros.
	arr := make([][]int, n)
	for i := range arr {
		arr[i] = make([]int, n)
	}

	// Define the starting and ending indices for filling the array.
	left, right, top, bottom := 0, n-1, 0, n-1
	// Start filling the array with a value 'k'.
	k := 1

	// Continue filling the array until all elements are set.
	for k <= n*n {
		// Fill from left to right.
		for x := left; x <= right; x++ {
			arr[top][x] = k
			k++
		}
		top++ // Move the top boundary down.

		// Fill from top to bottom.
		for y := top; y <= bottom; y++ {
			arr[y][right] = k
			k++
		}
		right-- // Move the right boundary left.

		// Fill from right to left.
		for x := right; x >= left; x-- {
			arr[bottom][x] = k
			k++
		}
		bottom-- // Move the bottom boundary up.

		// Fill from bottom to top.
		for y := bottom; y >= top; y-- {
			arr[y][left] = k
			k++
		}
		left++ // Move the left boundary right.
	}

	return arr
}

// PrintArray prints the elements of a 2D array in a formatted way.
func PrintArray(matrix [][]int) {
	for _, row := range matrix {
		for _, val := range row {
			// Print a space after numbers less than 10 for better formatting.
			if val < 10 {
				fmt.Printf("%d  ", val)
			} else {
				fmt.Printf("%d ", val)
			}
		}
		fmt.Println() // Print a newline after each row.
	}
}
