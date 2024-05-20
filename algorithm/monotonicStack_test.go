package algorithm

import (
	"reflect"
	"testing"
)

func TestMonotonicStack_Bigger(t *testing.T) {
	source := []int{2, 1, 2, 4, 3}
	want := []int{4, 2, 4, -1, -1}
	got := MonotonicStack(source, Comparator_Bigger)
	if !reflect.DeepEqual(want, got) {
		t.Error("not pass")
	}
}

func TestMonotonicStack_Smaller(t *testing.T) {
	source := []int{2, 1, 2, 4, 3}
	want := []int{1, -1, -1, 3, -1}
	got := MonotonicStack(source, Comparator_Smaller)
	if !reflect.DeepEqual(want, got) {
		t.Error("not pass")
	}
}
