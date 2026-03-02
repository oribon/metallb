// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	"reflect"
	"testing"

	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBGPAdvertisementConvertTo(t *testing.T) {
	var res metallbv1.BGPAdvertisement
	aggLen := int32(24)
	aggLenV6 := int32(64)

	src := BGPAdvertisement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "adv1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: BGPAdvertisementSpec{
			AggregationLength:   &aggLen,
			AggregationLengthV6: &aggLenV6,
			LocalPref:           100,
			Communities:         []string{"65535:65282"},
			IPAddressPools:      []string{"pool1"},
			IPAddressPoolSelectors: []metav1.LabelSelector{
				{MatchLabels: map[string]string{"env": "prod"}},
			},
			NodeSelectors: []metav1.LabelSelector{
				{MatchLabels: map[string]string{"node": "worker"}},
			},
			Peers: []string{"peer1"},
			ServiceSelectors: []metav1.LabelSelector{
				{MatchLabels: map[string]string{"app": "web"}},
			},
		},
	}

	expected := metallbv1.BGPAdvertisement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "adv1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.BGPAdvertisementSpec{
			AggregationLength:   &aggLen,
			AggregationLengthV6: &aggLenV6,
			LocalPref:           100,
			Communities:         []string{"65535:65282"},
			IPAddressPools:      []string{"pool1"},
			IPAddressPoolSelectors: []metav1.LabelSelector{
				{MatchLabels: map[string]string{"env": "prod"}},
			},
			NodeSelectors: []metav1.LabelSelector{
				{MatchLabels: map[string]string{"node": "worker"}},
			},
			Peers: []string{"peer1"},
			ServiceSelectors: []metav1.LabelSelector{
				{MatchLabels: map[string]string{"app": "web"}},
			},
		},
	}

	if err := src.ConvertTo(&res); err != nil {
		t.Fatalf("ConvertTo failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected BGPAdvertisement different than converted")
	}
}

func TestBGPAdvertisementConvertFrom(t *testing.T) {
	var res BGPAdvertisement
	aggLen := int32(24)

	src := metallbv1.BGPAdvertisement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "adv1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.BGPAdvertisementSpec{
			AggregationLength: &aggLen,
			LocalPref:         200,
			Communities:       []string{"65535:65283"},
			IPAddressPools:    []string{"pool2"},
			Peers:             []string{"peer2"},
		},
	}

	expected := BGPAdvertisement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "adv1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: BGPAdvertisementSpec{
			AggregationLength: &aggLen,
			LocalPref:         200,
			Communities:       []string{"65535:65283"},
			IPAddressPools:    []string{"pool2"},
			Peers:             []string{"peer2"},
		},
	}

	if err := res.ConvertFrom(&src); err != nil {
		t.Fatalf("ConvertFrom failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected BGPAdvertisement different than converted")
	}
}
