// SPDX-License-Identifier:Apache-2.0

package webhookv1beta1

import (
	"errors"

	metallbv1 "go.universe.tf/metallb/api/v1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type mockValidator struct {
	ipAddressPools *metallbv1.IPAddressPoolList
	bgpAdvs        *metallbv1.BGPAdvertisementList
	l2Advs         *metallbv1.L2AdvertisementList
	communities    *metallbv1.CommunityList
	nodes          *v1.NodeList
	forceError     bool
}

func (m *mockValidator) Validate(objects ...client.ObjectList) error {
	for _, obj := range objects { // assuming one object per type
		switch list := obj.(type) {
		case *metallbv1.BGPAdvertisementList:
			m.bgpAdvs = list
		case *metallbv1.L2AdvertisementList:
			m.l2Advs = list
		case *metallbv1.IPAddressPoolList:
			m.ipAddressPools = list
		case *metallbv1.CommunityList:
			m.communities = list
		case *v1.NodeList:
			m.nodes = list
		default:
			panic("unexpected type")
		}
	}

	if m.forceError {
		return errors.New("error!")
	}
	return nil
}
