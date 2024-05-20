package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	fmt.Printf("DOCKER_HOST :%s\n", os.Getenv("DOCKER_HOST"))
	fmt.Printf("DOCKER_API_VERSION :%s\n", os.Getenv("DOCKER_API_VERSION"))
	fmt.Printf("DOCKER_CERT_PATH :%s\n", os.Getenv("DOCKER_CERT_PATH"))
	fmt.Printf("DOCKER_TLS_VERIFY :%s\n", os.Getenv("DOCKER_TLS_VERIFY"))
	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// 获取所有容器的列表
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	// 打印容器信息
	for _, container := range containers {
		fmt.Printf("Container ID: %s, Image: %s, Status: %s\n", container.ID, container.Image, container.Status)
	}
}
