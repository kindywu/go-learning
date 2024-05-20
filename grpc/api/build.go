package api

// -I../build/grpc-gateway-v2.3.0/third_party/googleapis

//go:generate protoc -I. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative helloworld.proto

//go:generate protoc -I. --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative helloworld.proto
