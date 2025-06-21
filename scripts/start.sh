#!/bin/bash

# CohereDB Startup Script
# This script starts the DB Manager and multiple DB Servers

set -e

echo "🚀 Starting CohereDB v1..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

if [ ! -d "data" ]; then
    print_status "Creating data directory..."
    mkdir -p data
fi

if [ ! -f "bin/db_manager" ] || [ ! -f "bin/db_server" ] || [ ! -f "bin/client" ]; then
    print_status "Binaries not found, building components..."
    make build
fi

print_status "All components built successfully!"

mkdir -p bin
mkdir -p logs
mkdir -p config

if [ ! -f "config/manager.toml" ]; then
    print_status "Creating manager config..."
    cat > config/manager.toml << 'EOFCONFIG'
[manager]
grpc_addr = "127.0.0.1:9090"
http_addr = "127.0.0.1:8090"
EOFCONFIG
fi

start_component() {
    local name=$1
    local cmd=$2
    local log_file="logs/${name}.log"
    
    print_status "Starting ${name}..."
    nohup $cmd > "$log_file" 2>&1 &
    echo $! > "logs/${name}.pid"
    sleep 2
    
    if ps -p $(cat "logs/${name}.pid") > /dev/null; then
        print_status "${name} started successfully (PID: $(cat "logs/${name}.pid"))"
    else
        print_error "Failed to start ${name}"
        print_error "Check log file: $log_file"
        return 1
    fi
}

start_component "db_manager" "./bin/db_manager -config=config/manager.toml"

sleep 3

start_component "db_server_pune" "./bin/db_server -config=config/server0.toml -register=true"
sleep 2

start_component "db_server_mumbai" "./bin/db_server -config=config/server1.toml -register=true"
sleep 2

start_component "db_server_bangalore" "./bin/db_server -config=config/server2.toml -register=true"
sleep 2

print_status "All servers started successfully!"
echo
print_status "CohereDB v1 is now running!"
echo
echo -e "${BLUE}Usage Examples:${NC}"
echo "  Set a key:    ./bin/client -op=set -key=mykey -value=myvalue"
echo "  Get a key:    ./bin/client -op=get -key=mykey"
echo "  Delete a key: ./bin/client -op=delete -key=mykey"
echo
echo -e "${BLUE}Log files:${NC}"
echo "  DB Manager:        logs/db_manager.log"
echo "  DB Server (Pune):  logs/db_server_pune.log"
echo "  DB Server (Mumbai): logs/db_server_mumbai.log"
echo "  DB Server (Bangalore): logs/db_server_bangalore.log"
echo
echo -e "${YELLOW}To stop all services, run: make stop${NC}"
