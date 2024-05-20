package algorithm

import (
	"testing"
)

func TestMaximalRectangle(t *testing.T) {
	t.Run("4", func(t *testing.T) {
		matrix := [][]byte{
			{0, 1, 1, 0, 0},
			{0, 1, 0, 1, 1},
			{1, 1, 1, 0, 1},
			{1, 1, 0, 1, 0},
		}
		if MaximalRectangle3(matrix) != 4 {
			t.Fatal("fail 4")
		}
	})

	t.Run("5", func(t *testing.T) {
		matrix := [][]byte{
			{1, 0, 1, 0, 0},
			{1, 0, 1, 0, 1},
			{1, 1, 1, 1, 1},
			{1, 0, 0, 1, 0},
		}
		if MaximalRectangle3(matrix) != 5 {
			t.Fatal("fail 5")
		}
	})

	t.Run("6", func(t *testing.T) {
		matrix := [][]byte{
			{1, 0, 1, 0, 0},
			{1, 0, 1, 1, 1},
			{1, 1, 1, 1, 1},
			{1, 0, 0, 1, 0},
		}
		if MaximalRectangle3(matrix) != 6 {
			t.Fatal("fail 6")
		}
	})
}
