#!/bin/bash

MYPATH=$(dirname $0)
source $MYPATH/../functions.sh

if [ -z "$1" ]; then
    logme "$(basename $0) - EROR: EMS_ADDR not specified"
    exit 2
fi
if [ -z "$2" ]; then
    logme "$(basename $0) - EROR: NFS_ADDR not specified"
    exit 2
fi
EMS_ADDR=$1
NFS_ADDR=$2
EMS_USER=admin
EMS_PASS=changeme
NAMESPACE=default
PROJECT=$(gcloud config get-value project)
IMAGE=gcr.io/$PROJECT/provisioner:1.0

CONFIGMAP=elastifile-provisioner-deployer-config

exec_cmd kubectl delete configmap $CONFIGMAP
exec_cmd kubectl create configmap $CONFIGMAP --from-literal=namespace=$NAMESPACE --from-literal=imageProvisioner=$IMAGE --from-literal=emanageAddress=$EMS_ADDR --from-literal=nfsAddress=$NFS_ADDR --from-literal=emanageUser=$EMS_USER --from-literal=emanagePassword=$EMS_PASS

