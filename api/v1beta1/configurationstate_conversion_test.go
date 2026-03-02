// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	"reflect"
	"testing"

	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConfigurationStateConvertTo(t *testing.T) {
	var res metallbv1.ConfigurationState

	src := ConfigurationState{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "controller",
			Namespace: MetalLBTestNameSpace,
		},
		Status: ConfigurationStateStatus{
			Result:       ConfigurationResultValid,
			ErrorSummary: "",
			Conditions: []metav1.Condition{
				{
					Type:   "ConfigLoaded",
					Status: metav1.ConditionTrue,
					Reason: "Valid",
				},
			},
		},
	}

	expected := metallbv1.ConfigurationState{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "controller",
			Namespace: MetalLBTestNameSpace,
		},
		Status: metallbv1.ConfigurationStateStatus{
			Result:       metallbv1.ConfigurationResultValid,
			ErrorSummary: "",
			Conditions: []metav1.Condition{
				{
					Type:   "ConfigLoaded",
					Status: metav1.ConditionTrue,
					Reason: "Valid",
				},
			},
		},
	}

	if err := src.ConvertTo(&res); err != nil {
		t.Fatalf("ConvertTo failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected ConfigurationState different than converted")
	}
}

func TestConfigurationStateConvertFrom(t *testing.T) {
	var res ConfigurationState

	src := metallbv1.ConfigurationState{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "speaker-node1",
			Namespace: MetalLBTestNameSpace,
		},
		Status: metallbv1.ConfigurationStateStatus{
			Result:       metallbv1.ConfigurationResultInvalid,
			ErrorSummary: "BGPPeer invalid",
		},
	}

	expected := ConfigurationState{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "speaker-node1",
			Namespace: MetalLBTestNameSpace,
		},
		Status: ConfigurationStateStatus{
			Result:       ConfigurationResultInvalid,
			ErrorSummary: "BGPPeer invalid",
		},
	}

	if err := res.ConvertFrom(&src); err != nil {
		t.Fatalf("ConvertFrom failed: %s", err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Fatalf("expected ConfigurationState different than converted")
	}
}
