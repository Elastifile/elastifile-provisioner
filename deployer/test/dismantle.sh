#!/bin/bash

MYPATH=$(dirname $0)
CONFIGMAP=elastifile-provisioner-deployer-config

source $MYPATH/../functions.sh

logme "Deleting deployer pod (might not exist at this point)"
exec_cmd kubectl delete pod elastifile-provisioner-deployer

#logme "Deleting provisioner deployment (might not exist at this point)"
#exec_cmd kubectl delete deployment elastifile-provisioner

logme "Deleting fake config map"
exec_cmd kubectl delete configmap $CONFIGMAP

logme "Deleting provisioner deployment"
exec_cmd $MYPATH/../deploy.sh -d

logme "Deleting test PVC"
exec_cmd kubectl delete -f $MYPATH/pvc.yaml

logme "Deleting application CRD"
exec_cmd kubectl delete -f "$MYPATH/../config/app-crd.yaml"

