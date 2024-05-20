package main

import (
	"fmt"
	"sync"
)

func main() {
	// 封装好的计数器
	var counter Counter
	var count = 1

	var num = 10
	var wg sync.WaitGroup
	wg.Add(num)

	// 启动10个goroutine
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			// 执行10万次累加
			for j := 0; j < 1000; j++ {
				counter.Incr() // 受到锁保护的方法
				count++
			}
		}()
	}
	wg.Wait()
	fmt.Println(counter.Count())
	println(count)
}

// 线程安全的计数器类型
type Counter struct {
	mu    sync.Mutex
	count uint64
}

// 加1的方法，内部使用互斥锁保护
func (c *Counter) Incr() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

// 得到计数器的值，也需要锁保护
func (c *Counter) Count() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// go run -race counter.go
// go build -gcflags "-S" counter.go
