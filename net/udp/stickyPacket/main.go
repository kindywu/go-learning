package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
)

// 定义一个简单的协议，用于处理粘包半包问题
type UDPPacket struct {
	Length uint32 // 数据长度
	Data   []byte // 实际数据
}

// 读取UDP数据并处理粘包半包问题
func readPacket(conn *net.UDPConn) (UDPPacket, error) {
	var packet UDPPacket
	buf := make([]byte, 4) // 用于读取长度字段

	// 首先读取长度字段
	n, err := conn.Read(buf)
	if err != nil {
		return packet, err
	}
	if n != 4 {
		return packet, fmt.Errorf("invalid packet size: %d", n)
	}

	// 解析长度字段
	if err := binary.Read(bufio.NewReader(bufio.NewReader(conn)), binary.BigEndian, &packet.Length); err != nil {
		return packet, err
	}

	// 读取实际数据
	packet.Data = make([]byte, packet.Length)
	n, err = conn.Read(packet.Data)
	if err != nil {
		return packet, err
	}
	if int(packet.Length) != n {
		return packet, fmt.Errorf("incomplete packet: expected %d bytes, got %d", packet.Length, n)
	}

	return packet, nil
}

// 处理接收到的数据包
func handlePacket(packet UDPPacket) {
	fmt.Printf("Received packet: %s\n", packet.Data)
}

// UDP服务器主函数
func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 8080})
	if err != nil {
		fmt.Println("Failed to listen:", err)
		os.Exit(1)
	}
	defer conn.Close() // 确保在程序结束时关闭连接

	var wg sync.WaitGroup // 使用WaitGroup来等待所有Goroutine完成

	for {
		// 使用Goroutine来并发读取数据
		wg.Add(1)
		go func() {
			defer wg.Done()

			// 读取数据包
			packet, err := readPacket(conn)
			if err != nil {
				fmt.Println("Error reading packet:", err)
				return
			}

			// 处理数据包
			handlePacket(packet)
		}()
	}
}
