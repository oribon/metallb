// SPDX-License-Identifier:Apache-2.0

package v1beta2

import (
	"reflect"
	"testing"
	"time"

	metallbv1 "go.universe.tf/metallb/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	MetalLBTestNameSpace = "metallb-test-namespace"
)

func TestBGPPeerConvertTo(t *testing.T) {
	var resBGPPeer metallbv1.BGPPeer

	convertBGPPeer := BGPPeer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "peer1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: BGPPeerSpec{
			MyASN:        42,
			ASN:          142,
			DynamicASN:   ExternalASNMode,
			Address:      "1.2.3.4",
			Interface:    "eth0",
			Port:         1179,
			HoldTime:     &metav1.Duration{Duration: 180 * time.Second},
			RouterID:     "10.20.30.40",
			SrcAddress:   "10.20.30.40",
			EBGPMultiHop: true,
			Password:     "nopass",
			PasswordSecret: corev1.SecretReference{
				Name:      "secret1",
				Namespace: "metallb-system",
			},
			BFDProfile:             "default",
			KeepaliveTime:          &metav1.Duration{Duration: time.Second},
			ConnectTime:            &metav1.Duration{Duration: 10 * time.Second},
			EnableGracefulRestart:  true,
			VRFName:                "red",
			DualStackAddressFamily: true,
			NodeSelectors: []metav1.LabelSelector{
				{
					MatchLabels: map[string]string{"foo": "bar"},
				},
			},
		},
	}

	expectedBGPPeer := metallbv1.BGPPeer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "peer1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.BGPPeerSpec{
			MyASN:      42,
			ASN:        142,
			DynamicASN: metallbv1.ExternalASNMode,
			Address:    "1.2.3.4",
			Interface:  "eth0",
			Port:       1179,
			HoldTime:   &metav1.Duration{Duration: 180 * time.Second},
			RouterID:   "10.20.30.40",
			SrcAddress: "10.20.30.40",
			EBGPMultiHop: true,
			Password:     "nopass",
			PasswordSecret: corev1.SecretReference{
				Name:      "secret1",
				Namespace: "metallb-system",
			},
			BFDProfile:             "default",
			KeepaliveTime:          &metav1.Duration{Duration: time.Second},
			ConnectTime:            &metav1.Duration{Duration: 10 * time.Second},
			EnableGracefulRestart:  true,
			VRFName:                "red",
			DualStackAddressFamily: true,
			NodeSelectors: []metav1.LabelSelector{
				{
					MatchLabels: map[string]string{"foo": "bar"},
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
		ObjectMeta: metav1.ObjectMeta{
			Name:      "peer1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: metallbv1.BGPPeerSpec{
			MyASN:      42,
			ASN:        142,
			DynamicASN: metallbv1.InternalASNMode,
			Address:    "1.2.3.4",
			Interface:  "eth0",
			Port:       1179,
			HoldTime:   &metav1.Duration{Duration: 180 * time.Second},
			RouterID:   "10.20.30.40",
			SrcAddress: "10.20.30.40",
			EBGPMultiHop: true,
			Password:     "nopass",
			PasswordSecret: corev1.SecretReference{
				Name:      "secret1",
				Namespace: "metallb-system",
			},
			BFDProfile:             "default",
			KeepaliveTime:          &metav1.Duration{Duration: time.Second},
			ConnectTime:            &metav1.Duration{Duration: 10 * time.Second},
			EnableGracefulRestart:  true,
			VRFName:                "red",
			DualStackAddressFamily: true,
			NodeSelectors: []metav1.LabelSelector{
				{
					MatchLabels: map[string]string{"foo": "bar"},
				},
			},
		},
	}

	expectedBGPPeer := BGPPeer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "peer1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: BGPPeerSpec{
			MyASN:      42,
			ASN:        142,
			DynamicASN: InternalASNMode,
			Address:    "1.2.3.4",
			Interface:  "eth0",
			Port:       1179,
			HoldTime:   &metav1.Duration{Duration: 180 * time.Second},
			RouterID:   "10.20.30.40",
			SrcAddress: "10.20.30.40",
			EBGPMultiHop: true,
			Password:     "nopass",
			PasswordSecret: corev1.SecretReference{
				Name:      "secret1",
				Namespace: "metallb-system",
			},
			BFDProfile:             "default",
			KeepaliveTime:          &metav1.Duration{Duration: time.Second},
			ConnectTime:            &metav1.Duration{Duration: 10 * time.Second},
			EnableGracefulRestart:  true,
			VRFName:                "red",
			DualStackAddressFamily: true,
			NodeSelectors: []metav1.LabelSelector{
				{
					MatchLabels: map[string]string{"foo": "bar"},
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

func TestBGPPeerDisableMPDropped(t *testing.T) {
	src := BGPPeer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "peer1",
			Namespace: MetalLBTestNameSpace,
		},
		Spec: BGPPeerSpec{
			MyASN:     42,
			ASN:       142,
			Address:   "1.2.3.4",
			DisableMP: true,
		},
	}

	var hub metallbv1.BGPPeer
	if err := src.ConvertTo(&hub); err != nil {
		t.Fatalf("ConvertTo failed: %s", err)
	}

	if len(hub.Annotations) != 0 {
		t.Fatalf("expected no annotations on hub, got: %v", hub.Annotations)
	}

	var roundTripped BGPPeer
	if err := roundTripped.ConvertFrom(&hub); err != nil {
		t.Fatalf("ConvertFrom failed: %s", err)
	}

	if roundTripped.Spec.DisableMP {
		t.Fatal("expected DisableMP to be false after round-trip (field is dropped)")
	}
}
