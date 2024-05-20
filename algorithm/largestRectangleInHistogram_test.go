package algorithm

import (
	"testing"
)

var largestRectangleInHistogramCases = map[string]struct {
	input []int
	want  int
}{
	"18": {
		input: []int{2, 1, 5, 4, 9, 10},
		want:  18,
	},
	"10": {
		input: []int{2, 1, 5, 6, 2, 3},
		want:  10,
	},
	"25": {
		input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		want:  25,
	},
}

func TestLargestRectangleInHistogram(t *testing.T) {
	for name, case0 := range largestRectangleInHistogramCases {
		t.Run(name, func(t *testing.T) {
			if got := LargestRectangleInHistogram(case0.input); got != case0.want {
				t.Fatalf("case name:%s want:%d got:%d\n", name, case0.want, got)
			}
		})
	}
}

func TestLargestRectangleInHistogram_Scan(t *testing.T) {
	for name, case0 := range largestRectangleInHistogramCases {
		t.Run(name, func(t *testing.T) {
			if got := LargestRectangleInHistogram_Scan(case0.input); got != case0.want {
				t.Fatalf("case name:%s want:%d got:%d\n", name, case0.want, got)
			}
		})
	}
}

const Size = 100000
const Want = 2500000000

// go test -timeout 30s -run ^TestForMaxArea$ main/algorithm -v
func TestForMaxArea(t *testing.T) {
	heights := makeBigSlice(Size)
	t.Log("max_area", LargestRectangleInHistogram(heights))
}

func BenchmarkLargestRectangleInHistogram(b *testing.B) {
	heights := makeBigSlice(Size)
	if LargestRectangleInHistogram(heights) != Want {
		b.FailNow()
	}
}

func BenchmarkLargestRectangleInHistogram_Scan(b *testing.B) {
	heights := makeBigSlice(Size)
	if LargestRectangleInHistogram_Scan(heights) != Want {
		b.FailNow()
	}
}
