package main

import (
	"fmt"
	"main/concurrent/map/rwmap"
	"sync"
)

func main() {
	m := rwmap.NewRWMap[int, string](10)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := range 10 {
			m.Set(i, fmt.Sprintf("value %d", i))
		}
	}()
	go func() {
		defer wg.Done()
		for range 10 {
			for i := range 10 {
				fmt.Println(m.Get(i))
			}
		}

	}()
	wg.Wait()
}
