#!/bin/bash

function logme {
    msg=$*
    echo $(date '+%Y-%m-%d %H:%M:%S') === $msg
    topic="script"
    if [ -n "$MYNAME" ]; then
        topic=$MYNAME
    fi
    logger -t $topic "$msg"
}

function exec_cmd () {
    cmd=$*
    logme "Executing: $cmd"
    $cmd
    rv=$?
    if ((rv!=0)); then
        logme "Command execution failed with exit code: $rv"
    else
        logme Command executed successfully
    fi
}

