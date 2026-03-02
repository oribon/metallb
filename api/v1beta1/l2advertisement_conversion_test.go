// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	"reflect"
	"testing"

	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestL2AdvertisementConvertTo(t *testing.T) {
	var res metallbv1.L2Advertisement

	src := L2Advertisement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "l2adv1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: L2AdvertisementSpec{
			IPAddressPools: []string{"pool1"},
			IPAddressPoolSelectors: []metav1.LabelSelector{
				{MatchLabels: map[string]string{"env": "prod"}},
			},
			NodeSelectors: []metav1.LabelSelector{
				{MatchLabels: map[string]string{"node": "worker"}},
			},
			Interfaces: []string{"eth0", "eth1"},
			ServiceSelectors: []metav1.LabelSelector{
				{MatchLabels: map[string]string{"app": "web"}},
			},
		},
	}

	expected := metallbv1.L2Advertisement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "l2adv1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.L2AdvertisementSpec{
			IPAddressPools: []string{"pool1"},
			IPAddressPoolSelectors: []metav1.LabelSelector{
				{MatchLabels: map[string]string{"env": "prod"}},
			},
			NodeSelectors: []metav1.LabelSelector{
				{MatchLabels: map[string]string{"node": "worker"}},
			},
			Interfaces: []string{"eth0", "eth1"},
			ServiceSelectors: []metav1.LabelSelector{
				{MatchLabels: map[string]string{"app": "web"}},
			},
		},
	}

	if err := src.ConvertTo(&res); err != nil {
		t.Fatalf("ConvertTo failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected L2Advertisement different than converted")
	}
}

func TestL2AdvertisementConvertFrom(t *testing.T) {
	var res L2Advertisement

	src := metallbv1.L2Advertisement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "l2adv1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.L2AdvertisementSpec{
			IPAddressPools: []string{"pool1"},
			Interfaces:     []string{"eth0"},
		},
	}

	expected := L2Advertisement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "l2adv1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: L2AdvertisementSpec{
			IPAddressPools: []string{"pool1"},
			Interfaces:     []string{"eth0"},
		},
	}

	if err := res.ConvertFrom(&src); err != nil {
		t.Fatalf("ConvertFrom failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected L2Advertisement different than converted")
	}
}
