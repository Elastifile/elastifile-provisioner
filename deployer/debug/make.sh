#!/bin/bash

MYNAME=$(basename $0)
MYPATH=$(dirname $0)
source ${MYPATH}/../functions.sh

# Defaults
PROJECT=$(gcloud config get-value project)
TAG=latest
IMAGE=nfs-debug
DOCKER_FILE=Dockerfile

# Command-line arguments
OPTS=$(getopt -o p:t:i:d: -n ${MYNAME} -- "$@")
if [ $? != 0 ] ; then
    log_error "Failed parsing command line arguments"
    exit 2
fi

eval set -- "${OPTS}"
while [ $# -gt 0 ]; do
  case "$1" in
    -p) PROJECT="$2"
        shift 2 ;;
    -t) TAG="$2"
        shift 2 ;;
    -i) IMAGE="$2"
        shift 2 ;;
    -d) DOCKER_FILE="$2"
        shift 2 ;;
    *) break ;;
  esac
done

DOCKER_REPO=gcr.io/${PROJECT}
IMAGE_TAGGED=${DOCKER_REPO}/${IMAGE}:${TAG}
log_info "Building "${DOCKER_FILE}" and tagging the image as ${IMAGE_TAGGED}"
assert_exec_cmd docker build -t ${IMAGE_TAGGED} . -f ${DOCKER_FILE}
log_info "Pushing ${IMAGE_TAGGED}"
assert_exec_cmd docker push ${IMAGE_TAGGED}
log_info "Listing available images"
assert_exec_cmd gcloud container images list --repository=${DOCKER_REPO}

echo In order to free up the disk space:
echo gcloud container images delete ${IMAGE_TAGGED}
