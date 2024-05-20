package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	endpoints := []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	})
	if err != nil {
		fmt.Println("connect failed err : ", err)
		return
	}
	defer client.Close()

	key := "/demo/demo_key"
	client.Put(context.TODO(), key, "value1")

	go func() {
		//watch
		watchKey := client.Watch(context.TODO(), key)
		for resp := range watchKey {
			for _, event := range resp.Events {
				fmt.Printf("watch %s %q : %q \n", event.Type, event.Kv.Key, event.Kv.Value)
			}
		}
	}()

	log.Println("press enter to continue")
	bufio.NewReader(os.Stdin).ReadLine()
	log.Println("put to", key)
	if resp, err := client.Put(context.TODO(), key, "value2"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("put", resp)
	}
	bufio.NewReader(os.Stdin).ReadLine()
}
