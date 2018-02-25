#!/bin/bash
set -e

# prepare variables
WORKSPACE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [ "$#" -ne 0 ];
then
    TARGET=("$@")
else
    TARGET=("$(getent passwd "$USER" | awk -F: '{print $NF}')")
fi

# prepare environment and run command
set -x
cd "$WORKSPACE"
GOPATH="$WORKSPACE" exec ${TARGET[@]}
