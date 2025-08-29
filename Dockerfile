FROM golang:1.24 AS base

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

# ----------------------------------------------

FROM base AS dev
CMD ["go", "run", "."]

# ----------------------------------------------

FROM base AS bbb-builder
COPY . .
ENV GOOS=linux GOARCH=arm GOARM=7
RUN go build -buildvcs=false -o /out/app-bbb.bin .