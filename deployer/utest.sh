#!/usr/bin/env bash

# Fake MYPATH
MYPATH=$(dirname $0)

source ${MYPATH}/functions.sh
source ${MYPATH}/validators.sh

VAR1=blah
VAR2=" "
echo $(assert_var_not_empty VAR1)
echo $(assert_var_not_empty VAR2)
echo $(assert_var_not_empty VAR3)

