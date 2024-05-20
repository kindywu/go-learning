package main

import (
	"log"

	clientv2 "go.etcd.io/etcd/client/v2"
)

func main() {
	// ctx := context.Background()
	endpoints := []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"}
	cli, err := clientv2.New(clientv2.Config{
		// Endpoints: []string{"etcd:2379"},
		Endpoints: endpoints,
	})
	if err != nil {
		log.Fatalln(err)
	}

	println(cli)
}
