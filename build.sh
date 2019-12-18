#!/bin/sh

# building cli project
go test ./...
go build -o dist/lwwinscli cmd/cli/main.go

# TODO: add build for main project
# docker build -t lwwins:dev -f docker/Dockerfile 

