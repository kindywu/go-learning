package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vardius/gollback"
)

func main() {
	rs, errs := gollback.All( // 调用All方法
		context.Background(),
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(3 * time.Second)
			return 1, nil // 第一个任务没有错误，返回1
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("failed") // 第二个任务返回一个错误
		},
		func(ctx context.Context) (interface{}, error) {
			return 3, nil // 第三个任务没有错误，返回3
		},
	)

	fmt.Println(rs)   // 输出子任务的结果
	fmt.Println(errs) // 输出子任务的错误信息

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 尝试5次，或者超时返回
	count := 0
	res, err := gollback.Retry(ctx, 5, func(ctx context.Context) (interface{}, error) {
		count++
		if count < 3 {
			return nil, errors.New("failed")
		}
		return 1, nil
	})
	fmt.Println(res) // 输出结果
	fmt.Println(err) // 输出错误信息
}
