#!/usr/bin/env bash -x

MYPATH=$(dirname $0)
source ${MYPATH}/functions.sh
source /home/jeans/Documents/workspace/elastifile-provisioner/deployer/functions.sh

function string_starts_with {
    local string="$1"
    local substr="$2"

    echo "Checking $string"
    shopt -s nocasematch
    case "$string" in
     $substr* ) return 0;;
     *) return 1;;
    esac
}

function validate_https {
    local str="$1"
    local https="https://"
    string_starts_with "$str" ${https}
    local res=$?
    if [ $res -ne 0 ]; then
        log_error "'${str}' doesn't start with https:// "
    fi
    return ${res}
}
