package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("make bytes buffer")
		buf := make([]byte, 4096) // 假设我们预计每次发送的消息的最大长度为4096字节
		return &buf
	},
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// 从Pool中获取一个缓冲区
		pointer := bufferPool.Get().(*[]byte)
		buf := *pointer

		// 读取消息长度
		if _, err := io.ReadFull(reader, buf[:4]); err != nil {
			log.Println("读取消息长度失败:", err)
			return
		}

		// 转换长度字段为整数
		length := int(int32(buf[0])<<24 | int32(buf[1])<<16 | int32(buf[2])<<8 | int32(buf[3]))

		// 确保缓冲区足够大
		if length > cap(buf) {
			log.Println("消息长度超出缓冲区容量")
			return
		}

		// 读取消息内容
		if _, err := io.ReadFull(reader, buf[:length]); err != nil {
			log.Println("读取消息内容失败:", err)
			return
		}

		// 处理接收到的消息
		// fmt.Println("服务器接收到消息:", string(buf[:length]))
		if string(buf[:length]) != "你好" {
			fmt.Printf("len=%d,msg=%s\n", length, string(buf[:length]))
			break
		}

		// 将缓冲区放回Pool
		bufferPool.Put(&buf)
	}
}

func main() {
	// 监听本地8080端口
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("监听端口失败:", err)
	}
	defer listener.Close()

	fmt.Println("服务器正在监听8080端口...")
	for {
		// 接受客户端连接
		conn, err := listener.Accept()
		if err != nil {
			log.Println("接受客户端连接失败:", err)
			continue
		}

		// 处理客户端连接
		go handleClient(conn)
	}
}
