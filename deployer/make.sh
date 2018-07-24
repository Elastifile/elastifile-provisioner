#!/bin/bash

MYNAME=$(basename $0)
MYPATH=$(dirname $0)

source ${MYPATH}/functions.sh

# Defaults
PROJECT=$(gcloud config get-value project)
TAG=1.0

# Command-line arguments
OPTS=$(getopt -o p:t: -n ${MYNAME} -- "$@")
if [ $? != 0 ] ; then
    logme "ERROR: Failed parsing command line arguments"
    exit 2
fi

eval set -- "${OPTS}"
while [ $# -gt 0 ]; do
  case "$1" in
    -p) PROJECT="$2"
        shift 2 ;;
    -t) TAG="$2"
        shift 2 ;;
    *) break ;;
  esac
done

APP_IMAGE_NAME=provisioner
REGISTRY=gcr.io/${PROJECT}
DEPLOYER_IMAGE_TAGGED=${REGISTRY}/${APP_IMAGE_NAME}/deployer:${TAG}
IMAGE_CLONE_TAGGED=${REGISTRY}/${APP_IMAGE_NAME}/custom_script:${TAG} # Fake image used to run deploy.sh

logme "Building and tagging the deployer image as ${DEPLOYER_IMAGE_TAGGED}"
assert_exec_cmd docker build -t ${DEPLOYER_IMAGE_TAGGED} --build-arg REGISTRY=${REGISTRY} --build-arg TAG=${TAG} .
logme "Pushing ${DEPLOYER_IMAGE_TAGGED}"
assert_exec_cmd docker push ${DEPLOYER_IMAGE_TAGGED}

logme "Tagging and pushing ${IMAGE_CLONE_TAGGED}"
assert_exec_cmd docker tag ${DEPLOYER_IMAGE_TAGGED} ${IMAGE_CLONE_TAGGED}
assert_exec_cmd docker push ${IMAGE_CLONE_TAGGED}

logme "Listing available images"
exec_cmd gcloud container images list --repository=${REGISTRY}/${APP_IMAGE_NAME}
