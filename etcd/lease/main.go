package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// From：https://github.com/the-gigi/go-etcd3-demo/blob/master/main.go
var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

func main() {
	endpoints := []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	client, err := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   endpoints,
	})

	if err != nil {
		log.Fatal(err)
	}

	kv := clientv3.NewKV(client)

	key := "/demo/demo1_key"
	//Delete all keys
	kv.Delete(ctx, key, clientv3.WithPrefix())

	gr, _ := kv.Get(ctx, key)
	if len(gr.Kvs) == 0 {
		fmt.Println("no key")
	}

	lease, err := client.Grant(ctx, 3)
	if err != nil {
		log.Fatal(err)
	}

	//Insert key with a lease of 3 second TTL
	kv.Put(ctx, key, "demo1_value", clientv3.WithLease(lease.ID))

	gr, _ = kv.Get(ctx, key)
	if len(gr.Kvs) == 1 {
		fmt.Println("Found key")
	}

	//let the TTL expire
	time.Sleep(5 * time.Second)

	gr, _ = kv.Get(ctx, key)
	if len(gr.Kvs) == 0 {
		fmt.Println("no more key")
	}
}
