#!/usr/bin/env bash

MYNAME=$(basename $0)
MYPATH=$(dirname $0)
APP_TOOLS_DIR=${MYPATH}/marketplace-k8s-app-tools

. ${MYPATH}/functions.sh
. ${MYPATH}/validators.sh

# Defaults
PROJECT=$(gcloud config get-value project)
DOCKER_REPO=gcr.io/${PROJECT}
TAG=1.0
EMS_ADDR='https://FAKE_EMS_URL'
NFS_ADDR='FAKE_NFS_ADDR'
NAMESPACE='default'
SERVICE_ACCOUNT='default'
EMS_USER='admin'
EMS_PASS='changeme'
APP_NAME='test-prov-1'

OPTS=$(getopt -o e:n:m:t:s:a:u:p: -n ${MYNAME} -- "$@")
if [ $? != 0 ] ; then
    log_error "Failed parsing command line arguments"
    exit 2
fi

eval set -- "${OPTS}"
while [ $# -gt 0 ]; do
  case "$1" in
    -e) EMS_ADDR=$2
        shift 2 ;;
    -n) NFS_ADDR=$2
        shift 2 ;;
    -m) APP_TOOLS_DIR=$2
        shift 2 ;;
    -t) TAG=$2
        shift 2 ;;
    -s) NAMESPACE=$2
        shift 2 ;;
    -a) SERVICE_ACCOUNT=$2
        shift 2 ;;
    -u) EMS_USER="$2"
        shift 2 ;;
    -p) EMS_PASS="$2"
        shift 2 ;;
    *) break ;;
  esac
done

DEPLOYER_IMAGE=${DOCKER_REPO}/provisioner/deployer:${TAG}

IFS='' read -r -d '' PARAMETERS <<EOF
{"name":"${APP_NAME}","emanageAddress":"${EMS_ADDR}","emanagePassword":"${EMS_PASS}","emanageUser":"${EMS_USER}","imageProvisioner":"${DOCKER_REPO}/provisioner:${TAG}","imageCustomScript":"${DOCKER_REPO}/provisioner/custom_script:${TAG}","namespace":"${NAMESPACE}","nfsAddress":"${NFS_ADDR}","serviceAccount":"${SERVICE_ACCOUNT}"}
EOF

set -x
# Creating app resource
${MYPATH}/crd/set.sh

# Setting permissions - Failure for existing resource is ok here
kubectl create -f ${MYPATH}/test/admin-pod-rbac.yaml

# Executing start.sh
${APP_TOOLS_DIR}/scripts/start.sh --deployer=${DEPLOYER_IMAGE} --parameters=${PARAMETERS}
set +x
