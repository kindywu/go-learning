package main

import (
	"fmt"
	"net"
	"sync"
)

var writeBufferPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("make bytes buffer")
		buf := make([]byte, 4096) // 假设我们预计每次发送的消息的最大长度为4096字节
		return &buf
	},
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for range 10 {
		go func() {
			defer wg.Done()
			client()
		}()
	}
	wg.Wait()
}

func client() {
	// 连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	defer conn.Close()

	for {
		userMessage := "你好"
		// 从Pool中获取一个写缓冲区
		pointer := writeBufferPool.Get().(*[]byte)
		writeBuf := *pointer

		// 计算消息长度，并写入4个字节的长度字段
		messageLength := len(userMessage) // 加1是为了换行符
		copy(writeBuf[:4], []byte{byte(messageLength >> 24), byte(messageLength >> 16), byte(messageLength >> 8), byte(messageLength)})

		// 写入用户消息和换行符
		copy(writeBuf[4:4+messageLength], userMessage)

		// 发送消息
		if _, err := conn.Write(writeBuf[:4+messageLength]); err != nil {
			fmt.Println("发送消息失败:", err)
			return
		}

		// 将缓冲区放回Pool
		writeBufferPool.Put(&writeBuf)
	}
}
