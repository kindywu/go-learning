package main

import (
	"context"
	"fmt"
	"log"
	"main/grpc/api"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type ApiServer struct {
	api.UnimplementedGreeterServer
	addr string
}

func (s *ApiServer) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloReply, error) {
	return &api.HelloReply{Message: fmt.Sprintf("%s %s (from %s)", "Hello ", in.Name, s.addr)}, nil
}

var (
	addrs       = []string{":50051", ":50052"}
	clientAddrs = []string{"localhost:50051", "localhost:50052"}
)

func main() {
	go startServer()

	go startClient()

	// Respect OS stop signals.
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-c
}

func startClient() {
	resolver.Register(&greeterResolverBuilder{})

	pickfirstConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", greeterScheme, greeterServiceName),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer pickfirstConn.Close()

	fmt.Println("--- calling helloworld.Greeter/SayHello with round_robin ---")
	makeRPCs(pickfirstConn, 10)
}

const (
	greeterScheme      = "greeter"
	greeterServiceName = "lb.api.grpc.greeter.com"
)

func callUnaryEcho(c api.GreeterClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &api.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Message)
}

func makeRPCs(cc *grpc.ClientConn, n int) {
	hwc := api.NewGreeterClient(cc)
	for i := 0; i < n; i++ {
		callUnaryEcho(hwc, "this is examples/load_balancing")
	}
}

func startServer() {
	var wg sync.WaitGroup
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			grpcServer := grpc.NewServer()
			api.RegisterGreeterServer(grpcServer, &ApiServer{addr: addr})
			log.Printf("Starting server on :50051")
			if err := grpcServer.Serve(listener); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}(addr)
	}
	wg.Wait()
}

// Following is an example name resolver implementation. Read the name
// resolution example to learn more about it.

type greeterResolverBuilder struct{}

func (*greeterResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &greeterResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			greeterServiceName: clientAddrs,
		},
	}
	r.start()
	return r, nil
}
func (*greeterResolverBuilder) Scheme() string { return greeterScheme }

type greeterResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *greeterResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*greeterResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*greeterResolver) Close()                                  {}
