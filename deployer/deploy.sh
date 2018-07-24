#!/bin/bash -x

# Current script name
MYNAME=$(basename $0)
MYPATH=$(dirname $0)
# Directory where deployment files are found
DEPLOYDIR=/tmp/manifest4script
# Dir is set in deployment YAML (mountPath)
CONFIGDIR=/etc/config

source ${MYPATH}/functions.sh
source ${MYPATH}/validators.sh

function get_file_conf {
    FNAME=${CONFIGDIR}/$1
    if [ ! -f "$FNAME" ]; then
        log_error "File $FNAME doesn't exist"
        exit 11
    fi
    cat ${FNAME}
}

function deploy_provisioner {
    log_info "Deploying ECFS provisioner"
    DRY_RUN_FLAG=""
    if [ "$DRY_RUN" = true ]; then
        log_info "WARNING: DRY RUN"
        DRY_RUN_FLAG="--dry-run"
    fi

    assert_exec_cmd cat ${STORAGECLASS_YAML}

    log_info "Creating storageclass"
    assert_exec_cmd kubectl create -f ${STORAGECLASS_YAML} -n ${NAMESPACE} ${DRY_RUN_FLAG}

    # TODO: Check if special characters break this
    ECFS_PASS_BASE64=$(echo -n "${ECFS_PASS}" | base64)
    SECRET_NAME='elastifile-rest'

    log_info "Creating secret $SECRET_NAME"
    ECFS_PASS_BASE64=$(echo -n "${ECFS_PASS}" | base64)
    cat <<END_OF_SECRET_MANIFEST | kubectl create -f - ${DRY_RUN_FLAG}
    apiVersion: v1
    kind: Secret
    metadata:
      name: "${SECRET_NAME}"
      namespace: "${NAMESPACE}"
      ownerReferences:
      - apiVersion: ${APP_API_VERSION}
        blockOwnerDeletion: true
        kind: Application
        name: ${APP_NAME}
        uid: ${APP_UID}
    data:
      password.txt: "${ECFS_PASS_BASE64}"
END_OF_SECRET_MANIFEST
}

# Parse command line arguments
DRY_RUN=false
KEEP_ALIVE=false
OPTS=$(getopt -o nk -n ${MYNAME} -- "$@")
if [ $? != 0 ] ; then
    log_error "Failed parsing command line arguments"
    exit 2
fi

eval set -- "${OPTS}"
while true; do
  case "$1" in
    -n) DRY_RUN=true
        shift ;;
    -k) KEEP_ALIVE=true
        shift ;;
    *) break ;;
  esac
done

# Fetch user input (and other settings) from configMap
APP_NAME=$(get_file_conf name)
assert $? "Failed getting APP_NAME"
assert_var_not_empty APP_NAME

NAMESPACE=$(get_file_conf namespace)
assert $? "Failed getting NAMESPACE"
assert_var_not_empty NAMESPACE

NFS_ADDR=$(get_file_conf nfsAddress)
assert $? "Failed getting nfsAddress"
assert_var_not_empty NFS_ADDR
# TODO: Rename emanageAddress to emanageUrl

EMS_URL=$(get_file_conf emanageAddress)
assert $? "Failed getting emanageAddress"
assert_var_not_empty EMS_URL
validate_https ${EMS_URL}
assert $? "Management URL should start with HTTPS:// - received ${EMS_URL}"

ECFS_USER=$(get_file_conf emanageUser)
assert $? "Failed getting emanageUser"
assert_var_not_empty ECFS_USER

IS_PASSWORD_BASE64=false
if [ "$IS_PASSWORD_BASE64" = true ]; then
    # TODO: Use secrets to store eManage password, once supported in the Marketplace
    ECFS_PASS_BASE64=$(get_file_conf emanagePassword)
    assert $? "Failed getting emanagePassword"
    # Decode the password (if acquired from a secret)
    ECFS_PASS=$(echo -n "$ECFS_PASS_BASE64" | base64 -d)
else
    ECFS_PASS=$(get_file_conf emanagePassword)
    assert $? "Failed getting emanagePassword"
fi

# Update the configuration
APP_UID=$(kubectl get "applications/$APP_NAME" --namespace="$NAMESPACE" --output=jsonpath='{.metadata.uid}')
assert $? "Failed getting APP_UID"
assert_var_not_empty APP_UID

APP_API_VERSION=$(kubectl get "applications/$APP_NAME" --namespace="$NAMESPACE" --output=jsonpath='{.apiVersion}') # app.k8s.io/v1alpha1
assert $? "Failed getting APP_API_VERSION"
assert_var_not_empty APP_API_VERSION

STORAGECLASS_YAML=${DEPLOYDIR}/storageclass.yaml
STORAGECLASS_TMPL=${STORAGECLASS_YAML}.template

echo Running envsubst
app_api_version=${APP_API_VERSION} name=${APP_NAME} app_uid=${APP_UID} namespace=${NAMESPACE} nfs_addr=${NFS_ADDR} emanage_addr=${EMS_URL} emanage_user=${ECFS_USER} envsubst < ${STORAGECLASS_TMPL} > ${STORAGECLASS_YAML}

# Deploy provisioner
deploy_provisioner

echo "Deployment completed"
while [ "$KEEP_ALIVE" = true ]; do
    sleep 1
done
