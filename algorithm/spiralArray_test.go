package algorithm

import (
	"reflect"
	"testing"
)

func TestSpiralArray(t *testing.T) {
	n := 3
	arr := SpiralArray(n)

	except := [][]int{
		{1, 2, 3}, {8, 9, 4}, {7, 6, 5},
	}
	PrintArray(arr)
	PrintArray(except)
	if !reflect.DeepEqual(arr, except) {

		t.Fatalf("error")
	}
	t.Log("pass")
}

// go test -timeout 30s -run ^TestSpiralArray_JustPrint$ main/algorithm -v
func TestSpiralArray_JustPrint(t *testing.T) {
	n := 9
	arr := SpiralArray(n)
	PrintArray(arr)
	t.Log("pass")
}
