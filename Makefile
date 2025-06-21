DB_SERVER_PROTO=./proto/db_server.proto
DB_MANAGER_PROTO=./proto/db_manager.proto

COHERE_DB_DIR=.
DB_SERVER_DIR=$(COHERE_DB_DIR)/pb/db_server
DB_MANAGER_DIR=$(COHERE_DB_DIR)/pb/db_manager

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
	rm -rf bin/
	rm -rf data/
	rm -rf logs/

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

build: generate
	@echo "Building all CohereDB components..."
	@mkdir -p bin
	@echo "Building DB Manager..."
	@cd db_manager && go build -o ../bin/db_manager ./cmd/main.go
	@echo "Building DB Server..."
	@cd db_server && go build -o ../bin/db_server ./cmd/main.go
	@echo "Building Client..."
	@cd client && go build -o ../bin/client ./main.go
	@echo "All components built successfully!"

start: build
	@echo "Starting CohereDB v1..."
	@chmod +x scripts/start.sh
	@./scripts/start.sh

stop:
	@echo "Stopping CohereDB v1..."
	@chmod +x scripts/stop.sh
	@./scripts/stop.sh

test: build start
	@echo "Testing CohereDB..."
	@./scripts/test.sh

setup:
	@echo "Setting up CohereDB development environment..."
	@mkdir -p bin logs data scripts config
	@chmod +x scripts/*.sh || true
	@echo "Development environment ready!"

generate: generate-db-server generate-db-manager

all: generate build

.PHONY: generate-db-server generate-db-manager clean install lint build start stop test setup generate all
