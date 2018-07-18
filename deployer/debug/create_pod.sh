#!/usr/bin/env bash

if [ -n "$1" ]; then
    TAG="$1"
else
    TAG=latest
fi

PROJECT=$(gcloud config get-value project)
IMAGE_NAME=nfs-debug
YAML_FILE=debug-pod.yaml

echo Creating pod ${IMAGE_NAME}:${TAG} in project ${PROJECT}
PROJECT=${PROJECT} IMAGE_NAME=${IMAGE_NAME} TAG=${TAG} envsubst < templates/${YAML_FILE} > /tmp/${YAML_FILE}
kubectl create -f /tmp/${YAML_FILE}

