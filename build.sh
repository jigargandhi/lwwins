#!/bin/sh

# building cli project
cd services/
go generate
cd .. 
go test ./...
go build -o dist/lwwinscli cmd/cli/main.go
go build -o dist/lwwinsruntime cmd/runtime/main.go
docker build -t lwwins:dev -f docker/Dockerfile .

# TODO: add build for main project
# 

