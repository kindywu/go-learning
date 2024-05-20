package algorithm

import (
	"sync"
	"testing"
)

var smallestPositiveIntegerCases = map[string]struct {
	input []int
	want  int
}{
	"number": {
		input: []int{1, 3, 6, 4, 1, 2},
		want:  5,
	},
	"ordered number": {
		input: []int{1, 2, 3, 4},
		want:  5,
	},
	"negative number": {
		input: []int{-1, -3},
		want:  1,
	},
}

func TestSmallestPositiveInteger(t *testing.T) {
	for name, case0 := range smallestPositiveIntegerCases {
		t.Run(name, func(t *testing.T) {
			if got := SmallestPositiveInteger(case0.input); got != case0.want {
				t.Fatalf("case name:%s want:%d got:%d\n", name, case0.want, got)
			}
		})
	}
}

func BenchmarkSmallestPositiveInteger(b *testing.B) {
	size := 500000
	A := makeBigSlice(size)
	if SmallestPositiveInteger(A) != size {
		b.FailNow()
	}
}

func BenchmarkSmallestPositiveInteger2(b *testing.B) {
	size := 500000
	A := makeBigSlice(size)
	if SmallestPositiveInteger2(A) != size {
		b.FailNow()
	}
}

func makeBigSlice(size int) []int {
	var slice []int
	var once sync.Once
	once.Do(func() {
		slice = make([]int, size)
		for i := 0; i < size; i++ {
			slice[i] = i
		}
	})
	return slice
}
