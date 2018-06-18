#!/bin/bash

#PROJECT=elastifile-gce-lab-c304
#CLUSTER=test-gke-in-same-network3
#ZONE=europe-west1-b
PROJECT=launcher-poc-207208
CLUSTER=gke-cluster-jean
ZONE=us-west1-a

gcloud config set project $PROJECT
gcloud container clusters get-credentials "$CLUSTER" --zone "$ZONE"
kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user $(gcloud config get-value account)
gcloud auth configure-docker

