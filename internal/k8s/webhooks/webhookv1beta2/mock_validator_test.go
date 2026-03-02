// SPDX-License-Identifier:Apache-2.0

package webhookv1beta2

import (
	"errors"

	metallbv1 "go.universe.tf/metallb/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type mockValidator struct {
	bgpPeers   *metallbv1.BGPPeerList
	forceError bool
}

func (m *mockValidator) Validate(objects ...client.ObjectList) error {
	for _, obj := range objects { // assuming one object per type
		switch list := obj.(type) {
		case *metallbv1.BGPPeerList:
			m.bgpPeers = list
		default:
			panic("unexpected type")
		}
	}

	if m.forceError {
		return errors.New("error!")
	}
	return nil
}
