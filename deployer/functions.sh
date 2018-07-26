#!/usr/bin/env bash

function logme {
    local msg=$*
    echo $(date '+%Y-%m-%d %H:%M:%S') $msg
    local topic="script"
    if [ -n "${MYNAME}" ]; then
        topic=${MYNAME}
    fi
    logger -t ${topic} "$msg"
}

function log_error {
    logme "ERROR: $*"
}

function log_info {
    logme "INFO:  $*"
}

function log_debug {
    logme "DEBUG: $*"
}

function exec_cmd {
    local cmd=$@
    log_info "Executing: ${cmd}"
    ${cmd}
    local res=$?
    if ((res != 0)); then
        log_info "Command execution failed with exit code: ${res}"
    else
        log_info "Command executed successfully"
    fi
    return ${res}
}

function assert {
    local error=$1
    shift
    local desc="$*"
    if [ "$error" -ne 0 ]; then
        log_error "$desc"
        exit ${error}
    fi
}

function assert_exec_cmd {
    exec_cmd $@
    assert $? "Command failed with exit code $?: $@"
}
