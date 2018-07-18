#!/usr/bin/env bash

APP_TOOLS_DIR=/home/jeans/Documents/workspace/gke-repos/marketplace-k8s-app-tools
PROJECT=launcher-poc-207208
DEPLOYER_IMAGE=gcr.io/${PROJECT}/provisioner/deployer:1.0

set -x
${APP_TOOLS_DIR}/scripts/start.sh --deployer=${DEPLOYER_IMAGE} --parameters='{"name":"ecfs-provisioner-via-start-1","emanageAddress":"https://CHANGEME","emanagePassword":"changeme","emanageUser":"admin","imageProvisioner":"gcr.io/launcher-poc-207208/provisioner:1.0","imageCustomScript":"gcr.io/launcher-poc-207208/provisioner/custom_script:1.0","namespace":"default","nfsAddress":"CHANGEME","serviceAccount":"default"}'
set +x
