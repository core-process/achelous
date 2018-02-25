#!/bin/bash

# setup environment
export GOPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo "> setting GOPATH=\"$GOPATH\""

# run default shell
USER_SHELL="$(getent passwd "$USER" | awk -F: '{print $NF}')"
echo "> opening workspace with shell $USER_SHELL..."
eval "$USER_SHELL"
echo "> workspace closed with code $?"
