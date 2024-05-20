package main

import (
	"fmt"
	"net"
	"os"
)

func server() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9981})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Local: <%s> \n", listener.LocalAddr().String())
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during server read: %s", err)
		}
		fmt.Printf("read %s from <%s>\n", data[:n], remoteAddr)
		_, err = listener.WriteToUDP([]byte("world"), remoteAddr)
		if err != nil {
			fmt.Printf("error during server write: %s", err)
		}
	}
}

func client() {
	sip := net.ParseIP("127.0.0.1")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: sip, Port: 9981}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte("hello"))
	if err != nil {
		fmt.Printf("error during client write: %s", err)
	}
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Printf("error during client read: %s", err)
	}
	fmt.Printf("read %s from <%s>\n", data[:n], conn.RemoteAddr())
}

func main() {
	go func() {
		server()
	}()
	go func() {
		client()
	}()

	b := make([]byte, 1)
	os.Stdin.Read(b)
}
