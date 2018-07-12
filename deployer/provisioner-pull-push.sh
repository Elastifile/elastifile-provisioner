#!/usr/bin/env bash
TAG=1.0
IMAGE_NAME=provisioner
SRC_REPO=docker.io/elastifileio
DST_REPO=gcr.io/launcher-poc-207208

docker pull $SRC_REPO/$IMAGE_NAME:$TAG
docker tag $SRC_REPO/$IMAGE_NAME:$TAG $DST_REPO/$IMAGE_NAME:$TAG
docker push $DST_REPO/$IMAGE_NAME:$TAG

