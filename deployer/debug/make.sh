#!/bin/bash

MYPATH=$(dirname $0)
source $MYPATH/../functions.sh

# TODO: Make customizable values actually customaizable

# Customizable values
PROJECT=launcher-poc-207208
#PROJECT=elastifile-gce-lab-c304
TAG=latest
DOCKER_FILE=Dockerfile

# Static values
IMAGE=elastifile-provisioner-debug

IMAGE_TAGGED=gcr.io/$PROJECT/$IMAGE:$TAG
logme "Building "$DOCKER_FILE" and tagging the image as $IMAGE_TAGGED"
exec_cmd docker build -t $IMAGE_TAGGED . -f $DOCKER_FILE
logme Pushing $IMAGE_TAGGED
exec_cmd docker push $IMAGE_TAGGED
logme Listing available images
exec_cmd gcloud container images list --repository=gcr.io/$PROJECT

echo In order to free up the disk space:
echo gcloud container images delete $IMAGE_TAGGED

