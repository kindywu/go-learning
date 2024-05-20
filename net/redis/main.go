package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func readRESPMessage(rd *bufio.Reader) (string, error) {
	// 读取第一个字节以确定消息类型
	firstByte, err := rd.ReadByte()
	if err != nil {
		if err == io.EOF {
			return "", err
		}
		return "", fmt.Errorf("error reading first byte: %w", err)
	}

	switch firstByte {
	case '+': // 简单字符串
		return readSimpleString(rd)
	case '*': // 简单字符串
		return readSimpleString(rd)
	default:
		return "", fmt.Errorf("invalid first byte: %c", firstByte)
	}
}

func readSimpleString(rd *bufio.Reader) (string, error) {
	line, err := rd.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading simple string: %w", err)
	}
	line = strings.TrimSpace(line) // 移除前后的空白字符
	return line[1:], nil           // 移除开头的'+'
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	rd := bufio.NewReader(conn)
	for {
		line, err := readRESPMessage(rd)
		if err != nil {
			log.Println("Error reading from client:", err)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 打印接收到的命令（用于调试）
		fmt.Println("Received command:", line)

		// 根据命令执行相应的操作
		// 这里只是一个示例，实际应用中需要根据具体命令进行处理
		if strings.ToLower(line) == "info" {
			_, err := conn.Write([]byte("+OK\r\n"))
			if err != nil {
				log.Println("Error writing to client:", err)
				return
			}
		} else if strings.ToLower(line) == "ping" {
			_, err := conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				log.Println("Error writing to client:", err)
				return
			}
		} else {
			_, err := conn.Write([]byte("-ERR unknown command\r\n"))
			if err != nil {
				log.Println("Error writing to client:", err)
				return
			}
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("Error listening: %s", err)
	}
	defer listener.Close()

	log.Println("Listening on :6379")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting: %s", err)
			continue
		}

		go handleClient(conn) // 为每个连接创建一个goroutine
	}
}
