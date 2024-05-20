package main

import (
	"fmt"

	"github.com/elliotchance/orderedmap/v2"
)

func main() {
	m := orderedmap.NewOrderedMap[string, any]()

	m.Set("foo", "bar")
	m.Set("qux", 1.23)
	m.Set("123", true)

	m.Delete("qux")

	for _, key := range m.Keys() {
		value, _ := m.Get(key)
		fmt.Println(key, value)
	}
}
