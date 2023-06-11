#!/bin/bash

IMAGE_NAME="activity-app:v1.01"
CONTAINER_NAME="activity-api-server"

# Stop and remove the existing container if it exists
if docker ps -a --format '{{.Names}}' | grep -Eq "^${CONTAINER_NAME}\$"; then
  docker stop $CONTAINER_NAME
  docker rm $CONTAINER_NAME
fi
