DB_SERVER_PROTO=./proto/db_server.proto
DB_MANAGER_PROTO=./proto/db_manager.proto

COHERE_DB_DIR=.
DB_SERVER_DIR=$(COHERE_DB_DIR)/pb/db_server
DB_MANAGER_DIR=$(COHERE_DB_DIR)/pb/db_manager

# Generate the db_server code
generate-db-server:
	protoc --proto_path=./proto \
		--go_out=$(DB_SERVER_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(DB_SERVER_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(DB_SERVER_PROTO)

generate-db-manager:
	protoc --proto_path=./proto \
		--go_out=$(DB_MANAGER_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(DB_MANAGER_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(DB_MANAGER_PROTO)

clean:
	rm -f $(DB_SERVER_DIR)/*.pb.go
	rm -f $(DB_MANAGER_DIR)/*.pb.go

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
