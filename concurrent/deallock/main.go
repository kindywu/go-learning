package main

// import (
// 	"fmt"
// 	"sync"
// )

// type Counter struct {
// 	sync.Mutex
// 	Count int
// }

// func main() {
// 	var c Counter
// 	c.Lock()
// 	defer c.Unlock()
// 	c.Count++
// 	foo(&c) // 复制锁
// }

// // 这里Counter的参数是通过复制的方式传入的
// func foo(c *Counter) {
// 	c.Lock()
// 	defer c.Unlock()
// 	fmt.Println("in foo")
// }

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 派出所证明
	var psCertificate sync.Mutex
	// 物业证明
	var propertyCertificate sync.Mutex

	var wg sync.WaitGroup
	wg.Add(2) // 需要派出所和物业都处理

	// 派出所处理goroutine
	go func() {
		defer wg.Done() // 派出所处理完成

		psCertificate.Lock()
		defer psCertificate.Unlock()

		// 检查材料
		time.Sleep(5 * time.Second)
		// 请求物业的证明
		propertyCertificate.Lock()
		defer propertyCertificate.Unlock()
	}()

	// 物业处理goroutine
	go func() {
		defer wg.Done() // 物业处理完成

		propertyCertificate.Lock()
		defer propertyCertificate.Unlock()

		// 检查材料
		time.Sleep(5 * time.Second)
		// 请求派出所的证明
		psCertificate.Lock()
		defer psCertificate.Unlock()
	}()

	wg.Wait()
	fmt.Println("成功完成")
}
