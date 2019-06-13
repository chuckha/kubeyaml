#!/usr/bin/env bash

set -o errexit

go mod vendor
go test ./...
go run ./cmd/generate/main.go
npm run build-prod

docker build . -t "chuckdha/kubeyaml-web:latest"
docker build --file Dockerfile.backend . -t "chuckdha/kubeyaml-backend:latest"
