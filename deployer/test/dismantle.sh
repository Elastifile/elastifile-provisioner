#!/bin/bash

MYPATH=$(dirname $0)
CONFIGMAP=elastifile-provisioner-deployer-config

source $MYPATH/../functions.sh

log_info "Deleting deployer pod (might not exist at this point)"
exec_cmd kubectl delete pod elastifile-provisioner-deployer

#log_info "Deleting provisioner deployment (might not exist at this point)"
#exec_cmd kubectl delete deployment elastifile-provisioner

log_info "Deleting fake config map"
exec_cmd kubectl delete configmap $CONFIGMAP

log_info "Deleting provisioner deployment"
exec_cmd $MYPATH/../deploy.sh -d

log_info "Deleting test PVC"
exec_cmd kubectl delete -f $MYPATH/pvc.yaml

log_info "Deleting application CRD"
exec_cmd kubectl delete -f "$MYPATH/../config/app-crd.yaml"

