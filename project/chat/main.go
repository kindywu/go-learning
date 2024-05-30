package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"

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

	p := NewProxy()
	p.Start()
	defer p.Close()
	for {
		client, err := listener.Accept()
		if err != nil {
			log.Fatal("接受客户端失败", err)
		}
		go handle_client(client, &p)
	}
}

func handle_client(conn net.Conn, p *Proxy) {
	defer conn.Close()

	// 创建bufio.Scanner，设置分隔符为回车符
	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanLines)
	scanner.Buffer([]byte{}, bufio.MaxScanTokenSize)

	conn.Write([]byte("Input your name:\n"))
	var name string
	for scanner.Scan() {
		name = scanner.Text()
		if len(name) > 0 {
			break
		}
	}
	log.Printf("%s joined, addr=[%s]", name, conn.RemoteAddr())

	p.Register(name, conn)
	defer log.Printf("%s left, addr=[%s]", name, conn.RemoteAddr())
	defer p.Unregister(name, conn)

	for scanner.Scan() {
		content := scanner.Text()
		p.Broadcast(NewMessage(name, conn.RemoteAddr(), content))
	}
}

type Message struct {
	Sender     string
	SenderAddr net.Addr
	Content    string
}

func NewMessage(sender string,
	senderAddr net.Addr,
	content string) Message {
	return Message{
		Sender:     sender,
		SenderAddr: senderAddr,
		Content:    content,
	}
}

const SYSTEM = "system"

type Proxy struct {
	clients   sync.Map
	broadcast chan Message
}

func NewProxy() Proxy {
	return Proxy{
		clients:   sync.Map{},
		broadcast: make(chan Message, 1024),
	}
}

func (p *Proxy) Start() {
	go func() {
		for message := range p.broadcast {
			p.clients.Range(func(key, value interface{}) bool {
				addr := key.(net.Addr)
				if addr != message.SenderAddr {
					conn := value.(net.Conn)
					conn.Write([]byte(fmt.Sprintf("%s: %s\n", message.Sender, message.Content)))
				}
				return true
			})
		}
	}()
}

func (p *Proxy) Close() {
	close(p.broadcast)
	p = nil
}

func (p *Proxy) Register(name string, conn net.Conn) {
	p.clients.Store(conn.RemoteAddr(), conn)
	p.Broadcast(NewMessage(SYSTEM, conn.RemoteAddr(), fmt.Sprintf("[%s joined]", name)))
}

func (p *Proxy) Unregister(name string, conn net.Conn) {
	p.clients.Delete(conn.RemoteAddr())
	p.Broadcast(NewMessage(SYSTEM, conn.RemoteAddr(), fmt.Sprintf("[%s left]", name)))
}

func (p *Proxy) Broadcast(message Message) {
	p.broadcast <- message
}
