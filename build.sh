#!/bin/bash
set -e

# switch to workspace
cd "$(dirname "${BASH_SOURCE[0]}")"
WORKSPACE="$(pwd)"

# show commands
set -x

# install go dependencies
pushd src/achelous
GOPATH="$WORKSPACE" glide install
popd

# build go sources
GOPATH="$WORKSPACE" go install achelous/spring-core achelous/upstream-core

# build c sources
gcc src/achelous/wrapper/main.c -o bin/spring
cp bin/spring bin/upstream
