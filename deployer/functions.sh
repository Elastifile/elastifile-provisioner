#!/usr/bin/env bash

function logme {
    msg=$*
    echo $(date '+%Y-%m-%d %H:%M:%S') === $msg
    topic="script"
    if [ -n "$MYNAME" ]; then
        topic=$MYNAME
    fi
    logger -t $topic "$msg"
}

function log_info {
    logme "INFO: $*"
}

function log_error {
    logme "ERROR: $*"
}

function exec_cmd {
    cmd=$*
    logme "Executing: $cmd"
    ${cmd}
    rv=$?
    if ((rv != 0)); then
        logme "Command execution failed with exit code: $rv"
    else
        logme Command executed successfully
    fi
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
