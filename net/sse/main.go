package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var m = sync.Map{}

func sseHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置必要的头部
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		user_id := r.Header.Get("user_id") //not empty and unique

		var data_ready chan string
		if value, ok := m.Load(user_id); ok {
			data_ready = value.(chan string)
		} else {
			data_ready = make(chan string, 1000)
			m.Store(user_id, data_ready)
		}

		// 创建一个关闭的通道用于在客户端断开连接时退出循环
		done := make(chan struct{})

		// 启动一个协程来监听请求的关闭事件
		go func() {
			<-r.Context().Done()
			close(done)
		}()

		for {
			select {
			case data := <-data_ready:
				_, err := fmt.Print(w, data)
				if err != nil {
					log.Println("Error writing to response:", err)
					continue
				}
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			case <-done:
				{
					m.Delete(user_id)
					close(data_ready)
				}
				return
			}
		}
	}
}

func main() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			currentTime := time.Now().Format(time.RFC3339)
			m.Range(func(key, _ interface{}) bool {
				send_to_user(key.(string), currentTime)
				return true // 继续迭代直到结束
			})

		}
	}()

	http.HandleFunc("/stream", sseHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func send_to_user(user_id string, data string) {
	if value, ok := m.Load(user_id); ok {
		data_ready := value.(chan string)
		data_ready <- data
	}
}

// const eventSource = new EventSource('/stream');
// eventSource.onmessage = function(event) {
//     console.log('Received time:', event.data);
// };
