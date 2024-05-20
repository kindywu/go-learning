package main

import (
	"fmt"
	"time"
)

func asStream(done <-chan struct{}, values ...interface{}) <-chan interface{} {
	s := make(chan interface{}, 3) //创建一个unbuffered的channel
	go func() {                    // 启动一个goroutine，往s中塞数据
		defer close(s)             // 退出时关闭chan
		for _, v := range values { // 遍历数组
			select {
			case <-done:
				return
			case s <- v:
				println("s<-v") // 将数组元素塞入到chan中
			}
		}
	}()
	return s
}

func takeN(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{}) // 创建输出流
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ { // 只读取前num个元素
			select {
			case <-done:
				return
			case v, ok := <-valueStream:
				if !ok {
					return
				}
				fmt.Printf("takeStream<-v ok=%v\n", ok)
				takeStream <- v
			}
		}
	}()
	return takeStream
}
func main() {
	done := make(chan struct{})
	// 设置一个3秒后的定时器
	time.AfterFunc(20*time.Second, func() {
		// 向done通道发送一个空的结构体，表示时间到了
		// done <- struct{}{}
		close(done)
	})
	valueStream := asStream(done, 0, 1, 2, 3, 4, 5, 6, 7)

	time.Sleep(3 * time.Second)
	fmt.Println("开始")
	// valueStream := asStream(done, 0, 1, 2, 3)
	takenStream := takeN(done, valueStream, 5)

	for v := range takenStream {
		intValue, ok := v.(int)

		// 检查类型断言是否成功
		if ok {
			time.Sleep(time.Second)
			fmt.Println("The value is an int:", intValue)
		} else {
			fmt.Println("The value is NOT an int.")
		}
	}
}
