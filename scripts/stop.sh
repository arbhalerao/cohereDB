#!/bin/bash

echo "🛑 Stopping CohereDB v1..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

stop_component() {
    local name=$1
    local pid_file="logs/${name}.pid"
    
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if ps -p $pid > /dev/null; then
            print_status "Stopping ${name} (PID: $pid)..."
            kill $pid
            sleep 2
            
            if ps -p $pid > /dev/null; then
                print_warning "Process $pid still running, force killing..."
                kill -9 $pid
            fi
            
            print_status "${name} stopped successfully"
        else
            print_warning "${name} was not running"
        fi
        rm -f "$pid_file"
    else
        print_warning "No PID file found for ${name}"
    fi
}

stop_component "db_server_bangalore"
stop_component "db_server_mumbai"
stop_component "db_server_pune"
stop_component "db_manager"

print_status "Cleaning up any remaining processes..."
pkill -f "db_manager" || true
pkill -f "db_server" || true

print_status "All CohereDB components stopped successfully!"
print_status "CohereDB v1 has been stopped."
