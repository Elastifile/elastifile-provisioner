#!/usr/bin/env bash

MYPATH=$(dirname $0)
source ${MYPATH}/../functions.sh

assert_exec_cmd kubectl delete -f ${MYPATH}/app-crd.yaml
