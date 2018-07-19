#!/bin/bash

MYPATH=$(dirname $0)
source $MYPATH/../functions.sh

# TODO: Make customizable values actually customaizable

# Customizable values
if [ -n "$1" ]; then
    PROJECT="$1"
else
    PROJECT=$(gcloud config get-value project)
fi

if [ -n "$2" ]; then
    TAG="$2"
else
    TAG=latest
fi

DOCKER_FILE=Dockerfile

# Static values
IMAGE=nfs-test

DOCKER_REPO=gcr.io/${PROJECT}
IMAGE_TAGGED=${DOCKER_REPO}/$IMAGE:$TAG
logme "Building "$DOCKER_FILE" and tagging the image as $IMAGE_TAGGED"
exec_cmd docker build -t $IMAGE_TAGGED . -f $DOCKER_FILE
logme Pushing $IMAGE_TAGGED
exec_cmd docker push $IMAGE_TAGGED
logme Listing available images
exec_cmd gcloud container images list --repository=${DOCKER_REPO}

echo In order to free up the disk space:
echo gcloud container images delete $IMAGE_TAGGED

