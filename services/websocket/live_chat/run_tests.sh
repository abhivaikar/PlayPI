#!/bin/bash

set -e

# Function to kill background processes
cleanup() {
    echo "Cleaning up..."
    kill $WS_SERVER_PID $CLIENT_A_PID $CLIENT_B_PID
}
trap cleanup EXIT

echo "Running automated tests for websocket live-chat server..."
# Step 1: Run unit tests
echo "Running unit tests..."
go test . -v -count=1

# Step 2: Build and start the servers
echo "Building and starting the WebSocket server..."
go build -o websocket_server ./main/main.go
./websocket_server > websocket_server.log 2>&1 &
WS_SERVER_PID=$!

echo "Building and starting Client A..."
go build -o client_a ./e2e_tests/client/client.go
./client_a 8090 > client_a.log 2>&1 &
CLIENT_A_PID=$!

echo "Building and starting Client B..."
go build -o client_b ./e2e_tests/client/client.go
./client_b 8091 > client_b.log 2>&1 &
CLIENT_B_PID=$!

# Wait for servers to start
sleep 5

# Step 3: Run the e2e tests
echo "Running e2e tests..."
go test ./e2e_tests -v -count=1

# Step 4: Gracefully shutdown the servers
echo "Tests completed. Shutting down servers..."
cleanup