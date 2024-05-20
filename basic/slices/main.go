package main

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strconv"
)

func main() {
	newSlice := slices.Grow[[]int](nil, 1000)
	fmt.Printf("cap: %d, len: %d\n", cap(newSlice), len(newSlice))
	fmt.Println(newSlice)

	println("******************")
	useBinarySearch()
	println("******************")
	useBinarySearchFunc()
	println("******************")
	useClip()
	println("******************")
	useClone()
	println("******************")
	useContainsFunc()
	println("******************")
	fmt.Println(useEqualFunc())
	println("******************")
}

func useEqualFunc() bool {
	numbers := []int{0, 42, 8}
	strings := []string{"000", "42", "0o10"}
	equal := slices.EqualFunc(numbers, strings, func(n int, s string) bool {
		sn, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return false
		}
		return n == int(sn)
	})
	return equal
}

func useContainsFunc() {
	numbers := []int{0, 42, -10, 8}
	hasNegative := slices.ContainsFunc(numbers, func(n int) bool {
		return n < 0
	})
	fmt.Println("Has a negative:", hasNegative)
	hasOdd := slices.ContainsFunc(numbers, func(n int) bool {
		return n%2 != 0
	})
	fmt.Println("Has an odd number:", hasOdd)
}

func useClone() {
	names := []string{"路多辛的博客", "路多辛的所思所想"}
	namesCopy := slices.Clone(names)
	fmt.Println(namesCopy)
}

func useClip() {
	names := make([]string, 2, 5)
	fmt.Printf("长度：%d,容量：%d\n", len(names), cap(names))
	names = slices.Clip(names)
	fmt.Printf("长度：%d,容量：%d\n", len(names), cap(names))
}

func useBinarySearch() {
	names := []string{"Boss", "Jack", "Black", "Alice", "Bob", "Vera", "Blue"}
	sort.Strings(names)
	n, found := slices.BinarySearch(names, "Vera")
	fmt.Println("Vera:", n, found)

	n, found = slices.BinarySearch(names, "Bill")
	fmt.Println("Bill:", n, found)
}

func useBinarySearchFunc() {
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{"Alice", 55},
		{"Bob", 24},
		{"Gopher", 13},
		{"Blue", 24},
	}
	n, found := slices.BinarySearchFunc(people, Person{"Blue", 0}, func(a, b Person) int {
		return cmp.Compare(a.Name, b.Name)
	})
	fmt.Println("Bob:", n, found) // Bob: 1 true
}
