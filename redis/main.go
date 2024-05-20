package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func main() {
	var ctx = context.Background()

	// 配置主节点和从节点的连接信息
	primaryRedisCfg := &redis.Options{
		Addr:     "redis-redis-primary-1:6379",
		Password: "my_password", // 根据你的配置文件中的密码进行修改
		DB:       0,             // 使用默认数据库
	}

	secondaryRedisCfg := &redis.Options{
		Addr:     "redis-redis-secondary-1:6379",
		Password: "my_password", // 同上
		DB:       0,             // 使用默认数据库
	}

	// 连接到主节点
	primaryClient := redis.NewClient(primaryRedisCfg)
	defer primaryClient.Close()

	// 连接到从节点
	secondaryClient := redis.NewClient(secondaryRedisCfg)
	defer secondaryClient.Close()

	// 写入数据到主节点
	setResult := primaryClient.Set(ctx, "mykey", "myvalue", 0) // 设置键为"mykey"，值为"myvalue"
	if setResult.Err() != nil {
		log.Fatalf("Error setting key in primary Redis: %v", setResult.Err())
	}
	fmt.Printf("Set operation result: %v\n", setResult)

	// 从从节点读取数据
	getResult, err := secondaryClient.Get(ctx, "mykey").Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("Key not found in secondary Redis, it might not have replicated yet")
		} else {
			log.Fatalf("Error getting key from secondary Redis: %v", err)
		}
	} else {
		fmt.Printf("Get operation result: %s\n", getResult)
	}
}
