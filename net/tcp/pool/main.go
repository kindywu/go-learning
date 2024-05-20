package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	mrand "math/rand"
	"net"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/panjf2000/ants/v2"
)

const ClientNum = 100
const ClientSendNum = 1000
const ClientSendMessageMinLength = 2 << 10
const ClientSendMessageMaxLength = 4 << 12
const DefaultExpiredTime = 10 * time.Second

var pool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

var lenPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func server(wg *sync.WaitGroup) {
	p, _ := ants.NewMultiPool(10, 5e4/10, ants.RoundRobin, ants.WithExpiryDuration(DefaultExpiredTime))
	defer p.ReleaseTimeout(DefaultExpiredTime) //nolint:errcheck
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

		p.Submit(func() {
			handleClient(wg, conn)
		})
	}
}

func handleClient(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()
	defer func() {
		fmt.Println("HandleClient Exit " + conn.RemoteAddr().String())
		conn.Close()
	}()

	fmt.Println("HandleClient Enter " + conn.RemoteAddr().String())
	for {
		lenBuf := lenPool.Get().(*bytes.Buffer)
		lenBuf.Grow(4)
		// 读取消息长度
		lengthBytes := lenBuf.Bytes()[:4]
		_, err := conn.Read(lengthBytes)
		if err != nil {
			if err == io.EOF {
				//fmt.Println("Client disconnected:", err.Error())
				break
			}
			fmt.Println("Error reading length: ", err.Error())
			break
		}

		// 解码消息长度
		length := int(binary.BigEndian.Uint32(lengthBytes))

		lenBuf.Reset()
		lenPool.Put(lenBuf)
		if length <= 0 {
			fmt.Println("Invalid message length")
			break
		}

		buf := pool.Get().(*bytes.Buffer)
		buf.Grow(length)
		// 读取消息体

		msg := buf.Bytes()[:length]
		_, err = conn.Read(msg)
		if err != nil {
			fmt.Println("Error reading message: ", err.Error())
			break
		}

		// if length > 1000 && length < 2000 {
		// 	fmt.Println("Received message", string(msg[:10]), "length", length)
		// }

		buf.Reset()
		pool.Put(buf)

		// 处理消息
		// fmt.Printf("Received message length: %d\n", length)
		// fmt.Printf("Received message: %s\n", msg)
	}
}

const CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

	for range ClientSendNum {
		// 发送消息
		// msg := fmt.Sprintf("%s %d", "Hello, Server!", i)
		length := generateRandomInt(ClientSendMessageMinLength, ClientSendMessageMaxLength)

		lenBuf := lenPool.Get().(*bytes.Buffer)
		lenBuf.Grow(4)
		lengthBytes := lenBuf.Bytes()[:4]
		binary.BigEndian.PutUint32(lengthBytes, uint32(length))

		// 写入长度和消息
		_, err = conn.Write(lengthBytes)
		lenBuf.Reset()
		lenPool.Put(lenBuf)
		if err != nil {
			fmt.Println("Error writing length: ", err.Error())
			return
		}

		buf := pool.Get().(*bytes.Buffer)
		buf.Grow(length)
		for i := 0; i < length; i++ {
			// 随机选择一个字符
			randInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(CHARSET))))
			if err != nil {
				buf.Reset()
				pool.Put(buf) // 释放资源
				fmt.Println("Error writing length: ", err.Error())
				return
			}
			buf.WriteByte(CHARSET[randInt.Int64()])
		}

		// fmt.Println("Send message (string)", buf.String())
		// fmt.Println("Send message (string)", buf.Next(length), "length", length)
		_, err = conn.Write(buf.Next(length))
		if err != nil {
			fmt.Println("Error writing message: ", err.Error())
			return
		}
		buf.Reset()
		pool.Put(buf)
		// fmt.Println("Message sent to server.")
	}

	// println("client finish")
}

func main() {
	defer ants.Release()
	var wg sync.WaitGroup
	go func() {
		server(&wg)
	}()

	time.Sleep(3 * time.Second)

	for range ClientNum {
		wg.Add(2) //1 for client,1 for server
		ants.Submit(func() {
			client(&wg)
		})
	}

	wg.Wait()
	println("Finish: press Ctrl+C to exit")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-c
}
