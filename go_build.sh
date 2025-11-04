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
    -o type=docker,dest=image.tar \
    -t "$FULL_TAG" .
    # --push \
    # -t "$DOCKER_HUB/$IMAGE_NAME:latest" .
  
  echo "[✔] Image built successfully. Exported to image.tar"
  
  echo "[🔍] Running Trivy scan before push..."
  
  docker run --rm -v $PWD:/scan aquasec/trivy:latest image \
  --exit-code 1 \
  --severity HIGH,CRITICAL \
  --input /scan/image.tar \
  --format table --output /scan/report.txt

  echo "[✔] Trivy scan passed. Pushing image..."
  docker load -i image.tar
  docker push "$FULL_TAG"
  docker tag "$FULL_TAG" "$DOCKER_HUB/$IMAGE_NAME:latest"
  docker push "$DOCKER_HUB/$IMAGE_NAME:latest"

  echo "[🧹] Cleaning up tarball..."
  rm -f "$TAR_FILE"

  echo "[✔] Image pushed: $FULL_TAG"
  echo "[✔] Also tagged as :latest"

}

case "${1:-}" in
  bin) build_bin ;;
  image) build_image ;;
  *) echo "Usage: $0 {bin|image}"; exit 1 ;;
esac