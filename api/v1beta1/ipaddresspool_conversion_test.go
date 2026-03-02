// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	"reflect"
	"testing"

	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIPAddressPoolConvertTo(t *testing.T) {
	var res metallbv1.IPAddressPool
	autoAssign := true

	src := IPAddressPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pool1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: IPAddressPoolSpec{
			Addresses:     []string{"192.168.1.0/24"},
			AutoAssign:    &autoAssign,
			AvoidBuggyIPs: true,
			AllocateTo: &ServiceAllocation{
				Priority:   10,
				Namespaces: []string{"default"},
				NamespaceSelectors: []metav1.LabelSelector{
					{MatchLabels: map[string]string{"env": "prod"}},
				},
				ServiceSelectors: []metav1.LabelSelector{
					{MatchLabels: map[string]string{"app": "web"}},
				},
			},
		},
		Status: IPAddressPoolStatus{
			AssignedIPv4:  5,
			AssignedIPv6:  2,
			AvailableIPv4: 249,
			AvailableIPv6: 100,
		},
	}

	expected := metallbv1.IPAddressPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pool1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.IPAddressPoolSpec{
			Addresses:     []string{"192.168.1.0/24"},
			AutoAssign:    &autoAssign,
			AvoidBuggyIPs: true,
			AllocateTo: &metallbv1.ServiceAllocation{
				Priority:   10,
				Namespaces: []string{"default"},
				NamespaceSelectors: []metav1.LabelSelector{
					{MatchLabels: map[string]string{"env": "prod"}},
				},
				ServiceSelectors: []metav1.LabelSelector{
					{MatchLabels: map[string]string{"app": "web"}},
				},
			},
		},
		Status: metallbv1.IPAddressPoolStatus{
			AssignedIPv4:  5,
			AssignedIPv6:  2,
			AvailableIPv4: 249,
			AvailableIPv6: 100,
		},
	}

	if err := src.ConvertTo(&res); err != nil {
		t.Fatalf("ConvertTo failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected IPAddressPool different than converted")
	}
}

func TestIPAddressPoolConvertFrom(t *testing.T) {
	var res IPAddressPool

	src := metallbv1.IPAddressPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pool1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.IPAddressPoolSpec{
			Addresses:     []string{"10.0.0.0/16"},
			AvoidBuggyIPs: false,
		},
		Status: metallbv1.IPAddressPoolStatus{
			AssignedIPv4:  10,
			AvailableIPv4: 65516,
		},
	}

	expected := IPAddressPool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pool1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: IPAddressPoolSpec{
			Addresses:     []string{"10.0.0.0/16"},
			AvoidBuggyIPs: false,
		},
		Status: IPAddressPoolStatus{
			AssignedIPv4:  10,
			AvailableIPv4: 65516,
		},
	}

	if err := res.ConvertFrom(&src); err != nil {
		t.Fatalf("ConvertFrom failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected IPAddressPool different than converted")
	}
}
