package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Key string

func main() {

	ctx := context.TODO()
	fmt.Printf("Memory address of instance: %p\n", &ctx)

	ctx = context.WithValue(ctx, Key("key1"), "0001")
	fmt.Printf("Memory address of instance: %p\n", &ctx)

	ctx = context.WithValue(ctx, Key("key2"), "0002")
	fmt.Printf("Memory address of instance: %p\n", &ctx)

	ctx = context.WithValue(ctx, Key("key3"), "0003")
	fmt.Printf("Memory address of instance: %p\n", &ctx)

	ctx = context.WithValue(ctx, Key("key4"), "0004")
	fmt.Printf("Memory address of instance: %p\n", &ctx)

	fmt.Println(ctx.Value(Key("key1")))

	// 初始化随机数生成器
	time1 := time.Duration(rand.Intn(3)+1) * time.Second
	time.Sleep(time.Second)
	time2 := time.Duration(rand.Intn(3)+1) * time.Second

	fmt.Printf("time1=%v,time2=%v\n", time1, time2)

	ctx, cancel := context.WithTimeout(context.Background(), time1)

	go func() {
		defer func() {
			fmt.Println("goroutine exit")
		}()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("timeout")
				return
			default:

				time.Sleep(time.Second)
			}
		}
	}()

	time.Sleep(time2)
	fmt.Println("cancel")
	cancel()
	time.Sleep(2 * time.Second)
}
