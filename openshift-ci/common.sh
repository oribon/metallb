#!/usr/bin/bash
set -euo pipefail

metallb_dir="$(dirname $(readlink -f $0))"
source ${metallb_dir}/../../common.sh
source ${metallb_dir}/../../network.sh
