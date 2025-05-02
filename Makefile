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

lint:
	@echo "Running lint on db..."
	@(cd db && golangci-lint run)

	@echo "Running lint on db_manager..."
	@(cd db_manager && golangci-lint run)

	@echo "Running lint on db_server..."
	@(cd db_server && golangci-lint run)

	@echo "Running lint on utils..."
	@(cd utils && golangci-lint run)

	@echo "Linting completed!"


all: generate
