#!/usr/bin/env bash

set -o errexit

docker push "chuckdha/kubeyaml-web:latest"
docker push "chuckdha/kubeyaml-backend:latest"
