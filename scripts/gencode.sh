#!/bin/bash

# generate swagger
rm -rf server/web/gen/*
swagger generate server -f swagger.yml --principal=models.Principle -t server/web/gen --regenerate-configureapi

protoc \
  --go-grpc_out=./pb \
  --go_out=./pb \
  --proto_path=. online.proto
