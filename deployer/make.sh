#!/bin/bash

source functions.sh

# TODO: Make customizable values actually customaizable

# Customizable values
#PROJECT=elastifile-gce-lab-c304
PROJECT=launcher-poc-207208
TAG=latest
DOCKER_FILE=Dockerfile
REGISTRY=docker.io/elastifileio

if [ -n "$1" ]; then
    PROJECT="$1"
fi

if [ -n "$2" ]; then
    DOCKER_FILE="$2"
fi

# Static values
IMAGE=elastifile-provisioner-deployer

logme "Building "$DOCKER_FILE" and tagging the image as gcr.io/$PROJECT/$IMAGE:$TAG"
IMAGE_TAGGED=gcr.io/$PROJECT/$IMAGE:$TAG
exec_cmd docker build -t $IMAGE_TAGGED --build-arg REGISTRY=$REGISTRY --build-arg TAG=$TAG . -f $DOCKER_FILE
logme Pushing $IMAGE_TAGGED
exec_cmd docker push $IMAGE_TAGGED
logme Listing available images
exec_cmd gcloud container images list --repository=gcr.io/$PROJECT

echo In order to free up the disk space:
echo gcloud container images delete $IMAGE_TAGGED

