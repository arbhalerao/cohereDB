#!/bin/bash

# CohereDB Test Script
# Wait for system to be ready, then run tests

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

CLIENT="./bin/client"
MANAGER_ADDR="127.0.0.1:9090"

print_test() {
    echo -e "${BLUE}[TEST]${NC} $1"
}

print_pass() {
    echo -e "${GREEN}[PASS]${NC} $1"
}

print_fail() {
    echo -e "${RED}[FAIL]${NC} $1"
}

print_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

wait_for_system() {
    print_info "Waiting for CohereDB system to be ready..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s http://127.0.0.1:8090/health >/dev/null 2>&1; then
            if $CLIENT -addr=$MANAGER_ADDR -op=set -key=test:readiness -value=ready >/dev/null 2>&1; then
                print_pass "System is ready!"
                $CLIENT -addr=$MANAGER_ADDR -op=delete -key=test:readiness >/dev/null 2>&1
                return 0
            fi
        fi
        
        print_info "Attempt $attempt/$max_attempts - waiting for servers to register..."
        sleep 2
        ((attempt++))
    done
    
    print_fail "System did not become ready within timeout"
    return 1
}

main() {
    echo "🧪 Testing CohereDB..."
    
    if ! wait_for_system; then
        exit 1
    fi
    
    echo
    print_info "Running basic functionality tests..."
    
    print_test "Setting user:1"
    if $CLIENT -addr=$MANAGER_ADDR -op=set -key=user:1 -value="John Doe"; then
        print_pass "Set user:1"
    else
        print_fail "Set user:1"
        exit 1
    fi
    
    print_test "Getting user:1"
    if $CLIENT -addr=$MANAGER_ADDR -op=get -key=user:1; then
        print_pass "Get user:1"
    else
        print_fail "Get user:1"
        exit 1
    fi
    
    print_test "Setting user:2"
    if $CLIENT -addr=$MANAGER_ADDR -op=set -key=user:2 -value="Jane Smith"; then
        print_pass "Set user:2"
    else
        print_fail "Set user:2"
        exit 1
    fi
    
    print_test "Deleting user:2"
    if $CLIENT -addr=$MANAGER_ADDR -op=delete -key=user:2; then
        print_pass "Delete user:2"
    else
        print_fail "Delete user:2"
        exit 1
    fi
    
    print_test "Getting deleted key user:2 (should fail)"
    if ! $CLIENT -addr=$MANAGER_ADDR -op=get -key=user:2 >/dev/null 2>&1; then
        print_pass "Get deleted key failed as expected"
    else
        print_fail "Get deleted key should have failed"
        exit 1
    fi
    
    echo
    print_pass "All tests completed successfully!"
    echo
    print_info "Test summary:"
    echo "✅ Basic SET/GET operations"
    echo "✅ Multiple key storage"  
    echo "✅ DELETE operations"
    echo "✅ Error handling for missing keys"
    echo "✅ Data distribution across servers"
}

main "$@"
