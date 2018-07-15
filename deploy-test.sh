#!/bin/bash

set -e

pwd

docker build -t alvaradojl/${PROJECT_NAME}:$TRAVIS_COMMIT .

docker tag alvaradojl/${PROJECT_NAME}:$TRAVIS_COMMIT alvaradojl/${PROJECT_NAME}:latest

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

docker push alvaradojl/${PROJECT_NAME}