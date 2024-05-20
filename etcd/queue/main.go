package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

func main() {

	// 解析etcd地址
	endpoints := []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"}
	queueName := "queue"

	// 创建etcd的client
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	go func() {
		// 创建/获取队列
		q := recipe.NewQueue(cli, queueName)
		for i := range 10 {
			item := fmt.Sprintf("item %d", i)
			err := q.Enqueue(item)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("enqueue", "<-", item)
		}
	}()

	go func() {
		// 创建/获取队列
		q := recipe.NewQueue(cli, queueName)
		for {
			item, err := q.Dequeue()
			if err != nil {
				log.Println(err)
			}
			fmt.Println("dequeue", "->", item)
		}
	}()

	bufio.NewReader(os.Stdin).ReadLine()
}
