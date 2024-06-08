package main

import (
	"fmt"
	"net"
	"testing"
)

// go test -bench=BenchmarkClientSend -benchmem -cpu 1,2,4,8,16
func BenchmarkClientSend(b *testing.B) {
	// message_size := 512
	message_size := 2 * K
	message := randomString(message_size)

	b.RunParallel(func(pb *testing.PB) {
		conn, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			fmt.Println("连接服务器失败:", err)
			return
		}
		defer conn.Close()
		for pb.Next() {
			if err := write(conn, message); err != nil {
				fmt.Println("发送消息失败:", err)
			}
		}
	})
}
