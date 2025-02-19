DB_SERVER_PROTO=./proto/db_server.proto
DB_MANAGER_PROTO=./proto/db_manager.proto

DB_SERVER_DIR=./db_server
DB_MANAGER_DIR=./db_manager

generate-db-server:
	protoc --proto_path=./proto \
		--go_out=$(DB_SERVER_DIR)/grpc/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(DB_SERVER_DIR)/grpc/pb \
		--go-grpc_opt=paths=source_relative \
		$(DB_SERVER_PROTO)

generate-db-manager:
	protoc --proto_path=./proto \
		--go_out=$(DB_MANAGER_DIR)/grpc/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(DB_MANAGER_DIR)/grpc/pb \
		--go-grpc_opt=paths=source_relative \
		$(DB_MANAGER_PROTO)

clean:
	rm -f $(DB_SERVER_DIR)/grpc/pb/*.pb.go
	rm -f $(DB_MANAGER_DIR)/grpc/pb/*.pb.go

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

# Run all proto generations
generate: generate-db-server generate-db-manager

all: generate
