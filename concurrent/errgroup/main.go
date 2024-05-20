package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(1)

	// 启动第一个子任务,它执行成功
	g.Go(func() error {
		time.Sleep(5 * time.Second)
		fmt.Println("exec #1")
		return nil
	})
	// 启动第二个子任务，它执行失败
	g.Go(func() error {
		time.Sleep(1 * time.Second)
		fmt.Println("exec #2")
		return errors.New("failed to exec #2")
	})

	// 启动第三个子任务，它执行成功
	g.Go(func() error {
		time.Sleep(3 * time.Second)
		fmt.Println("exec #3")
		return nil
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			fmt.Println("#4 canceled")
			return ctx.Err()
		case <-time.After(6 * time.Second):
			fmt.Println("exec #4")
			return nil
		}
	})

	g.TryGo(func() error {
		time.Sleep(2 * time.Second)
		fmt.Println("exec #5")
		return nil
	})

	// 等待四个任务都完成
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully exec all")
	} else {
		fmt.Println("failed:", err)
	}
}
