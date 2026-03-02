// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	"reflect"
	"testing"
	"time"

	metallbv1 "go.universe.tf/metallb/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	MetalLBTestNameSpace = "metallb-test-namespace"
)

func TestBGPPeerConvertTo(t *testing.T) {
	var resBGPPeer metallbv1.BGPPeer

	convertBGPPeer := BGPPeer{
		ObjectMeta: v1.ObjectMeta{
			Name:      "peer1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: BGPPeerSpec{
			MyASN:         42,
			ASN:           142,
			Address:       "1.2.3.4",
			Port:          1179,
			HoldTime:      v1.Duration{Duration: 180 * time.Second},
			RouterID:      "10.20.30.40",
			SrcAddress:    "10.20.30.40",
			EBGPMultiHop:  true,
			BFDProfile:    "default",
			KeepaliveTime: v1.Duration{Duration: time.Second},
			Password:      "nopass",
			NodeSelectors: []NodeSelector{
				{
					MatchLabels: map[string]string{
						"foo": "bar",
					},
					MatchExpressions: []MatchExpression{
						{
							Operator: "In",
							Values:   []string{"quux"},
						},
					},
				},
			},
		},
	}

	expectedBGPPeer := metallbv1.BGPPeer{
		ObjectMeta: v1.ObjectMeta{
			Name:      "peer1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.BGPPeerSpec{
			MyASN:      42,
			ASN:        142,
			Address:    "1.2.3.4",
			Port:       1179,
			HoldTime:   &v1.Duration{Duration: 180 * time.Second},
			RouterID:   "10.20.30.40",
			SrcAddress: "10.20.30.40",
			EBGPMultiHop: true,
			Password:     "nopass",
			BFDProfile:   "default",
			KeepaliveTime: &v1.Duration{Duration: time.Second},
			NodeSelectors: []v1.LabelSelector{
				{
					MatchLabels: map[string]string{
						"foo": "bar",
					},
					MatchExpressions: []v1.LabelSelectorRequirement{
						{
							Operator: "In",
							Values:   []string{"quux"},
						},
					},
				},
			},
		},
	}

	err := convertBGPPeer.ConvertTo(&resBGPPeer)
	if err != nil {
		t.Fatalf("failed converting BGPPeer to v1: %s", err)
	}

	if !reflect.DeepEqual(resBGPPeer, expectedBGPPeer) {
		t.Fatalf("expected BGPPeer different than converted")
	}
}

func TestBGPPeerConvertFrom(t *testing.T) {
	var resBGPPeer BGPPeer

	convertBGPPeer := metallbv1.BGPPeer{
		ObjectMeta: v1.ObjectMeta{
			Name:      "peer1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.BGPPeerSpec{
			MyASN:        42,
			ASN:          142,
			Address:      "1.2.3.4",
			Port:         1179,
			HoldTime:     &v1.Duration{Duration: 180 * time.Second},
			RouterID:     "10.20.30.40",
			SrcAddress:   "10.20.30.40",
			EBGPMultiHop: true,
			Password:     "nopass",
			BFDProfile:   "default",
			KeepaliveTime: &v1.Duration{Duration: time.Second},
			NodeSelectors: []v1.LabelSelector{
				{
					MatchLabels: map[string]string{
						"foo": "bar",
					},
					MatchExpressions: []v1.LabelSelectorRequirement{
						{
							Operator: "In",
							Values:   []string{"quux"},
						},
					},
				},
			},
		},
	}

	expectedBGPPeer := BGPPeer{
		ObjectMeta: v1.ObjectMeta{
			Name:      "peer1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: BGPPeerSpec{
			MyASN:         42,
			ASN:           142,
			Address:       "1.2.3.4",
			Port:          1179,
			HoldTime:      v1.Duration{Duration: 180 * time.Second},
			RouterID:      "10.20.30.40",
			SrcAddress:    "10.20.30.40",
			EBGPMultiHop:  true,
			BFDProfile:    "default",
			KeepaliveTime: v1.Duration{Duration: time.Second},
			Password:      "nopass",
			NodeSelectors: []NodeSelector{
				{
					MatchLabels: map[string]string{
						"foo": "bar",
					},
					MatchExpressions: []MatchExpression{
						{
							Operator: "In",
							Values:   []string{"quux"},
						},
					},
				},
			},
		},
	}

	err := resBGPPeer.ConvertFrom(&convertBGPPeer)
	if err != nil {
		t.Fatalf("failed converting v1 BGPPeer: %s", err)
	}

	if !reflect.DeepEqual(resBGPPeer, expectedBGPPeer) {
		t.Fatalf("expected BGPPeer different than converted")
	}
}
