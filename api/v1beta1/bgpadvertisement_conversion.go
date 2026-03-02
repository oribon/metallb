// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	metallbv1 "go.universe.tf/metallb/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this BGPAdvertisement to the Hub version (v1).
func (src *BGPAdvertisement) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metallbv1.BGPAdvertisement)
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.AggregationLength = src.Spec.AggregationLength
	dst.Spec.AggregationLengthV6 = src.Spec.AggregationLengthV6
	dst.Spec.LocalPref = src.Spec.LocalPref
	dst.Spec.Communities = src.Spec.Communities
	dst.Spec.IPAddressPools = src.Spec.IPAddressPools
	dst.Spec.IPAddressPoolSelectors = copyLabelSelectors(src.Spec.IPAddressPoolSelectors)
	dst.Spec.NodeSelectors = copyLabelSelectors(src.Spec.NodeSelectors)
	dst.Spec.Peers = src.Spec.Peers
	dst.Spec.ServiceSelectors = copyLabelSelectors(src.Spec.ServiceSelectors)

	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *BGPAdvertisement) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metallbv1.BGPAdvertisement)
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.AggregationLength = src.Spec.AggregationLength
	dst.Spec.AggregationLengthV6 = src.Spec.AggregationLengthV6
	dst.Spec.LocalPref = src.Spec.LocalPref
	dst.Spec.Communities = src.Spec.Communities
	dst.Spec.IPAddressPools = src.Spec.IPAddressPools
	dst.Spec.IPAddressPoolSelectors = copyLabelSelectors(src.Spec.IPAddressPoolSelectors)
	dst.Spec.NodeSelectors = copyLabelSelectors(src.Spec.NodeSelectors)
	dst.Spec.Peers = src.Spec.Peers
	dst.Spec.ServiceSelectors = copyLabelSelectors(src.Spec.ServiceSelectors)

	return nil
}

func copyLabelSelectors(selectors []metav1.LabelSelector) []metav1.LabelSelector {
	if selectors == nil {
		return nil
	}
	res := make([]metav1.LabelSelector, len(selectors))
	for i, sel := range selectors {
		res[i] = *sel.DeepCopy()
	}
	return res
}
