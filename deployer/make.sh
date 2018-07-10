#!/bin/bash

source functions.sh

# TODO: Make customizable values actually customaizable

# Customizable values
#PROJECT=elastifile-gce-lab-c304
PROJECT=launcher-poc-207208
DOCKER_FILE=Dockerfile
TAG=latest
REGISTRY=docker.io/elastifileio

if [ -n "$1" ]; then
    PROJECT="$1"
fi

if [ -n "$2" ]; then
    DOCKER_FILE="$2"
fi

# Static values
IMAGE=elastifile-provisioner-deployer
STAGING_CONTAINER_REGISTRY=gcr.io/$PROJECT
IMAGE_TAGGED=$STAGING_CONTAINER_REGISTRY/$IMAGE:$TAG
logme "Building "$DOCKER_FILE" and tagging the image as $IMAGE_TAGGED"
exec_cmd docker build -t $IMAGE_TAGGED --build-arg REGISTRY=$REGISTRY --build-arg TAG=$TAG . -f $DOCKER_FILE
logme Pushing $IMAGE_TAGGED
exec_cmd docker push $IMAGE_TAGGED
logme Listing available images
exec_cmd gcloud container images list --repository=$STAGING_CONTAINER_REGISTRY

echo In order to free up the disk space:
echo gcloud container images delete $IMAGE_TAGGED

