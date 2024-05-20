package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"main/proto/gateway/generated"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcServerAddr = flag.String("grpc-server", ":50051", "gRPC server address")
	gatewayAddr    = flag.String("gateway", ":50052", "Gateway server address")
)

type server struct {
	generated.UnimplementedGreeterServer
}

// SayHello implements generated.GreeterServer
func (s *server) SayHello(ctx context.Context, req *generated.HelloRequest) (*generated.HelloReply, error) {
	log.Printf("Server Received: %v", req.GetName())
	return &generated.HelloReply{Message: "Hello " + req.GetName()}, nil
}

func start_gRPC_Server() {
	lis, err := net.Listen("tcp", *grpcServerAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Service listening on %s\n", *grpcServerAddr)

	s := grpc.NewServer()
	generated.RegisterGreeterServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func start_gRPC_Gateway() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock()}

	err := generated.RegisterGreeterHandlerFromEndpoint(ctx, mux, "localhost"+*grpcServerAddr, opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Gateway listening on %s\n", *gatewayAddr)
	if err := http.ListenAndServe(*gatewayAddr, mux); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

func main() {
	flag.Parse()

	go start_gRPC_Server()
	go start_gRPC_Gateway()

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
