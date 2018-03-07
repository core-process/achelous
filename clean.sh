#!/bin/bash
set -e

# switch to workspace
cd "$(dirname "${BASH_SOURCE[0]}")"

# show commands
set -x

# clean
rm -r -f "bin"
rm -r -f "pkg"
rm -r -f "src/achelous/vendor"
