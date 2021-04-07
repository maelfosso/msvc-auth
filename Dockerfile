FROM golang:1.16.2 AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# Move to working directory /build
WORKDIR /app

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 5000

# ENTRYPOINT CompileDaemon --command="go run models.go handlers.go main.go"
ENTRYPOINT [ "./bin/entry.sh" ]

# ENTRYPOINT CompileDaemon --build="go build commands/runserver.go" --command=./runserver
