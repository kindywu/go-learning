package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"main/proto/cs/generated"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// server is used to implement generated.GreeterServer.
type server struct {
	generated.UnimplementedGreeterServer
}

// SayHello implements generated.GreeterServer
func (s *server) SayHello(ctx context.Context, req *generated.HelloRequest) (*generated.HelloReply, error) {
	log.Printf("Server Received: %v", req.GetName())
	return &generated.HelloReply{Message: "Hello " + req.GetName()}, nil
}

func startServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on :50051")

	s := grpc.NewServer()
	generated.RegisterGreeterServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

const (
	address = "localhost:50051"
)

func startClient() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	log.Printf("connected to %s", address)

	c := generated.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &generated.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Client: %s", r.Message)
}

func main() {
	go startServer()
	time.Sleep(3 * time.Second)
	go startClient()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			fmt.Printf("Captured %v, shutting down.\n", sig)
			break
		}
	}()

	// 等待用户输入或捕获到信号
	fmt.Println("Press Enter to exit...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}
