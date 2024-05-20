package main

import (
	"bufio"
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
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
	useLock(cli) // 测试锁
}

func useLock(cli *clientv3.Client) {
	ctx := context.Background()
	ttl := 3
	count := 0
	for i := range 10 {
		go func(id int) {
			// 为锁生成session
			s1, err := concurrency.NewSession(cli, concurrency.WithTTL(ttl))
			if err != nil {
				log.Fatal(err)
			}
			defer s1.Close()
			// lockName := fmt.Sprintf("lock-%d", id)
			lockName := "lock"
			//得到一个分布式锁
			locker := concurrency.NewMutex(s1, lockName)

			// 请求锁
			log.Println("acquiring lock", id)
			locker.Lock(ctx)
			log.Println("acquired lock", id)

			// 等待一段时间
			second := rand.Intn(10)
			log.Println("do", id, "wait", second)
			time.Sleep(time.Duration(second) * time.Second)

			count++

			locker.Unlock(ctx) // 释放锁
			log.Println("released lock", id)

		}(i)
	}

	bufio.NewReader(os.Stdin).ReadString('\n')
	println(count)
	bufio.NewReader(os.Stdin).ReadString('\n')

}
