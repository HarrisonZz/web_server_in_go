#!/bin/bash
set -euo pipefail

IMAGE_NAME="go_build"
CONTAINER_NAME="go_build_for_bbb"
DOCKER_HUB="harrisonchen0418"
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "nogit")

function build_bin() {
  echo "[*] Building ARMv7 binary only..."
  if ! docker image inspect "$IMAGE_NAME" >/dev/null 2>&1; then
    echo "Image '$IMAGE_NAME' not found, building..."
    docker build -f Dockerfile.build -t "$IMAGE_NAME" --target bbb-builder .
  fi

  docker run --rm \
    -v "$PWD:/src" \
    -w /src \
    --name "$CONTAINER_NAME" \
    "$IMAGE_NAME" \
    bash -c '
      echo "[*] Running go build..."
      go mod tidy
      GOOS=linux GOARCH=arm GOARM=7 go build -buildvcs=false -o out/app-bbb.bin .
    '

  echo "[✔] Binary built successfully: ./out/app-bbb.bin"
}

function build_image() {
  local FULL_TAG="$DOCKER_HUB/$IMAGE_NAME:git-$GIT_COMMIT"
  echo "[*] Building full BBB runtime image: $FULL_TAG"
  docker buildx build \
    --platform linux/arm/v7 \
    -f Dockerfile.runtime \
    --push \
    -t "$FULL_TAG" \
    -t "$DOCKER_HUB/$IMAGE_NAME:latest" .

  echo "[✔] Image pushed: $FULL_TAG"
  echo "[✔] Also tagged as :latest"

}

case "${1:-}" in
  bin) build_bin ;;
  image) build_image ;;
  *) echo "Usage: $0 {bin|image}"; exit 1 ;;
esac