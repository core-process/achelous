#!/bin/bash
set -e

# prepare variables
WORKSPACE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# build go sources
echo "Building Go sources..."
"$WORKSPACE/env.sh" go install achelous/spring-core achelous/upstream-core

# build c sources
echo "Building C sources..."
cd "$WORKSPACE" && mkdir -p bin && gcc src/achelous/wrapper.c -o bin/spring && cp bin/spring bin/upstream

echo "Done!"
