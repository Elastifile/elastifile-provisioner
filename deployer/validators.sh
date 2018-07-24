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

# Takes variable description and value and asserts is the latter is empty - can be used with out-of-scope values
function assert_var_not_empty_by_val {
    local var_desc="$1"
    local var_value="$2"

    # Here we check if the variable is non-zero after spaces are removed
    [[ -n "${var_value## }" ]]
    assert $? "${var_desc} is empty"
}

# Takes global variable name and asserts is the value is empty
function assert_var_not_empty () {
    local var_name=$1
    local var_value=${!var_name}
    assert_var_not_empty_by_val "Variable ${var_name}" ${var_value}
}
