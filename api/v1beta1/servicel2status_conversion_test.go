// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	"reflect"
	"testing"

	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestServiceL2StatusConvertTo(t *testing.T) {
	var res metallbv1.ServiceL2Status

	src := ServiceL2Status{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "status1",
			Namespace: MetalLBTestNameSpace,
		},
		Status: MetalLBServiceL2Status{
			Node:             "node1",
			ServiceName:      "svc1",
			ServiceNamespace: "default",
			Interfaces: []InterfaceInfo{
				{Name: "eth0"},
				{Name: "eth1"},
			},
		},
	}

	expected := metallbv1.ServiceL2Status{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "status1",
			Namespace: MetalLBTestNameSpace,
		},
		Status: metallbv1.MetalLBServiceL2Status{
			Node:             "node1",
			ServiceName:      "svc1",
			ServiceNamespace: "default",
			Interfaces: []metallbv1.InterfaceInfo{
				{Name: "eth0"},
				{Name: "eth1"},
			},
		},
	}

	if err := src.ConvertTo(&res); err != nil {
		t.Fatalf("ConvertTo failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected ServiceL2Status different than converted")
	}
}

func TestServiceL2StatusConvertFrom(t *testing.T) {
	var res ServiceL2Status

	src := metallbv1.ServiceL2Status{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "status1",
			Namespace: MetalLBTestNameSpace,
		},
		Status: metallbv1.MetalLBServiceL2Status{
			Node:             "node2",
			ServiceName:      "svc2",
			ServiceNamespace: "kube-system",
			Interfaces: []metallbv1.InterfaceInfo{
				{Name: "br0"},
			},
		},
	}

	expected := ServiceL2Status{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "status1",
			Namespace: MetalLBTestNameSpace,
		},
		Status: MetalLBServiceL2Status{
			Node:             "node2",
			ServiceName:      "svc2",
			ServiceNamespace: "kube-system",
			Interfaces: []InterfaceInfo{
				{Name: "br0"},
			},
		},
	}

	if err := res.ConvertFrom(&src); err != nil {
		t.Fatalf("ConvertFrom failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected ServiceL2Status different than converted")
	}
}
