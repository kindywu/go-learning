package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func doCleanup() {
	println("清理")
}

func main() {
	go func() {
		ch := time.Tick(time.Second)
		for d := range ch {
			fmt.Printf("time %v\n", d)
		}
	}()

	tempChan := make(chan os.Signal, 1)
	signal.Notify(tempChan, syscall.SIGINT, syscall.SIGTERM)
	<-tempChan

	doCleanup()
	fmt.Println("优雅退出")
}
