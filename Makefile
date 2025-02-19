PROTO_FILE=./proto/db_server.proto

SERVER_DIR=./db_server

generate:
	protoc --proto_path=./proto \
		--go_out=$(SERVER_DIR)/grpc/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(SERVER_DIR)/grpc/pb \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)

clean:
	rm -f $(SERVER_DIR)/grpc/db_server.pb.go
	rm -f $(SERVER_DIR)/grpc/db_server_grpc.pb.go

install:
	go get google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

all: generate
