package main

import (
	"runtime"
	"sync"
	"sync/atomic"
)

const MaxBackoff = 16

type SpinLock struct {
	flag int32
}

func (s *SpinLock) Lock() {
	backoff := 1
	for !atomic.CompareAndSwapInt32(&s.flag, 0, 1) {
		// Leverage the exponential backoff algorithm, see https://en.wikipedia.org/wiki/Exponential_backoff.
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		if backoff < MaxBackoff {
			backoff <<= 1
		}
	}
}

func (s *SpinLock) Unlock() {
	atomic.StoreInt32(&s.flag, 0)
}

func main() {
	lock := SpinLock{}

	wg := sync.WaitGroup{}
	count := 0
	for range 1000 {
		wg.Add(1)
		go func() {
			lock.Lock()
			for range 1000 {
				count++
			}
			defer lock.Unlock()
			defer wg.Done()
		}()
	}

	wg.Wait()
	println("count =", count)
}
