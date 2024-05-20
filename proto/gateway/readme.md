go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway

# generate code
protoc -I ./proto --go_out=. --go-grpc_out=. --grpc-gateway_out=. ./proto/service.proto

# post
Invoke-WebRequest -Method POST `
                  -Body (@{"name"="kindywu";}|ConvertTo-Json) `
                  -Uri http://localhost:50052/v1/greeter/sayhello `
                  -ContentType application/json