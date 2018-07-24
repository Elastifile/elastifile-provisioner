#!/usr/bin/env bash

# This script can be used to pull an updated version of the provisioner docker image
# and push it to the Marketplace staging project

TAG=1.0
IMAGE_NAME=provisioner
PROJECT=$(gcloud config get-value project)
SRC_REPO=docker.io/elastifileio
DST_REPO=gcr.io/$PROJECT

docker pull $SRC_REPO/$IMAGE_NAME:$TAG
docker tag $SRC_REPO/$IMAGE_NAME:$TAG $DST_REPO/$IMAGE_NAME:$TAG
docker push $DST_REPO/$IMAGE_NAME:$TAG

