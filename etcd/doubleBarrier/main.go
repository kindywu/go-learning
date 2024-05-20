package main

import (
	"bufio"
	"log"
	"math/rand/v2"
	"os"
	"time"

	v3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

func main() {
	// etcd地址
	endpoints := []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"}
	// 生成一个etcd client
	cli, err := v3.New(v3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	count := 10
	for i := range count {
		go func(id int) {
			session, err := concurrency.NewSession(cli)
			if err != nil {
				log.Fatal(err)
			}
			defer session.Close()
			barrierName := "barrier"

			barrier := recipe.NewDoubleBarrier(session, barrierName, count)
			second := rand.IntN(3)
			time.Sleep(time.Duration(second) * time.Second)
			err = barrier.Enter()
			if err != nil {
				log.Fatal(err)
			}
			log.Println("do", id, "wait second", second, time.Now())
			time.Sleep(time.Duration(rand.IntN(5)) * time.Second)
			err = barrier.Leave()
			if err != nil {
				log.Fatal(err)
			}
			log.Println("finish", id, time.Now())
		}(i)
	}

	bufio.NewReader(os.Stdin).ReadLine()
}
