package main

import (
	"fmt"
	"io"
	"os"
)

// loggingWriter 是一个实现了 io.Writer 接口的结构体
type loggingWriter struct {
	io.Writer
}

// Write 方法实现了 io.Writer 接口的 Write 方法
func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Got '%s'\n", string(b)) // 打印接收到的数据
	return w.Writer.Write(b)            // 调用底层 Writer 的 Write 方法
}

func main() {
	// 创建一个 loggingWriter 实例，底层 Writer 是 os.Stdout
	logger := loggingWriter{Writer: os.Stdout}

	// 使用 loggingWriter 实例来写入数据
	_, err := logger.Write([]byte("Hello, World!"))
	if err != nil {
		fmt.Println("Error writing to logger:", err)
	}
}
