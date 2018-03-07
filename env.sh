#!/bin/bash
set -e

# switch to workspace
cd "$(dirname "${BASH_SOURCE[0]}")"
WORKSPACE="$(pwd)"

# detect target command
if [ "$#" -ne 0 ];
then
    TARGET=("$@")
else
    TARGET=("$(getent passwd "$USER" | awk -F: '{print $NF}')")
fi

# show commands
set -x

# run target command
GOPATH="$WORKSPACE" exec ${TARGET[@]}
