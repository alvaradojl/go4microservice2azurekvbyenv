#!/bin/bash

set -e

pwd

docker build -t alvaradojl/${PROJECT_NAME}:$TRAVIS_COMMIT .

docker tag alvaradojl/${PROJECT_NAME}:$TRAVIS_COMMIT alvaradojl/${PROJECT_NAME}:latest

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

docker push alvaradojl/${PROJECT_NAME}

az login --service-principal --username "$AZURE_USERNAME" --password "$AZURE_PASSWORD" --tenant "$AZURE_TENANT"

az aks get-credentials --resource-group RelayServices --name RelayServicesTest

kubectl apply -f ./kubernetes.yaml