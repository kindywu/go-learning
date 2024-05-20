package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// 创建一个管道
	reader, writer, err := os.Pipe()
	if err != nil {
		log.Fatal("Failed to create pipe: ", err)
	}

	// 定义子进程要执行的命令
	cmd := exec.Command("echo", "Hello, world!")

	// 设置子进程的标准输入和输出
	cmd.Stdin = os.Stdin
	cmd.Stdout = writer
	cmd.Stderr = os.Stderr

	// 启动子进程
	err = cmd.Start()
	if err != nil {
		log.Fatal("Failed to start command: ", err)
	}

	// 等待子进程结束
	err = cmd.Wait()
	if err != nil {
		log.Fatal("Command exited with error: ", err)
	}

	// 关闭管道的写端，因为子进程已经启动，不需要再写入数据
	writer.Close()

	// 读取子进程的输出
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// 检查扫描器是否有错误发生
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// 关闭读端管道
	reader.Close()
}
