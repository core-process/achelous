#!/bin/bash
set -e

# setup environment
export GOPATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# run default shell
eval "$(getent passwd "$USER" | awk -F: '{print $NF}')"
