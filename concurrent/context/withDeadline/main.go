package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 设置任务的截止时间为当前时间后的10秒
	deadline := time.Now().Add(10 * time.Second)

	// 使用 WithDeadline 创建一个带有截止时间的上下文
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel() // 清理上下文
	// 启动任务执行的协程
	go func(ctx context.Context) {
		fmt.Println("Task started.")
		select {
		case <-time.After(15 * time.Second): // 模拟任务需要15秒完成
			fmt.Println("Task completed successfully.")
		case <-ctx.Done():
			fmt.Printf("Task canceled before completion: %v\n", ctx.Err())
		}
	}(ctx)

	// 主协程等待任务完成或者截止时间到达
	<-ctx.Done()

	fmt.Println("Main goroutine finished.")

	// Respect OS stop signals.
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a termination signal.
	<-c
}
