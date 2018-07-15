#!/usr/bin/env bash

MYPATH=$(dirname $0)
source ${MYPATH}/../functions.sh

assert_exec_cmd kubectl apply -f ${MYPATH}/app-crd.yaml
