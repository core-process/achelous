#!/bin/bash
set -e

# prepare variables
WORKSPACE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# remove all
set -x

rm -r -f "$WORKSPACE/bin"
rm -r -f "$WORKSPACE/pkg"
rm -f "$WORKSPACE/src/achelous/debug"
