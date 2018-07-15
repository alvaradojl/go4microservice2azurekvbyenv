#!/bin/bash

set -e

docker build -t alvaradojl/${PROJECT_NAME}:$TRAVIS_COMMIT ./cmd/keyvault

docker tag alvaradojl/${PROJECT_NAME}:$TRAVIS_COMMIT alvaradojl/${PROJECT_NAME}:latest

docker push alvaradojl/${PROJECT_NAME}