// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	"reflect"
	"testing"

	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestServiceBGPStatusConvertTo(t *testing.T) {
	var res metallbv1.ServiceBGPStatus

	src := ServiceBGPStatus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "status1",
			Namespace: MetalLBTestNameSpace,
		},
		Status: MetalLBServiceBGPStatus{
			Node:             "node1",
			ServiceName:      "svc1",
			ServiceNamespace: "default",
			Peers:            []string{"peer1", "peer2"},
		},
	}

	expected := metallbv1.ServiceBGPStatus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "status1",
			Namespace: MetalLBTestNameSpace,
		},
		Status: metallbv1.MetalLBServiceBGPStatus{
			Node:             "node1",
			ServiceName:      "svc1",
			ServiceNamespace: "default",
			Peers:            []string{"peer1", "peer2"},
		},
	}

	if err := src.ConvertTo(&res); err != nil {
		t.Fatalf("ConvertTo failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected ServiceBGPStatus different than converted")
	}
}

func TestServiceBGPStatusConvertFrom(t *testing.T) {
	var res ServiceBGPStatus

	src := metallbv1.ServiceBGPStatus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "status1",
			Namespace: MetalLBTestNameSpace,
		},
		Status: metallbv1.MetalLBServiceBGPStatus{
			Node:             "node2",
			ServiceName:      "svc2",
			ServiceNamespace: "kube-system",
			Peers:            []string{"peer3"},
		},
	}

	expected := ServiceBGPStatus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "status1",
			Namespace: MetalLBTestNameSpace,
		},
		Status: MetalLBServiceBGPStatus{
			Node:             "node2",
			ServiceName:      "svc2",
			ServiceNamespace: "kube-system",
			Peers:            []string{"peer3"},
		},
	}

	if err := res.ConvertFrom(&src); err != nil {
		t.Fatalf("ConvertFrom failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected ServiceBGPStatus different than converted")
	}
}
