package main

import (
	"sync"
	"sync/atomic"
	"time"
)

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 0 {
		return o.doSlow(f)
	}
	return nil
}

func (o *Once) doSlow(f func() error) error {
	o.m.Lock()
	defer o.m.Unlock()
	var err error
	if o.done == 0 {
		err = f()
		// 只有没有 error 的时候，才修改 done 的值
		if err == nil {
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}

func main() {
	var wg sync.WaitGroup
	wg.Add(10)
	var once Once
	for range 10 {
		go func() {
			defer wg.Done()
			err := once.Do(func() error {
				time.Sleep(3 * time.Second)
				println("once done")
				return nil
			})
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
}
