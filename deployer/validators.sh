#!/usr/bin/env bash

MYPATH=$(dirname $0)
source ${MYPATH}/functions.sh

function string_starts_with {
    local string="$1"
    local substr="$2"

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

function assert_non_empty {
    local desc="$1"
    local value="$2"
    [ "${value}"x != ""x ]
    assert $? "${desc} is empty"
}
