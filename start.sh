#!/bin/bash

IMAGE_NAME="activity-app:v1.01"
CONTAINER_NAME="activity-api-server"

# Stop and remove the existing container if it exists
if docker ps -a --format '{{.Names}}' | grep -Eq "^${CONTAINER_NAME}\$"; then
  docker stop $CONTAINER_NAME
  docker rm $CONTAINER_NAME
fi

# Build the Docker image
docker build -t $IMAGE_NAME .

# Run the Docker container and publish port
docker run -d -p 9000:9000 --name $CONTAINER_NAME $IMAGE_NAME

# Output container logs for debugging
docker logs $CONTAINER_NAME
