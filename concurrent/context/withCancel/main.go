package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context) {
	fmt.Println("Worker started.")

	// 执行一些工作
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker exiting due to cancel signal. Err is", ctx.Err().Error())
			return
		default:
			fmt.Printf("Worker is doing work... (%v)\n", time.Now().Format(time.Stamp))
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// 启动一个工作协程
	go worker(ctx)

	// 主协程做其他事情或者等待
	time.Sleep(5 * time.Second)

	// 5秒后，取消工作协程的上下文
	fmt.Println("Main goroutine will cancel the worker.")
	cancel()

	// 等待工作协程退出
	time.Sleep(1 * time.Second)
	fmt.Println("Main goroutine finished.")
}
