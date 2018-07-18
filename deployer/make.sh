#!/bin/bash

source functions.sh

if [ -n "$1" ]; then
    PROJECT="$1"
else
    PROJECT=$(gcloud config get-value project)
fi

if [ -n "$2" ]; then
    TAG="$2"
else
    TAG=1.0
fi

# Customizable values
REGISTRY=gcr.io/${PROJECT}

# Static values
APP_IMAGE_NAME=provisioner
DEPLOYER_IMAGE_TAGGED=${REGISTRY}/${APP_IMAGE_NAME}/deployer:${TAG}
IMAGE_CLONE_TAGGED=${REGISTRY}/${APP_IMAGE_NAME}/custom_script:${TAG} # Fake image used to run deploy.sh

logme "Building and tagging the deployer image as ${DEPLOYER_IMAGE_TAGGED}"
exec_cmd docker build -t ${DEPLOYER_IMAGE_TAGGED} --build-arg REGISTRY=${REGISTRY} --build-arg TAG=${TAG} .
logme "Pushing ${DEPLOYER_IMAGE_TAGGED}"
exec_cmd docker push ${DEPLOYER_IMAGE_TAGGED}

logme "Tagging and pushing ${IMAGE_CLONE_TAGGED}"
exec_cmd docker tag ${DEPLOYER_IMAGE_TAGGED} ${IMAGE_CLONE_TAGGED}
exec_cmd docker push ${IMAGE_CLONE_TAGGED}

logme "Listing available images"
exec_cmd gcloud container images list --repository=${REGISTRY}/${APP_IMAGE_NAME}

echo "In order to free up the disk space:"
echo "gcloud container images delete ${DEPLOYER_IMAGE_TAGGED}"
