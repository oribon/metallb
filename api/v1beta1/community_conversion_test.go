// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	"reflect"
	"testing"

	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCommunityConvertTo(t *testing.T) {
	var res metallbv1.Community

	src := Community{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "community1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: CommunitySpec{
			Communities: []CommunityAlias{
				{Name: "no-advertise", Value: "65535:65282"},
				{Name: "large-test", Value: "large:123:456:789"},
			},
		},
	}

	expected := metallbv1.Community{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "community1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.CommunitySpec{
			Communities: []metallbv1.CommunityAlias{
				{Name: "no-advertise", Value: "65535:65282"},
				{Name: "large-test", Value: "large:123:456:789"},
			},
		},
	}

	if err := src.ConvertTo(&res); err != nil {
		t.Fatalf("ConvertTo failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected Community different than converted")
	}
}

func TestCommunityConvertFrom(t *testing.T) {
	var res Community

	src := metallbv1.Community{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "community1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.CommunitySpec{
			Communities: []metallbv1.CommunityAlias{
				{Name: "no-export", Value: "65535:65281"},
			},
		},
	}

	expected := Community{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "community1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: CommunitySpec{
			Communities: []CommunityAlias{
				{Name: "no-export", Value: "65535:65281"},
			},
		},
	}

	if err := res.ConvertFrom(&src); err != nil {
		t.Fatalf("ConvertFrom failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected Community different than converted")
	}
}
