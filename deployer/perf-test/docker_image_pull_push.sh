#!/usr/bin/env bash
SRC_IMAGE=$1
PROJECT=$2

# This script will take an image from an external repo and push it to gcr.io
if [ "$PROJECT"x == x ]; then
    echo Missing arguments. Usage: $0 IMAGE-NAME[:TAG] PROJECT
    exit 2
fi

DST_REPO=gcr.io/$PROJECT

IMAGE_NAME=$SRC_IMAGE

if [[ $IMAGE_NAME != *":"* ]]; then
    echo Fixing up image name with default tag
    IMAGE_NAME=${IMAGE_NAME}:latest
fi

set -x
docker pull $SRC_IMAGE
docker tag $IMAGE_NAME $DST_REPO/$IMAGE_NAME
docker push $DST_REPO/$IMAGE_NAME
set +x
