download https://github.com/protocolbuffers/protobuf/releases

go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

# generate code
protoc --go_out=. --go-grpc_out=. ./proto/message.proto