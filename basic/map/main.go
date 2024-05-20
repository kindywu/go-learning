package main

import (
	"fmt"
	"sync"
)

func main() {
	m := make(map[string]int)
	m["a"] = 1
	fmt.Printf("a=%d; b=%d\n", m["a"], m["b"])
	av, aexisted := m["a"]
	bv, bexisted := m["b"]
	fmt.Printf("a=%d, existed: %t; b=%d, existed: %t\n", av, aexisted, bv, bexisted)

	m["b"] = 2
	m["c"] = 3
	m["d"] = 4

	delete(m, "c")
	delete(m, "e")

	for k, v := range m {
		fmt.Printf("%s=%d\n", k, v)
	}

	var m2 map[int]int
	fmt.Println(m2[100]) // will not panic when we visit a nil map
	// m2[100] = 100

	// fatal error: concurrent map read and map write
	var wg sync.WaitGroup
	wg.Add(2)
	var m3 = make(map[int]int, 10) // 初始化一个map
	go func() {
		defer wg.Done()
		for {
			m3[1] = 1 //设置key
		}
	}()

	go func() {
		defer wg.Done()
		for {
			_ = m3[2] //访问这个map
		}
	}()
	wg.Wait()
}
