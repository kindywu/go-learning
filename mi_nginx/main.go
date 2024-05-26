package main

import (
	"io"
	"log"
	"main/mi_nginx/common"
	"net"
)

func handleClient(client net.Conn, config *common.Config) {
	defer client.Close()

	// 连接目标服务器
	upstream, err := net.Dial("tcp", config.UpstreamAddr)
	if err != nil {
		log.Println("连接服务器失败:", err)
		return
	}
	defer upstream.Close()

	// 创建读写管道
	go io.Copy(upstream, client) // 将客户端数据转发到服务器
	io.Copy(client, upstream)    // 将服务器数据转发到客户端
}

func main() {
	config, err := common.ReadConfig("./config.yml")
	if err != nil {
		log.Fatal("读取配置失败:", err)
	}
	listener, err := net.Listen("tcp", config.ListenAddr)
	if err != nil {
		log.Fatal("监听端口失败:", err)
	}
	defer listener.Close()

	log.Printf("TCP代理服务器正在监听%s...", config.ListenAddr)

	// 接受客户端连接
	for {
		client, err := listener.Accept()
		if err != nil {
			log.Println("接受客户端连接失败:", err)
			continue
		}

		// 处理客户端连接
		go handleClient(client, config)
	}
}
