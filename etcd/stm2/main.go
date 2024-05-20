package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	endpoints := []string{"etcd1:2379"}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		log.Println("connect failed err: ", err)
		return
	}
	defer client.Close()

	var w sync.WaitGroup
	count := 10
	w.Add(count)
	key10 := "setnx"
	for i := 0; i < count; i++ {
		go func(i int) {
			defer w.Done()
			time.Sleep(5 * time.Millisecond)
			//通过key的Create_Revision 是否为 0 来判断key是否存在。其中If，Then 以及 Else 分支都可以包含多个操作。
			//返回的数据包含一个successed字段，当为 true 时代表 If 为真
			txn0_resp, err := client.Txn(context.TODO()).Then(clientv3.OpGet(key10)).Commit()
			if err != nil {
				log.Fatalln(err)
			} else {
				log.Println("txn0_resp", txn0_resp)
			}

			fromKV := txn0_resp.Responses[0].GetResponseRange().Kvs[0]
			fromValue, _ := strconv.Atoi(string(fromKV.Value))
			txn1_resp, err := client.Txn(context.TODO()).
				If(clientv3.Compare(clientv3.ModRevision(key10), "=", fromKV.ModRevision),
					clientv3.Compare(clientv3.CreateRevision(key10), ">", 0)).
				Then(clientv3.OpPut(key10, fmt.Sprintf("%d", fromValue+i))).Commit()
			if err != nil {
				log.Fatalln(err)
			} else {
				log.Println("txn1_resp", txn1_resp)
			}
		}(i)
	}
	w.Wait()

	if get_resp, err := client.Get(context.TODO(), key10); err != nil {
		log.Println(err)
	} else {
		log.Println("get_resp", get_resp)
	}
}
