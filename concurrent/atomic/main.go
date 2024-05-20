package main

import (
	"fmt"
	"main/concurrent/atomic/queue"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	NodeName string
	Addr     string
	Count    int32
}

func loadNewConfig() Config {
	return Config{
		NodeName: "北京",
		Addr:     "10.77.95.27",
		Count:    rand.Int31(),
	}
}
func main() {
	q := queue.NewLKQueue[int]()
	q.Enqueue(100)
	q.Enqueue(101)
	q.Enqueue(102)
	fmt.Println(*q.Dequeue())
	fmt.Println(*q.Dequeue())
	fmt.Println(*q.Dequeue())

	// configChanged()
}

func configChanged() {
	// 设置新的config
	// 通知等待着配置已变更
	// 等待变更信号
	// 读取新的配置

	var config atomic.Value
	config.Store(loadNewConfig())
	var cond = sync.NewCond(&sync.Mutex{})

	go func() {
		for {
			time.Sleep(time.Duration(5+rand.Int63n(5)) * time.Second)
			config.Store(loadNewConfig())
			cond.Broadcast()
		}
	}()

	go func() {
		cond.L.Lock()
		defer cond.L.Unlock()

		for {
			cond.Wait()
			c := config.Load().(Config)
			fmt.Printf("new config: %+v\n", c)
		}
	}()

	select {}
}
