package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"sync"
	"time"

	_ "go.uber.org/automaxprocs"
)

const K = 1024
const MAX_SIZE = 4 * K

var bufferPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("make bytes buffer")
		buf := make([]byte, MAX_SIZE) // 假设我们预计每次发送的消息的最大长度为4096字节
		return &buf
	},
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// _, err := read(reader, process)
		data, err := read(reader, process)
		if str, ok := data.(string); ok {
			_ = len(str)
		}
		if err != nil {
			if err != io.EOF {
				fmt.Printf("读取数据失败：%v 时间：%v\r\n", err, time.Now())
			}
			return
		}

		// println(result.(string))
	}
}

// func process(buf []byte) (interface{}, error) {
// 	return len(buf), nil
// }

func process(buf []byte) (interface{}, error) {
	return string(buf), nil
}

func read(reader *bufio.Reader, process func(buf []byte) (interface{}, error)) (interface{}, error) {
	pointer := bufferPool.Get().(*[]byte)
	defer bufferPool.Put(pointer)

	buf := *pointer

	// 读取消息长度
	if _, err := io.ReadFull(reader, buf[:4]); err != nil {
		return nil, err
	}

	// 转换长度字段为整数
	length := int(int32(buf[0])<<24 | int32(buf[1])<<16 | int32(buf[2])<<8 | int32(buf[3]))

	// 确保缓冲区足够大
	if length > MAX_SIZE-4 {
		return nil, errors.New("消息长度超出缓冲区容量")
	}

	// 读取消息内容
	if _, err := io.ReadFull(reader, buf[:length]); err != nil {
		return nil, err
	}

	result, err := process(buf[:length])

	if err != nil {
		return nil, err
	}

	return result, nil
}

func main() {
	println(runtime.GOMAXPROCS(16))
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
