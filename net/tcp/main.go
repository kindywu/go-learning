package main

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	mrand "math/rand"
	"net"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func server() {

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		fmt.Println("HandleClient Exit " + conn.RemoteAddr().String())
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	fmt.Println("HandleClient Enter " + conn.RemoteAddr().String())
	for {
		// 读取消息长度
		lengthBytes := make([]byte, 4)
		_, err := reader.Read(lengthBytes)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected:", err.Error())
				break
			}
			fmt.Println("Error reading length: ", err.Error())
			break
		}

		// 解码消息长度
		length := int(binary.BigEndian.Uint32(lengthBytes))
		if length <= 0 {
			fmt.Println("Invalid message length")
			break
		}

		// 读取消息体
		msg := make([]byte, length)
		_, err = reader.Read(msg)
		if err != nil {
			fmt.Println("Error reading message: ", err.Error())
			break
		}

		// 处理消息
		// fmt.Printf("Received message length: %d\n", length)
		// fmt.Printf("Received message: %s\n", msg)
	}
}

// generateRandomString 生成一个指定长度的随机字符串
func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var bytes = make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes), nil
}

func generateRandomInt(min, max int) int {
	// 创建一个随机数源，使用当前时间作为种子
	source := mrand.NewSource(time.Now().UnixNano())

	// 使用随机数源创建一个随机数生成器
	randomGenerator := mrand.New(source)

	// 生成随机整数
	return randomGenerator.Intn(max-min) + min
}

func client(wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for range 10 {
		// 发送消息
		// msg := fmt.Sprintf("%s %d", "Hello, Server!", i)
		msg, _ := generateRandomString(generateRandomInt(10, 100000))
		// 计算长度并编码
		length := len(msg)
		lengthBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(lengthBytes, uint32(length))

		// 写入长度和消息
		_, err = conn.Write(lengthBytes)
		if err != nil {
			fmt.Println("Error writing length: ", err.Error())
			return
		}

		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error writing message: ", err.Error())
			return
		}

		// fmt.Println("Message sent to server.")
	}

	println("client finish")
}

func main() {
	go func() {
		server()
	}()

	time.Sleep(3 * time.Second)
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go client(&wg)
	}

	wg.Wait()
	println("Finish: press Ctrl+C to exit")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-c
}
