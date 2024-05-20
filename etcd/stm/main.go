package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"

	v3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var totalAccounts = 5

func main() {
	// 可能出错mvcc: required revision is a future revision
	endpoints := []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"}
	// endpoints := []string{"etcd1:2379", "etcd2:2379"}
	// endpoints := []string{"etcd1:2379"}

	cli, err := v3.New(v3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 设置5个账户，每个账号都有100元，总共500元

	v := 100
	sum := 0
	for i := 0; i < totalAccounts; i++ {
		k := fmt.Sprintf("accts/%d", i)
		if _, err = cli.Put(context.TODO(), k, fmt.Sprintf("%d", v)); err != nil {
			log.Fatal(err)
		}
		log.Printf("account %s: %d", k, v)
		sum += v
	}
	log.Println("account sum is", sum) // 总数

	fmt.Println("STM")

	count := 100
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := range count {
		go func(id int) {
			defer wg.Done()
			log.Println("stm", id)
			// 随机得到两个转账账号
			from, to := rand.IntN(totalAccounts), rand.IntN(totalAccounts)
			if from == to {
				// 自己不和自己转账
				log.Println("skip")
				return
			}
			// 读取账号的值
			fromK, toK := fmt.Sprintf("accts/%d", from), fmt.Sprintf("accts/%d", to)
			xfer := rand.IntN(50)
			exchange := createExchange(fromK, toK, xfer)
			if _, serr := concurrency.NewSTM(cli, exchange); serr != nil {
				log.Println(serr)
			}
		}(i)
	}

	wg.Wait()
	// 检查账号最后的数目
	sum = 0
	accts, err := cli.Get(context.TODO(), "accts/", v3.WithPrefix()) // 得到所有账号
	if err != nil {
		log.Fatal(err)
	}
	for _, kv := range accts.Kvs { // 遍历账号的值
		v := 0
		fmt.Sscanf(string(kv.Value), "%d", &v)
		sum += v
		log.Printf("account %s: %d", kv.Key, v)
	}

	log.Println("account sum is", sum) // 总数
}

func createExchange(fromK, toK string, xfer int) func(stm concurrency.STM) error {
	exchange := func(stm concurrency.STM) error {
		fromV, toV := stm.Get(fromK), stm.Get(toK)
		fromInt, toInt := 0, 0
		fmt.Sscanf(fromV, "%d", &fromInt)
		fmt.Sscanf(toV, "%d", &toInt)

		if fromInt < xfer {
			return fmt.Errorf("account %s remain %d is not enought for the xfer %d", fromK, fromInt, xfer)
		}
		fromInt, toInt = fromInt-xfer, toInt+xfer

		// 把转账后的值写回
		stm.Put(fromK, fmt.Sprintf("%d", fromInt))
		stm.Put(toK, fmt.Sprintf("%d", toInt))
		return nil
	}
	return exchange
}
