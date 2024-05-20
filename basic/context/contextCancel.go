package basic

import (
	"context"
	"fmt"
)

func GenerateNumber(ctx context.Context) <-chan int {
	dst := make(chan int)
	n := 1
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("GenerateNumber return")
				return // returning not to leak the goroutine
			case dst <- n:
				fmt.Println("Generate number is", n)
				n++
			}
		}
	}()
	return dst
}
