protoc --go_out=./protogen \
    --go-grpc_out=./protogen \
    ./protoc/brand/*.proto

protoc --go_out=./protogen \
    --go-grpc_out=./protogen \
    ./protoc/car/*.proto