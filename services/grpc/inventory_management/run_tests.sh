#!/bin/bash

set -e

echo "Running automated tests for grPC inventory management API..."

# Step 1: Run unit tests
echo "Running unit tests..."
go test . -v -count=1

echo "Completed successfully!"