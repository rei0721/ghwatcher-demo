#!/bin/bash
set -e

# Configuration
IMAGE_NAME="rei0721-server"
CONTAINER_NAME="rei0721-server"
PORT_MAPPING="9900-9999:9900-9999"
VOLUME_MAPPING="/root/.rei0721:/app"


# Check if running locally with source or remotely
if [ -f "Dockerfile" ]; then
    BUILD_CONTEXT="."
    echo ">>> Found Dockerfile locally. Building from local source..."
else
    BUILD_CONTEXT="https://github.com/rei0721/ghhook-server.git"
    echo ">>> Dockerfile not found. Building from remote repository: $BUILD_CONTEXT..."
fi

echo ">>> Building Docker image: $IMAGE_NAME..."
docker build -t "$IMAGE_NAME" "$BUILD_CONTEXT"

echo ">>> Checking for existing container '$CONTAINER_NAME'..."
if [ "$(docker ps -aq -f name=^/${CONTAINER_NAME}$)" ]; then
    echo ">>> Stopping and removing existing container..."
    docker rm -f "$CONTAINER_NAME"
fi

echo ">>> Starting container '$CONTAINER_NAME'..."
docker run -d \
  --restart unless-stopped \
  --name "$CONTAINER_NAME" \
  -p "$PORT_MAPPING" \
  -v "$VOLUME_MAPPING" \
  "$IMAGE_NAME"

echo ">>> Deployment completed successfully."
docker ps -f name="$CONTAINER_NAME"
