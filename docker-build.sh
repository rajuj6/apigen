#!/bin/bash
set -e
if [ -z "$1" ]
then
      echo "please pass target as 1st argument"
      exit
fi

DOCKER_BUILDKIT=1 docker buildx build --platform linux/amd64,linux/arm64  . -t ghcr.io/middleware-labs/loadgen-api:$1 --push
