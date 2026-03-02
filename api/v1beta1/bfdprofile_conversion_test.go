// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	"reflect"
	"testing"

	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBFDProfileConvertTo(t *testing.T) {
	var res metallbv1.BFDProfile
	recvInterval := uint32(300)
	transmitInterval := uint32(300)
	detectMultiplier := uint32(3)
	echoInterval := uint32(50)
	echoMode := false
	passiveMode := true
	minimumTTL := uint32(254)

	src := BFDProfile{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bfd1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: BFDProfileSpec{
			ReceiveInterval:  &recvInterval,
			TransmitInterval: &transmitInterval,
			DetectMultiplier: &detectMultiplier,
			EchoInterval:     &echoInterval,
			EchoMode:         &echoMode,
			PassiveMode:      &passiveMode,
			MinimumTTL:       &minimumTTL,
		},
	}

	expected := metallbv1.BFDProfile{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bfd1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.BFDProfileSpec{
			ReceiveInterval:  &recvInterval,
			TransmitInterval: &transmitInterval,
			DetectMultiplier: &detectMultiplier,
			EchoInterval:     &echoInterval,
			EchoMode:         &echoMode,
			PassiveMode:      &passiveMode,
			MinimumTTL:       &minimumTTL,
		},
	}

	if err := src.ConvertTo(&res); err != nil {
		t.Fatalf("ConvertTo failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected BFDProfile different than converted")
	}
}

func TestBFDProfileConvertFrom(t *testing.T) {
	var res BFDProfile
	recvInterval := uint32(500)
	passiveMode := false

	src := metallbv1.BFDProfile{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bfd1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.BFDProfileSpec{
			ReceiveInterval: &recvInterval,
			PassiveMode:     &passiveMode,
		},
	}

	expected := BFDProfile{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bfd1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: BFDProfileSpec{
			ReceiveInterval: &recvInterval,
			PassiveMode:     &passiveMode,
		},
	}

	if err := res.ConvertFrom(&src); err != nil {
		t.Fatalf("ConvertFrom failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected BFDProfile different than converted")
	}
}
