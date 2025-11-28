#!/bin/bash

set -e 

echo "Building ProxyX ..."

GOOS=linux GOARCH=amd64 go build -o bin/proxy cmd/proxy/main.go

echo "Done! -> bin/proxy"

