#!/bin/bash

set -e

# build docker image and tag it to commit identifier
docker build -t alvaradojl/${PROJECT_NAME}:$TRAVIS_COMMIT .

# create a 'latest' tag for docker image
docker tag alvaradojl/${PROJECT_NAME}:$TRAVIS_COMMIT alvaradojl/${PROJECT_NAME}:latest

# login to docker cli
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

# push docker image to registry
docker push alvaradojl/${PROJECT_NAME}

# login to azure cli
az login --service-principal --username "$AZURE_USERNAME" --password "$AZURE_PASSWORD" --tenant "$AZURE_TENANT"

# install kubectl from azure
sudo az aks install-cli

# embed credentials before using kubectl
az aks get-credentials --resource-group "$AZURE_RESOURCEGROUPNAME" --name "$AZURE_KUBERNETESCLUSTERNAME"

# create kubernetes deployment
kubectl apply -f ./k8s-deployment.yaml
# create kubernetes service
kubectl apply -f ./k8s-service.yaml