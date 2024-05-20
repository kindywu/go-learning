package main

import (
	"fmt"
	"reflect"
	"sort"
)

func main() {
	slice := make([]int, 10, 20)
	printSlice(slice)
	slice = append(slice, 1, 2, 3)
	printSlice(slice)

	slice2 := make([]int, 15, 20)
	fmt.Println(copy(slice2, slice))
	printSlice(slice2)

	slice3 := slice[0:12:18]
	printSlice(slice3)

	for i := range 10 {
		slice = append(slice, i+10)
	}

	printSlice(slice)

	sort.Slice(slice, func(i, j int) bool { return slice[i] > slice[j] })
	printSlice(slice)
	sort.Ints(slice)
	printSlice(slice)
	sort.SliceStable(slice, func(i, j int) bool { return slice[i] > slice[j] })
	printSlice(slice)
}

func printSlice(slice []int) {
	t := reflect.TypeOf(slice)
	fmt.Printf("Slice: %s Value: %v Length: %d, Capacity: %d \n", t.Name(), slice, len(slice), cap(slice))
}
