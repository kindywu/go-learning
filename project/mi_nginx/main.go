package main

import (
	"io"
	"log"
	"main/project/mi_nginx/common"
	"net"
	"sync"
	"time"
)

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

	pool := NewConnectionPool(config.UpstreamAddr, 10)
	// 接受客户端连接
	for {
		client, err := listener.Accept()
		if err != nil {
			log.Println("接受客户端连接失败:", err)
			continue
		}

		// 处理客户端连接
		go handleClient(client, pool)
	}
}

type ConnectionPool struct {
	connections  chan net.Conn
	upstreamAddr string
}

func NewConnectionPool(upstreamAddr string, poolSize int) *ConnectionPool {
	pool := &ConnectionPool{
		connections:  make(chan net.Conn, poolSize),
		upstreamAddr: upstreamAddr,
	}
	// 初始化连接池
	for i := 0; i < poolSize; i++ {
		dialer := net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 10 * time.Minute,
		}
		conn, err := dialer.Dial("tcp", upstreamAddr)
		if err != nil {
			log.Println("无法初始化连接池:", err)
			continue
		}
		log.Println("创建新连接，存入连接池")
		pool.connections <- conn
	}
	return pool
}

func (p *ConnectionPool) GetConnection() (net.Conn, error) {
	select {
	case conn := <-p.connections:
		log.Println("返回连接池里的连接")
		return conn, nil
	default:
		log.Println("连接池为空，创建新连接")
		dialer := net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 10 * time.Minute,
		}
		conn, err := dialer.Dial("tcp", p.upstreamAddr)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

func (p *ConnectionPool) ReleaseConnection(conn net.Conn) {
	select {
	case p.connections <- conn:
		log.Println("连接返回连接池")
	default:
		log.Println("连接池已满，关闭连接")
		conn.Close()
	}
}

//lint:ignore U1000 reason: it's a demo impl to compare with io.Copy
func copyBuffer(src, dst net.Conn) (written int64, err error) {
	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

func handleClient(client net.Conn, pool *ConnectionPool) {
	defer client.Close()

	// 从连接池获取连接
	upstream, err := pool.GetConnection()
	if err != nil {
		log.Println("无法从连接池获取连接", err)
		return
	}
	defer pool.ReleaseConnection(upstream)

	wg := sync.WaitGroup{}
	wg.Add(1)
	// 创建读写管道
	go func() {
		defer wg.Done()
		// 使用自定义的copyBuffer函数复制数据
		log.Println("从client复制到upstream")
		// n, err := copyBuffer(client, upstream)
		n, err := io.Copy(upstream, client)
		if err != nil && err != io.EOF {
			log.Println("从client复制到upstream，拷贝错误:", err)
		}
		log.Println("从client复制到upstream", n)
	}()

	// 也可以从upstream复制到client
	log.Println("从upstream复制到client")
	// m, err := copyBuffer(upstream, client)
	m, err := io.Copy(client, upstream)
	if err != nil && err != io.EOF {
		log.Println("从upstream复制到client，拷贝错误:", err)
	}

	log.Println("从upstream复制到client", m)
	wg.Wait()
}
