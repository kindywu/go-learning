package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有CORS请求
		return true
	},
}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				fmt.Println("Client has closed the WebSocket connection.")
			} else {
				fmt.Println("Error:", err)
			}
			break
		}
		fmt.Printf("\nServer received message: %s\n", message)

		// Echo the message back to the client
		err = conn.WriteMessage(mt, message)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
	}
}

func server() {
	http.HandleFunc("/echo", echo)
	log.Println("WebSocket server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func client() {
	url := "ws://localhost:8080/echo"
	fmt.Printf("connecting to %s\n", url)

	// 创建一个WebSocket连接
	var ws *websocket.Conn
	var err error
	ws, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// 捕获中断信号，优雅地关闭连接
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go func() {
		for sig := range c {
			fmt.Println("Captured signal:", sig)
			ws.WriteMessage(websocket.CloseMessage, nil)
			os.Exit(0)
		}
	}()

	// 读取消息
	go func() {
		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Client received message: %s type:%d .\n", message, mt)
		}
	}()

	// 发送消息
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Client send message: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		message := []byte(text[:len(text)-1]) // 去掉换行符

		err = ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	go func() {
		time.Sleep(3 * time.Second)
		client()
	}()
	server()
}
