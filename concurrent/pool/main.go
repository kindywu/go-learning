package main

import (
	"bytes"
	"fmt"
	"sync"
)

var buffers = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	buffers.Put(buf)
}

func main() {
	// 创建一个 sync.Pool 用于存储和重用 bytes.Buffer 对象
	var bufferPool = sync.Pool{
		New: func() interface{} {
			fmt.Println("New bytes.Buffer")
			return new(bytes.Buffer)
		},
	}

	// 从池中获取一个 bytes.Buffer 对象
	buf := bufferPool.Get().(*bytes.Buffer)

	// 使用 buf 对象...
	buf.WriteString("Hello, World!")

	// 打印 buf 的内容
	fmt.Println(buf.String())

	buf.Reset()
	// 使用完毕后，将 buf 放回池中
	bufferPool.Put(buf)

	// 再次从池中获取一个 bytes.Buffer 对象
	buf = bufferPool.Get().(*bytes.Buffer)

	// 由于我们之前已经放回了一个对象，这里会重用之前的 buf
	// 此时 buf 应该是一个空的 bytes.Buffer，因为没有调用 Put 之前的内容已经被清空
	buf.WriteString("Reused buffer")
	fmt.Println(buf.String())

	// 再次放回池中
	bufferPool.Put(buf)
}

// 程序输出：
// Hello, World!
// Reused buffer
