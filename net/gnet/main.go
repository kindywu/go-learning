package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"main/net/gnet/protocol"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

type simpleServer struct {
	gnet.BuiltinEventEngine
	eng          gnet.Engine
	network      string
	addr         string
	multicore    bool
	connected    int32
	disconnected int32
}

func (s *simpleServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logging.Infof("running server on %s with multi-core=%t",
		fmt.Sprintf("%s://%s", s.network, s.addr), s.multicore)
	s.eng = eng
	return
}

func (s *simpleServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	c.SetContext(new(protocol.SimpleCodec))
	atomic.AddInt32(&s.connected, 1)
	out = []byte("sweetness\r\n")
	return
}

func (s *simpleServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Infof("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	disconnected := atomic.AddInt32(&s.disconnected, 1)
	connected := atomic.AddInt32(&s.connected, -1)
	if connected == 0 {
		logging.Infof("all %d connections are closed, shut it down", disconnected)
		action = gnet.Shutdown
	}
	return
}

func (s *simpleServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	codec := c.Context().(*protocol.SimpleCodec)
	var packets [][]byte
	for {
		data, err := codec.Decode(c)
		if err == protocol.ErrIncompletePacket {
			break
		}
		if err != nil {
			logging.Errorf("invalid packet: %v", err)
			return gnet.Close
		}
		packet, _ := codec.Encode(data)
		packets = append(packets, packet)
	}
	if n := len(packets); n > 1 {
		_, _ = c.Writev(packets)
	} else if n == 1 {
		_, _ = c.Write(packets[0])
	}
	return
}

func server() {
	var port int
	var multicore bool

	// Example command: go run server.go --port 9000 --multicore=true
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore=true")
	flag.Parse()
	ss := &simpleServer{
		network:   "tcp",
		addr:      fmt.Sprintf(":%d", port),
		multicore: multicore,
	}
	err := gnet.Run(ss, ss.network+"://"+ss.addr, gnet.WithMulticore(multicore))
	logging.Infof("server exits with error: %v", err)
}

func logErr(err error) {
	logging.Error(err)
	if err != nil {
		panic(err)
	}
}

func client() {
	var (
		network     string
		addr        string
		concurrency int
		packetSize  int
		packetBatch int
		packetCount int
	)

	// Example command: go run client.go --network tcp --address ":9000" --concurrency 100 --packet_size 1024 --packet_batch 20 --packet_count 1000
	flag.StringVar(&network, "network", "tcp", "--network tcp")
	flag.StringVar(&addr, "address", "127.0.0.1:9000", "--address 127.0.0.1:9000")
	flag.IntVar(&concurrency, "concurrency", 10, "--concurrency 500")
	flag.IntVar(&packetSize, "packet_size", 1024, "--packe_size 256")
	flag.IntVar(&packetBatch, "packet_batch", 100, "--packe_batch 100")
	flag.IntVar(&packetCount, "packet_count", 10000, "--packe_count 10000")
	flag.Parse()

	logging.Infof("start %d clients...", concurrency)
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			runClient(network, addr, packetSize, packetBatch, packetCount)
			wg.Done()
		}()
	}
	wg.Wait()
	logging.Infof("all %d clients are done", concurrency)
}

func runClient(network, addr string, packetSize, batch, count int) {

	c, err := net.Dial(network, addr)
	logErr(err)
	logging.Infof("connection=%s starts...", c.LocalAddr().String())
	defer func() {
		logging.Infof("connection=%s stops...", c.LocalAddr().String())
		c.Close()
	}()
	rd := bufio.NewReader(c)
	msg, err := rd.ReadBytes('\n')
	logErr(err)
	expectMsg := "sweetness\r\n"
	if string(msg) != expectMsg {
		logging.Fatalf("the first response packet mismatches, expect: %s, but got: %s", expectMsg, msg)
	}

	for i := 0; i < count; i++ {
		batchSendAndRecv(c, rd, packetSize, batch)
	}
}

func batchSendAndRecv(c net.Conn, rd *bufio.Reader, packetSize, batch int) {
	codec := protocol.SimpleCodec{}
	var (
		requests  [][]byte
		buf       []byte
		packetLen int
	)
	for i := 0; i < batch; i++ {
		req := make([]byte, packetSize)
		_, err := rand.Read(req)
		logErr(err)
		requests = append(requests, req)
		packet, _ := codec.Encode(req)
		packetLen = len(packet)
		buf = append(buf, packet...)
	}
	_, err := c.Write(buf)
	logErr(err)
	respPacket := make([]byte, batch*packetLen)
	_, err = io.ReadFull(rd, respPacket)
	logErr(err)
	for i, req := range requests {
		rsp, err := codec.Unpack(respPacket[i*packetLen:])
		logErr(err)
		if !bytes.Equal(req, rsp) {
			logging.Fatalf("request and response mismatch, conn=%s, packet size: %d, batch: %d",
				c.LocalAddr().String(), packetSize, batch)
		}
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		server()
	}()
	time.Sleep(3 * time.Second)
	go func() {
		defer wg.Done()
		go client()
	}()

	wg.Wait()
}
