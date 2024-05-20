package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

var (
	c       = flag.Int("c", 10, "concurrency")
	prefork = flag.Bool("prefork", false, "prefork")
	child   = flag.Bool("child", false, "child proc")
)

const Host = ":9876"

// go run main.go -prefork -c 5

// ps -a
// pstree -c -p 48321

// apt update && apt -y install ncat
// nc 127.0.0.1 9876

func main() {
	flag.Parse()

	var ln net.Listener
	var err error

	if *prefork {
		ln = doPrefork(*c)
	} else {
		fmt.Println("Single")
		ln, err = net.Listen("tcp", Host)
		if err != nil {
			panic(err)
		}
	}

	start(ln)
}

func start(ln net.Listener) {
	log.Println("child started")
	for {
		conn, err := ln.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok {
				log.Printf("accept net err: %v\n", nerr)
				continue
			}

			log.Printf("accept err:%v\n", err)
			return
		}

		go io.Copy(conn, conn)
	}
}

func doPrefork(c int) net.Listener {
	var listener net.Listener
	var err error
	if !*child {
		log.Println("main started")
		// 主进程
		addr, err := net.ResolveTCPAddr("tcp", Host)
		if err != nil {
			log.Fatal(err)
		}
		tcpListener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		fl, err := tcpListener.File()
		if err != nil {
			log.Fatal(err)
		}

		children := make([]*exec.Cmd, c)
		for i := range children {
			children[i] = exec.Command(os.Args[0], "-prefork", "-child")
			children[i].Stdout = os.Stdout
			children[i].Stderr = os.Stderr
			children[i].ExtraFiles = []*os.File{fl}

			err = children[i].Start()
			if err != nil {
				log.Fatal(err)
			}
		}

		for _, child := range children {
			if err := child.Wait(); err != nil {
				log.Printf("failed to wait child's starting: %v", err)
			}
		}
		os.Exit(0)
	} else {
		// 通常文件描述符0是stdin, 1是stdout, 2是stderr, 3是ExtraFiles[0]
		listener, err = net.FileListener(os.NewFile(3, ""))
		if err != nil {
			log.Fatal(err)
		}
	}
	return listener
}
