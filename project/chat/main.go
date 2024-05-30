package main

import (
	"log"
	"net"

	toml "github.com/pelletier/go-toml"
)

func main() {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		log.Fatal("读取配置失败", err)
	}
	listen_addr := config.Get("listen_addr").(string)
	listener, err := net.Listen("tcp", listen_addr)
	if err != nil {
		log.Fatal("监听端口失败", err)
	}
	log.Println("服务启动成功", listen_addr)
	defer listener.Close()
}
