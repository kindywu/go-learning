package main

import (
	"fmt"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

var queue *recipe.Queue

func TestMain(m *testing.M) {
	// 创建etcd的client
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	// 创建/获取队列
	queue = recipe.NewQueue(cli, queueName)

	println(m.Run())
}

// 解析etcd地址
var endpoints = []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"}
var queueName = "queue"

// go test -benchmem -benchtime=10s -run=^$ -bench ^BenchmarkQueue$ main/etcd/queue
func BenchmarkQueue(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			item := fmt.Sprintf("item %d", b.N)
			err := queue.Enqueue(item)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
