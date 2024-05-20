package main

import (
	"bufio"
	"log"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

func main() {
	// etcd地址
	endpoints := []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"}
	// 生成一个etcd client
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 创建/获取栅栏
	barrierName := "barrier"
	barrier := recipe.NewBarrier(cli, barrierName)

	go func() {
		barrier.Hold()
		log.Println("hold")
		bufio.NewReader(os.Stdin).ReadLine()
		log.Println("release")
		barrier.Release()
	}()

	time.Sleep(3 * time.Second)

	for i := range 10 {
		go func(id int) {
			log.Println("ready")
			err = barrier.Wait()
			if err != nil {
				log.Fatal(err)
			}
			log.Println("go", id)
		}(i)
	}

	time.Sleep(10 * time.Second)

}
