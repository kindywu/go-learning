package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "main/grpc/api"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ApiServer struct {
	api.UnimplementedGreeterServer
}

func (s *ApiServer) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloReply, error) {
	return &api.HelloReply{Message: "Hello " + in.Name}, nil
}

var grpc_port = 50051

func grpc_client() {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", grpc_port), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &api.HelloRequest{Name: "Kimi"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}

func grpc_server() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpc_port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterGreeterServer(grpcServer, &ApiServer{})
	log.Printf("Starting server on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func gateway_server() {
	// 启动 gRPC-Gateway
	httpPort := 9999
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	grpcClient, err := grpc.DialContext(ctx, fmt.Sprintf("localhost:%d", grpc_port), opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	api.RegisterGreeterHandler(ctx, mux, grpcClient)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	})
	log.Printf("Starting HTTP server on port %d", httpPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), handler); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	go grpc_server()

	go gateway_server()

	time.AfterFunc(3*time.Second, grpc_client)

	// Respect OS stop signals.
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-c
}
