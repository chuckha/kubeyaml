#!/usr/bin/env bash

set -o errexit

IMAGE_NAME=${IMAGE_NAME:-chuckdha/kubeyaml}
IMAGE_TAG=${IMAGE_TAG:-latest}

docker push "${IMAGE_NAME}:${IMAGE_TAG}"
