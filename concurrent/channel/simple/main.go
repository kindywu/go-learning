package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ch := make(chan bool)
	ch2 := make(chan bool)
	ch3 := make(chan bool)

	go func() {
		for {
			time.Sleep(time.Second)
			println(1)
			ch <- true
			<-ch3
		}
	}()

	go func() {
		for {
			<-ch
			time.Sleep(time.Second)
			println(2)
			ch2 <- true
		}
	}()

	go func() {
		for {
			<-ch2
			time.Sleep(time.Second)
			println(3)
			ch3 <- true
		}
	}()

	// 处理CTRL+C等中断信号
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	println("Quit")
}
