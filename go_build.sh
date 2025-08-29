#!/bin/bash

IMAGE_NAME="go_build"
CONTAINER_NAME="go_build_for_bbb"
DOCKER_HUB="harrisonchen0418"


case "$1" in
  bin)
    echo "[*] Building ARMv7 binary only..."
    if [ -z "$(docker images -q $IMAGE_NAME)" ]; then
      echo "image '$IMAGE_NAME' not found, will build it first:" && \
      docker build -f Dockerfile.build -t $IMAGE_NAME --target bbb-builder .
    fi
    docker run --rm -v $PWD:/src --name $CONTAINER_NAME \
    -w /src $IMAGE_NAME go build -buildvcs=false -o app-bbb.bin . && \
    mv app-bbb.bin ./out/
    ;;
  image)
    echo "Building full BBB runtime image..."
    docker buildx build --platform linux/arm/v7 -f Dockerfile.runtime --push -t $DOCKER_HUB/$IMAGE_NAME:BBB .
    ;;
  *)
    echo "Usage: $0 {bin|image}"
    exit 1
    ;;
esac