package main

import (
	"fmt"
	"log/slog"
	cpool "main/concurrent/goroutine_pool/cpool"
	"time"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	// slog.SetLogLoggerLevel(slog.LevelDebug)

	// 创建一个协程池，设置协程数为3
	pool := cpool.NewPool(3)

	// 提交任务到协程池
	tasks := make([]*cpool.Task, 10)
	for i := 0; i < 10; i++ {
		taskName := fmt.Sprintf("task_%d", i)
		t := cpool.NewTask(taskName, func() (interface{}, error) {
			time.Sleep(1 * time.Second)
			if i > 0 && i%3 == 0 {
				panic(fmt.Sprintf("panic for 3 from %s", taskName))
			}
			if i > 0 && i%4 == 0 {
				return nil, fmt.Errorf("error for 4 from %s", taskName)
			}
			return i, nil
		})
		tasks[i] = t
		pool.Submit(t)
	}

	for _, t := range tasks {
		if result, ok := <-t.Done; ok {
			fmt.Printf("t result:%v error:%v\n", result.Result, result.Err)
		} else {
			fmt.Println("nothing from", t.Name)
		}
	}

	// 等待所有任务完成
	pool.Stop()
	time.Sleep(3 * time.Second)
}
