package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	closing := make(chan struct{})
	closed := make(chan struct{})

	go func() {
		for {
			select {
			case <-closing:
				return
			default:
				println("doing task")
				time.Sleep(time.Second)
			}
		}
	}()

	tmpChan := make(chan os.Signal, 1)
	signal.Notify(tmpChan, syscall.SIGTERM, syscall.SIGINT)
	<-tmpChan

	close(closing)

	go doCleanup(closed)

	select {
	case <-closed:
	case <-time.After(3 * time.Second):
		println("超时，不等了")
	}

	println("优雅退出")
}

func doCleanup(closed chan struct{}) {
	time.Sleep(2 * time.Second)
	println("清理干净")
	close(closed)
}
