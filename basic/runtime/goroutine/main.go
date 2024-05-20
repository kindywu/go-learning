package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"
)

func rung(n int) {
	v := int32(0)
	var wg sync.WaitGroup
	wg.Add(n + 1)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Wait()
			defer wg.Done()
			atomic.AddInt32(&v, 1)
		}()
	}

	for {
		if atomic.LoadInt32(&v) == int32(n) {
			fmt.Println("Goroutine Number: ", runtime.NumGoroutine())
			break
		}
	}
	wg.Done()
	wg.Wait()
}

func main() {
	f, err := os.Create("pprof.out")
	if err != nil {
		panic(err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	fmt.Println("Goroutine Number: ", runtime.NumGoroutine())

	fmt.Println("Rung 1000")
	rung(1000)
	time.Sleep(1 * time.Second)
	fmt.Println("After Rung 1000 Goroutine Number: ", runtime.NumGoroutine())

	fmt.Println("Rung 100")
	rung(100)
	time.Sleep(1 * time.Second)
	fmt.Println("After Rung 100 Goroutine Number: ", runtime.NumGoroutine())

	rung(10000000)
	time.Sleep(1 * time.Second)
	fmt.Println("After Rung 10000000 Goroutine Number: ", runtime.NumGoroutine())
}
