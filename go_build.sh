#!/bin/bash

IMAGE_NAME="go_build"
CONTAINER_NAME="go_build_for_bbb"

if [ -z "$(docker images -q $IMAGE_NAME)" ]; then
  echo "image '$IMAGE_NAME' not found, will build it first:" && \
  docker build -t $IMAGE_NAME --target bbb-builder .
fi

docker run -v $PWD:/src --name $CONTAINER_NAME \
    -w /src \
    $IMAGE_NAME \
    go build -buildvcs=false -o app-bbb.bin . && \
docker logs $CONTAINER_NAME && \
docker rm $CONTAINER_NAME

mv app-bbb.bin ./out/