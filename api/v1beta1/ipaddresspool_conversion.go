// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	metallbv1 "go.universe.tf/metallb/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this IPAddressPool to the Hub version (v1).
func (src *IPAddressPool) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metallbv1.IPAddressPool)
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Addresses = src.Spec.Addresses
	dst.Spec.AutoAssign = src.Spec.AutoAssign
	dst.Spec.AvoidBuggyIPs = src.Spec.AvoidBuggyIPs
	if src.Spec.AllocateTo != nil {
		dst.Spec.AllocateTo = &metallbv1.ServiceAllocation{
			Priority:           src.Spec.AllocateTo.Priority,
			Namespaces:         src.Spec.AllocateTo.Namespaces,
			NamespaceSelectors: copyLabelSelectors(src.Spec.AllocateTo.NamespaceSelectors),
			ServiceSelectors:   copyLabelSelectors(src.Spec.AllocateTo.ServiceSelectors),
		}
	}

	dst.Status.AssignedIPv4 = src.Status.AssignedIPv4
	dst.Status.AssignedIPv6 = src.Status.AssignedIPv6
	dst.Status.AvailableIPv4 = src.Status.AvailableIPv4
	dst.Status.AvailableIPv6 = src.Status.AvailableIPv6

	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *IPAddressPool) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metallbv1.IPAddressPool)
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Addresses = src.Spec.Addresses
	dst.Spec.AutoAssign = src.Spec.AutoAssign
	dst.Spec.AvoidBuggyIPs = src.Spec.AvoidBuggyIPs
	if src.Spec.AllocateTo != nil {
		dst.Spec.AllocateTo = &ServiceAllocation{
			Priority:           src.Spec.AllocateTo.Priority,
			Namespaces:         src.Spec.AllocateTo.Namespaces,
			NamespaceSelectors: copyLabelSelectors(src.Spec.AllocateTo.NamespaceSelectors),
			ServiceSelectors:   copyLabelSelectors(src.Spec.AllocateTo.ServiceSelectors),
		}
	}

	dst.Status.AssignedIPv4 = src.Status.AssignedIPv4
	dst.Status.AssignedIPv6 = src.Status.AssignedIPv6
	dst.Status.AvailableIPv4 = src.Status.AvailableIPv4
	dst.Status.AvailableIPv6 = src.Status.AvailableIPv6

	return nil
}
