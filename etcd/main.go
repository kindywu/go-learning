package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	ctx := context.Background()
	endpoints := []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"}
	cli, err := clientv3.New(clientv3.Config{
		// Endpoints: []string{"etcd:2379"},
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer cli.Close()

	for _, endpoint := range cli.Endpoints() {
		s_resp, err := cli.Status(ctx, endpoint)
		if err != nil {
			log.Fatalln(err)
		}
		if s_resp.Header.MemberId == s_resp.Leader {
			println("etcd leader is", s_resp.Leader, endpoint)
		}
	}

	session, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatalln(err)
	}
	defer session.Close()

	election := concurrency.NewElection(session, "mysql")

	// Elect a leader (or wait that the leader resign)
	node := "node-name"
	if err := election.Campaign(ctx, node); err != nil {
		log.Fatal(err)
	}
	fmt.Println("leader election for", node)

	g_resp, err := election.Leader(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("leader is ", g_resp)
	rev := election.Rev()
	fmt.Println("current rev:", rev)
}
