#!/bin/bash

set -e

pwd

docker build -t alvaradojl/${PROJECT_NAME}:$TRAVIS_COMMIT .

docker tag alvaradojl/${PROJECT_NAME}:$TRAVIS_COMMIT alvaradojl/${PROJECT_NAME}:latest

docker push alvaradojl/${PROJECT_NAME}