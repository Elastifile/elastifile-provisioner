In order to test the Elastifile Provisioner deployer, you need to execute run.sh

Example:
./run.sh -e 1.2.3.4 -n 1.2.3.5

Requirements:
* gcloud installed
* kubectl installed
* kubectl configured to use the relevant project by default (use "gcloud container clusters get-credentials")
* default service account configured with necessary permissions, see admin-pod-rbac.yaml for an example

In order to clean-up the environment created by run.sh, run dismantle.sh

