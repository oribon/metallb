// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this ConfigurationState to the Hub version (v1).
func (src *ConfigurationState) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metallbv1.ConfigurationState)
	dst.ObjectMeta = src.ObjectMeta

	dst.Status.Result = metallbv1.ConfigurationResult(src.Status.Result)
	dst.Status.ErrorSummary = src.Status.ErrorSummary
	dst.Status.Conditions = copyConditions(src.Status.Conditions)

	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *ConfigurationState) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metallbv1.ConfigurationState)
	dst.ObjectMeta = src.ObjectMeta

	dst.Status.Result = ConfigurationResult(src.Status.Result)
	dst.Status.ErrorSummary = src.Status.ErrorSummary
	dst.Status.Conditions = copyConditions(src.Status.Conditions)

	return nil
}

func copyConditions(conditions []metav1.Condition) []metav1.Condition {
	if conditions == nil {
		return nil
	}
	res := make([]metav1.Condition, len(conditions))
	copy(res, conditions)
	return res
}
