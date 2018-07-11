#!/usr/bin/env bash

# Current script name
MYNAME=$(basename $0)
MYPATH=$(dirname $0)
# Directory where deployment files are found
DEPLOYDIR=/tmp/manifest4script
# Name is derived from the app name
CONFIGMAP=elastifile-provisioner-deployer-config
# Dir is set in deployment YAML (mountPath)
CONFIGDIR=/etc/config

function log_info {
    echo "INFO: $*"
}

function log_error {
    echo "ERROR: $*" >&2
}

function get_file_conf {
    FNAME=$CONFIGDIR/$1
    if [ ! -f "$FNAME" ]; then
        log_error "File $FNAME doesn't exist"
        exit 11
    fi
    cat $FNAME
}

function assert {
    local error=$1
    shift
    local desc="$*"
    if [ "$error" -ne 0 ]; then
        log_error "$desc"
        exit $error
    fi
}

function run_cmd {
    log_info "Executing $@"
    "$@"
}

function assert_run_cmd {
    run_cmd $@
    assert $? "Command failed with exit code $exitcode: $@"
}

function deploy_provisioner {
    log_info "Deploying ECFS provisioner"
    CUSTOM_FLAGS=""
    if [ "$DRY_RUN" = true ]; then
        echo "WARNING: DRY RUN"
        CUSTOM_FLAGS="--dry-run"
    fi
    #assert_run_cmd kubectl create -f $DEPLOYDIR/serviceaccount.yaml -n $NAMESPACE $CUSTOM_FLAGS
    assert_run_cmd kubectl create -f $DEPLOYDIR/storageclass.yaml -n $NAMESPACE $CUSTOM_FLAGS
    #assert_run_cmd kubectl create -f $DEPLOYDIR/clusterrole.yaml -n $NAMESPACE $CUSTOM_FLAGS
    #assert_run_cmd kubectl create -f $DEPLOYDIR/clusterrolebinding.yaml -n $NAMESPACE $CUSTOM_FLAGS
    #assert_run_cmd kubectl create -f $DEPLOYDIR/deployment.yaml $CUSTOM_FLAGS
    #assert_run_cmd kubectl patch deployment elastifile-provisioner -p '{"spec":{"template":{"spec":{"serviceAccount":"elastifile-provisioner"}}}}'
    assert_run_cmd kubectl create secret generic elastifile-rest --from-literal="password.txt=$ECFS_PASS"
}

function destroy_provisioner_configuration {
    log_info "Destroying ECFS provisioner deployment"
    set -x
    run_cmd kubectl delete -f $DEPLOYDIR/deployment.yaml
    run_cmd kubectl delete -f $DEPLOYDIR/clusterrolebinding.yaml -n $NAMESPACE
    run_cmd kubectl delete -f $DEPLOYDIR/clusterrole.yaml -n $NAMESPACE
    run_cmd kubectl delete -f $DEPLOYDIR/serviceaccount.yaml -n $NAMESPACE
    run_cmd kubectl delete -f $DEPLOYDIR/storageclass.yaml -n $NAMESPACE
    run_cmd kubectl delete secret elastifile-rest
    set +x
}

DRY_RUN=false
DESTROY=false
KEEP_ALIVE=false

OPTS=$(getopt -o dnk -n $MYNAME -- "$@")
if [ $? != 0 ] ; then
    log_error "Failed parsing command line arguments"
    exit 2
fi

eval set -- "$OPTS"
while true; do
  case "$1" in
    -n) DRY_RUN=true
        shift ;;
    -d) DESTROY=true
        DEPLOYDIR=$MYPATH/config
        shift ;;
    -k) KEEP_ALIVE=true
        shift ;;
    *) break ;;
  esac
done

if [ "$DESTROY" != true ]; then
    # Fetch user input from configMap
    NAMESPACE=$(get_file_conf namespace)
    #assert $? "Failed getting NAMESPACE"
    NFS_ADDR=$(get_file_conf nfsAddress)
    assert $? "Failed getting nfsAddress"
    EMS_ADDR=$(get_file_conf emanageAddress)
    assert $? "Failed getting emanageAddress"
    ECFS_USER=$(get_file_conf emanageUser)
    assert $? "Failed getting emanageUser"

    IS_PASSWORD_BASE64=false
    if [ "$IS_PASSWORD_BASE64" = true ]; then
        # TODO: Use secrets to store eManage password
        ECFS_PASS_BASE64=$(get_file_conf emanagePassword)
        assert $? "Failed getting emanagePassword"
        # Decode the password (if acquired from a secret)
        ECFS_PASS=$(echo -n "$ECFS_PASS_BASE64" | base64 -d)
    else
        ECFS_PASS=$(get_file_conf emanagePassword)
        assert $? "Failed getting emanagePassword"
    fi
fi

if [ -z $NAMESPACE ]; then
    NAMESPACE=default
fi

# Update the configuration
if [ -n "$NAMESPACE" ]; then
    log_info "Updating YAML files with namespace: $NAMESPACE"
    sed -ie "s/^\\(\s*namespace:\s*\\)default/\\1$NAMESPACE/" $DEPLOYDIR/clusterrolebinding.yaml
    sed -ie "s/^\\(\s*secretNamespace:\s*\\).*/\\1\"$NAMESPACE\"/" $DEPLOYDIR/storageclass.yaml
fi
if [ -n "$NFS_ADDR" ]; then
    log_info "Updating YAML files with NFS Server: $NFS_ADDR"
    sed -ie "s/^\\(\s*nfsServer:\s*\\).*/\\1\"$NFS_ADDR\"/" $DEPLOYDIR/storageclass.yaml
fi
if [ -n "$EMS_ADDR" ]; then
    EMS_URL="https:\\/\\/$EMS_ADDR"
    log_info "Updating YAML files with eManage URL: $EMS_URL"
    sed -ie "s/^\\(\s*restURL:\s*\\).*/\\1\"$EMS_URL\"/" $DEPLOYDIR/storageclass.yaml
fi
if [ -n "$ECFS_USER" ]; then
    log_info "Updating YAML files with username: $ECFS_USER"
    sed -ie "s/^\\(\s*username:\s*\\).*/\\1\"$ECFS_USER\"/" $DEPLOYDIR/storageclass.yaml
fi

# Deploy/destroy provisioner
if [ ! "$DESTROY" = true ]; then 
    deploy_provisioner
else
    destroy_provisioner_configuration
fi

echo "Deployment completed"
while [ "$KEEP_ALIVE" = true ]; do
    sleep 1
done

