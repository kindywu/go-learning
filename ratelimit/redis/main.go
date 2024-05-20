package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	var wg sync.WaitGroup
	wg.Add(2)

	for i := 0; i < 20; i++ {
		go func(id int) {
			defer wg.Done()

			ctx := context.Background()
			rdb := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: []string{
					"redis-cluster-redis-node-1-1:6379",
					"redis-cluster-redis-node-2-1:6379",
					"redis-cluster-redis-node-3-1:6379",
					"redis-cluster-redis-node-4-1:6379",
					"redis-cluster-redis-node-5-1:6379"},
				Password: "bitnami",
			})
			_ = rdb.FlushDB(ctx).Err()

			limiter := redis_rate.NewLimiter(rdb)
			for j := 0; j < 10; j++ {
				res, err := limiter.Allow(ctx, "token:123", redis_rate.PerSecond(1))
				if err != nil {
					panic(err)
				}

				if res.Allowed == 0 {
					//log.Println("g:", id, "allowed", res.Allowed, "remaining", res.Remaining, "retry after", res.ResetAfter)
					time.Sleep(res.ResetAfter)
				} else {
					log.Println("g:", id, "allowed", res.Allowed, "time", time.Now())
				}
			}
		}(i)
	}

	wg.Wait()
}
