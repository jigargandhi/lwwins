FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app; protobuf compilation is done outside
RUN go build -o dist/lwwinsruntime cmd/runtime/main.go

EXPOSE 3334

ENTRYPOINT ["./dist/lwwinsruntime"]
