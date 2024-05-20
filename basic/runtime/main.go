package main

import (
	"os/exec"
	"runtime/debug"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	debug.SetMaxThreads(10)
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := exec.Command("bash", "-c", "sleep 3").Output()
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
}
