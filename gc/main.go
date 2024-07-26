package main

import (
	"fmt"
	"net/http"
	"sync"
)

const SIZE = 1024 * 5

func indexHandler1(w http.ResponseWriter, r *http.Request) {
	// 当没有业务逻辑，空架子
	// 当业务逻辑，只有栈内存分配，没有堆内存分配，跟GC无关
	// w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, World!")

	// 有些框架的特性里有zero memory allocation，意思就是框架本身没有进行堆内存分配，高并发情况下，框架自身不会触发GC
}

func indexHandler2(w http.ResponseWriter, r *http.Request) {
	// 模拟执行业务
	// 用new/make确保堆内存分配
	var byteArray = make([]byte, SIZE) // 在堆里创建的对象，比如反序列化JSON对象；用完丢弃等GC回收
	len := len(byteArray)              // 业务计算
	fmt.Fprintf(w, "byte array len is %d", len)
}

var pool = sync.Pool{
	New: func() interface{} {
		arr := make([]byte, SIZE)
		return &arr
	},
}

func indexHandler3(w http.ResponseWriter, r *http.Request) {
	// 模拟执行业务
	var byteArray = pool.Get().(*[]byte) // 通过pool从堆里创建的对象
	defer pool.Put(byteArray)            // 用完对象，归还pool，重复使用，使用 defer 确保无论后续代码如何，都会归还对象
	len := len(*byteArray)               // 业务计算
	fmt.Fprintf(w, "byte array len is %d", len)
}

func indexHandler4(w http.ResponseWriter, r *http.Request) {
	// 模拟执行业务
	pool := New(SIZE, 1*1000*1000)
	var byteArray = pool.Get() // 通过pool从堆里创建的对象
	defer pool.Put(byteArray)  // 用完对象，归还pool，重复使用，使用 defer 确保无论后续代码如何，都会归还对象
	len := len(*byteArray)     // 业务计算
	fmt.Fprintf(w, "byte array len is %d", len)
}

func main() {
	http.HandleFunc("/1", indexHandler1)
	http.HandleFunc("/2", indexHandler2)
	http.HandleFunc("/3", indexHandler3)
	http.HandleFunc("/4", indexHandler4)
	http.ListenAndServe(":8080", nil)
}
