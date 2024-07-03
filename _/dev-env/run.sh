#!/bin/sh
set -e

# Env. variables
DOCKER_IMG_NAME="sap-grant-dev-env-img:latest"
# path to the dir that will be mounted @ /var/lib/sap-grant-dev-root/ inside the container
PROJECT_DIR="/Users/bojinovic/Documents/MVPw-Projects/3327/sap-grant"

# ---------------------------------------

# Image build
docker buildx build --platform=linux/amd64 -t ${DOCKER_IMG_NAME} .

# Dev. Env. setup
export "DOCKER_IMG_NAME"=${DOCKER_IMG_NAME} "PROJECT_DIR"=${PROJECT_DIR} && docker compose up