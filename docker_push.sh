#!/bin/bash

set -e  # Exit if any command fails

# Define variables
IMAGE_NAME="abhijeetvaikar/playpi"
TAG="latest"

echo "Building Docker image..."
docker build -t $IMAGE_NAME:$TAG .

echo "Pushing Docker image to Docker Hub..."
docker push $IMAGE_NAME:$TAG

echo "Docker image pushed successfully!"