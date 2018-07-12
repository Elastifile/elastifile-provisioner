EMS_URL="https://EEE"
NFS_ADDR="n.n.n.n"

APP_TOOLS_DIR=/home/jeans/Documents/workspace/gke-repos/marketplace-k8s-app-tools-0.5

set -x
$APP_TOOLS_DIR/scripts/start.sh --deployer=gcr.io/launcher-poc-207208/elastifile-provisioner-deployer --parameters='{"name":"elastifile-provisioner-1","emanageAddress":"EMS_URL","emanagePassword":"changeme","emanageUser":"admin","imageProvisioner":"docker.io/elastifileio/provisioner:defaultstorageclass","namespace":"default","nfsAddress":"NFS_ADDR"}'
set +x
