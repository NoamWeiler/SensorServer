#!/bin/bash

echo "Running 3 server instances"
go run cmd/server.go -port=50051 &
go run cmd/server.go -port=50052 &
go run cmd/server.go -port=50053 &
