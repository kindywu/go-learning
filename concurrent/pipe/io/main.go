package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	// 创建一个管道
	reader, writer := io.Pipe()

	// 启动一个goroutine来读取管道中的数据
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(reader)
		if scanner.Scan() {
			fmt.Println("Received:", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal("Reading error:", err)
		}
	}()

	// 写入数据到管道
	fmt.Println("Writing to the pipe...")
	_, err := writer.Write([]byte("Hello, pipe!"))
	if err != nil {
		log.Fatal("Write error:", err)
	}

	// 关闭写入端，表示没有更多的数据会被写入
	err = writer.Close()
	if err != nil {
		log.Fatal("Write close error:", err)
	}

	wg.Wait()

	// 等待goroutine完成读取操作
	reader.Close() // 关闭读取端也会触发goroutine的退出
}
