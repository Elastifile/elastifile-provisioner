#!/bin/bash

#PROJECT=elastifile-gce-lab-c304
#CLUSTER=test-gke-in-same-network3
#ZONE=europe-west1-b

#PROJECT=launcher-poc-207208
#CLUSTER=gke-cluster-jean
#ZONE=us-west1-a

#PROJECT=launcher-poc-207208
#CLUSTER=launcher-poc-cluster
#ZONE=europe-west1-b

#PROJECT=launcher-poc-207208
#CLUSTER=cluster-2
#ZONE=us-central1-a

PROJECT=launcher-poc-207208
CLUSTER=cluster-tmp
ZONE=us-central1-a

set -x
gcloud config set project $PROJECT
gcloud container clusters get-credentials "$CLUSTER" --zone "$ZONE"
kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user $(gcloud config get-value account)
gcloud auth configure-docker
set +x
