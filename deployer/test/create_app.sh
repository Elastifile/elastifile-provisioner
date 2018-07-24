#!/usr/bin/env bash

MYPATH=$(dirname $0)

APP_TOOLS_DIR=/home/jeans/Documents/workspace/gke-repos/marketplace-k8s-app-tools
PROJECT=$(gcloud config get-value project)
DOCKER_REPO=gcr.io/${PROJECT}
DEPLOYER_IMAGE=${DOCKER_REPO}/provisioner/deployer:1.0

IFS='' read -r -d '' PARAMETERS <<EOF
{"name":"ecfs-provisioner-via-start-1","emanageAddress":"https://CHANGEME","emanagePassword":"changeme","emanageUser":"admin","imageProvisioner":"${DOCKER_REPO}/provisioner:1.0","imageCustomScript":"${DOCKER_REPO}/provisioner/custom_script:1.0","namespace":"default","nfsAddress":"CHANGEME","serviceAccount":"default"}
EOF

set -x

# Creating app resource
${MYPATH}/../crd/set.sh

# Setting permissions - Failure for existing resource is ok here
kubectl create -f ${MYPATH}/admin-pod-rbac.yaml

# Executing start.sh
${APP_TOOLS_DIR}/scripts/start.sh --deployer=${DEPLOYER_IMAGE} --parameters=${PARAMETERS}
set +x

