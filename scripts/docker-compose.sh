#!/bin/bash

# Default Docker Compose file path
COMPOSE_FILE="./build/docker-compose.yaml"
CMD="up --build"
# Check if the --storages-only option is passed
if [ "$1" == "--storages-only" ]; then
  COMPOSE_FILE="./build/docker-compose-storages-only.yaml"
fi

if [ "$1" == "down" ]; then
  CMD="down -v"
fi

# Run docker-compose with the selected YAML file
docker-compose -f "$COMPOSE_FILE" $CMD