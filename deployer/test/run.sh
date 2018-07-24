#!/bin/bash

MYNAME=$(basename $0)
MYPATH=$(dirname $0)
. $MYPATH/../functions.sh
. $MYPATH/../validators.sh

OPTS=$(getopt -o e:n: -n $MYNAME -- "$@")
if [ $? != 0 ] ; then
    logme "ERROR: Failed parsing command line arguments"
    exit 2
fi

eval set -- "$OPTS"
while [ $# -gt 0 ]; do
  case "$1" in
    -e) EMS_ADDR=$2
        shift 2 ;;
    -n) NFS_ADDR=$2
        shift 2 ;;
    *) break ;;
  esac
done

assert_var_not_empty EMS_ADDR
assert_var_not_empty NFS_ADDR

logme "Configuring application CRD (Custom Resource Definition)"
exec_cmd kubectl apply -f "$MYPATH/../config/app-crd.yaml"

logme "Faking config map"
$MYPATH/create_config_map.sh $EMS_ADDR $NFS_ADDR

logme "Starting deployer pod"
# TODO: Update the image location with $PROJECT
exec_cmd kubectl create -f $MYPATH/deployer-pod.yaml

logme "You should be able to monitor the pods' status via 'kubectl get pod' and 'kubectl logs <pod name>'"
