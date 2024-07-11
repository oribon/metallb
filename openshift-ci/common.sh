#!/usr/bin/bash
set -euo pipefail

pushd ../../
git log -1 || true # just printing commit in the test output
# we move to content of https://github.com/openshift-metal3/dev-scripts.git repo
# we need to change folder as {common,network}.sh have source files
# shellcheck source=network.sh #https://github.com/koalaman/shellcheck/wiki/SC1090
source ./common.sh
# shellcheck source=network.sh
source ./network.sh
popd
