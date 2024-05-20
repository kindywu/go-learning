package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	var buf [1024]byte
	for {
		n, err := reader.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		println("收到client发来的数据：", recvStr)
	}
}

func server() {
	isClosed := false
	listen, err := net.Listen("tcp", "127.0.0.1:30000")

	if err != nil {
		println("listen failed, err:", err)
		return
	}
	defer func() {
		if !isClosed {
			if err := listen.Close(); err != nil {
				println("error closing listener:", err)
			}
		}
	}()

	// 设置5秒后关闭监听器的计时器
	time.AfterFunc(5*time.Second, func() {
		if err := listen.Close(); err != nil {
			println("Error closing listener:", err)
		}
		isClosed = true
		println("5s timeout. Server is shutting down.")
	})

	// 无限循环，接受连接直到监听器关闭
	for {
		conn, err := listen.Accept()
		if err != nil {
			if err != net.ErrClosed {
				println("accept failed, err:", err)
			}
			break // 跳出循环，因为监听器已经关闭
		}
		go process(conn)
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		server()
		println("\nserver exit")
		wg.Done()
	}()

	go func() {
		time.Sleep(3 * time.Second)
		//client
		conn, err := net.Dial("tcp", "127.0.0.1:30000")
		if err != nil {
			println("dial failed, err", err)
			return
		}
		defer conn.Close()
		for i := 0; i < 20; i++ {
			msg := `Hello, Hello. How are you?`
			fmt.Println(msg)
			conn.Write([]byte(msg))
		}
		println("\nclient exit")
		wg.Done()
	}()
	wg.Wait()
}
