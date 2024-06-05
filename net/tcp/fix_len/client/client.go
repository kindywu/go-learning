package main

import (
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"strings"
	"sync"

	_ "go.uber.org/automaxprocs"
)

const MAX_SIZE = 4 * 1024
const CLIENT_SIZE = 100

var writeBufferPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("make bytes buffer")
		buf := make([]byte, MAX_SIZE) // 假设我们预计每次发送的消息的最大长度为4096字节
		return &buf
	},
}

func main() {
	runtime.GOMAXPROCS(16)
	wg := sync.WaitGroup{}
	wg.Add(CLIENT_SIZE)
	for range CLIENT_SIZE {
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
		userMessage := randomString(100)
		if err := write(conn, userMessage); err != nil {
			fmt.Println("发送消息失败:", err)
			break
		}
	}
}

func write(conn net.Conn, message string) error {
	pointer := writeBufferPool.Get().(*[]byte)
	defer writeBufferPool.Put(pointer)

	writeBuf := *pointer

	// 计算消息长度，并写入4个字节的长度字段
	messageLength := len(message)
	if messageLength > MAX_SIZE-4 {
		return fmt.Errorf("消息超过最大长度 %d", MAX_SIZE)
	}

	copy(writeBuf[:4], []byte{byte(messageLength >> 24), byte(messageLength >> 16), byte(messageLength >> 8), byte(messageLength)})

	// 写入用户消息和换行符
	copy(writeBuf[4:4+messageLength], message)

	// 发送消息
	if _, err := conn.Write(writeBuf[:4+messageLength]); err != nil {
		return err
	}

	return nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
