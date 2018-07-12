PROJECT=launcher-poc-207208
APP_TOOLS_DIR=/home/jeans/Documents/workspace/gke-repos/marketplace-k8s-app-tools-0.5

set -x
$APP_TOOLS_DIR/scripts/start.sh --deployer=gcr.io/$PROJECT/elastifile-provisioner-deployer --parameters='{"name":"elastifile-provisioner-1","emanageAddress":"https://CHANGEME","emanagePassword":"changeme","emanageUser":"admin","imageProvisioner":"docker.io/elastifileio/provisioner:defaultstorageclass","namespace":"default","nfsAddress":"CHANGEME"}'
set +x

