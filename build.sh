#!/bin/bash
set -e

# prepare variables
WORKSPACE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# run command in environment
exec "$WORKSPACE/env.sh" go install achelous
