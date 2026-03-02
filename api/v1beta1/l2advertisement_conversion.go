// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	metallbv1 "go.universe.tf/metallb/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this L2Advertisement to the Hub version (v1).
func (src *L2Advertisement) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metallbv1.L2Advertisement)
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.IPAddressPools = src.Spec.IPAddressPools
	dst.Spec.IPAddressPoolSelectors = copyLabelSelectors(src.Spec.IPAddressPoolSelectors)
	dst.Spec.NodeSelectors = copyLabelSelectors(src.Spec.NodeSelectors)
	dst.Spec.Interfaces = src.Spec.Interfaces
	dst.Spec.ServiceSelectors = copyLabelSelectors(src.Spec.ServiceSelectors)

	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *L2Advertisement) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metallbv1.L2Advertisement)
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.IPAddressPools = src.Spec.IPAddressPools
	dst.Spec.IPAddressPoolSelectors = copyLabelSelectors(src.Spec.IPAddressPoolSelectors)
	dst.Spec.NodeSelectors = copyLabelSelectors(src.Spec.NodeSelectors)
	dst.Spec.Interfaces = src.Spec.Interfaces
	dst.Spec.ServiceSelectors = copyLabelSelectors(src.Spec.ServiceSelectors)

	return nil
}
