#!/usr/bin/bash

set -euo pipefail

metallb_dir="$(dirname $(readlink -f $0))"
source ${metallb_dir}/common.sh

FRR_IMAGE_BASE=${FRR_IMAGE_BASE:-$(echo "${OPENSHIFT_RELEASE_IMAGE}" | sed -e 's/release/stable/g' | sed -e 's/@.*$//g')}
FRR_IMAGE_TAG=${FRR_IMAGE_TAG:-"metallb-frr"}

# - Change the feature gate to enable frrk8s deployed by CNO
${metallb_dir}/enable_frrk8s_on_cno.sh

FRRK8S_NAMESPACE="openshift-frr-k8s"

oc patch networks.operator.openshift.io cluster --type json  -p '[{"op": "add", "path": "/spec/additionalRoutingCapabilities", "value": {providers: ["FRR"]}}]'

wait_for_pods $FRRK8S_NAMESPACE "app=frr-k8s"

sudo ip route add 192.168.10.0/24 dev ${BAREMETAL_NETWORK_NAME}
