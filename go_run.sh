#!/bin/bash

IMAGE_NAME="go_run_env"
CONTAINER_NAME="go_app"
PORT=8080

if [ -z "$(docker images -q $IMAGE_NAME)" ]; then
  echo "image '$IMAGE_NAME' not found, will build it first:" && \
  docker build -f Dockerfile.build -t $IMAGE_NAME --target dev .
fi

if [ "$1" = "stop" ]; then
  if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}\$"; then
    docker rm -f "$CONTAINER_NAME" >/dev/null 2>&1
    echo "container '$CONTAINER_NAME' stopped and removed"
  else
    echo "container not found"
  fi
  exit 0
fi

CID=$(docker run -d \
  --name $CONTAINER_NAME \
  -v "$(pwd)":/usr/src/app \
  -p $PORT:8080 \
  -w /usr/src/app \
  $IMAGE_NAME)

if [ -n "$CID" ]; then
  echo "App have been execute (container ID: $CID)"
  echo "You can access it at http://localhost:8080"
else
  echo "Failed to start app"
fi