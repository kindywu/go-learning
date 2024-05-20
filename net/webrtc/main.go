package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pion/webrtc/v3"
)

// 定义一个全局变量来存储offer和answer
var offer *webrtc.SessionDescription

func createPeerConnection(config webrtc.Configuration) (*webrtc.PeerConnection, error) {
	api := webrtc.NewAPI()
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}
	return peerConnection, nil
}

func createOffer(peerConnection *webrtc.PeerConnection) (*webrtc.SessionDescription, error) {
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		return nil, err
	}
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

func setAnswer(peerConnection *webrtc.PeerConnection, answer *webrtc.SessionDescription) error {
	err := peerConnection.SetRemoteDescription(*answer)
	if err != nil {
		return err
	}
	_, err = peerConnection.CreateAnswer(nil)
	if err != nil {
		return err
	}
	return nil
}

func sender() {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}
	peerConnection, err := createPeerConnection(config)
	if err != nil {
		log.Fatal(err)
	}
	defer peerConnection.Close()

	offer, err = createOffer(peerConnection)
	if err != nil {
		log.Fatal(err)
	}

	// 创建数据通道
	dataChannel, err := peerConnection.CreateDataChannel("gameChannel", nil)
	if err != nil {
		log.Fatal(err)
	}

	isOpen := false
	dataChannel.OnOpen(func() {
		isOpen = true
		fmt.Println("Data channel is open!")
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Received message: %s\n", string(msg.Data))
	})

	// 模拟发送游戏状态更新
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		if isOpen {
			gameState := "Player position: (100, 100)"
			err := dataChannel.SendText(gameState)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func receiver() {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}
	peerConnection, err := createPeerConnection(config)
	if err != nil {
		log.Fatal(err)
	}
	defer peerConnection.Close()

	// 等待接收offer
	time.Sleep(1 * time.Second) // 确保发送者已经创建了offer
	if offer == nil {
		log.Fatal("No offer received")
	}
	if err := setAnswer(peerConnection, offer); err != nil {
		log.Fatal(err)
	}

	// 创建数据通道
	dataChannel, err := peerConnection.CreateDataChannel("gameChannel", nil)
	if err != nil {
		log.Fatal(err)
	}

	dataChannel.OnOpen(func() {
		fmt.Println("Data channel is open!")
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Received message: %s\n", string(msg.Data))
	})
}

func main() {
	// 启动发送者和接收者
	go sender()
	time.Sleep(3 * time.Second)
	go receiver()

	// 等待接收者准备好
	time.Sleep(3 * time.Second)
	// 这里可以添加其他逻辑，例如设置远程描述等
}
