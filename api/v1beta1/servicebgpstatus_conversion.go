// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	metallbv1 "go.universe.tf/metallb/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this ServiceBGPStatus to the Hub version (v1).
func (src *ServiceBGPStatus) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metallbv1.ServiceBGPStatus)
	dst.ObjectMeta = src.ObjectMeta

	dst.Status.Node = src.Status.Node
	dst.Status.ServiceName = src.Status.ServiceName
	dst.Status.ServiceNamespace = src.Status.ServiceNamespace
	dst.Status.Peers = src.Status.Peers

	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *ServiceBGPStatus) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metallbv1.ServiceBGPStatus)
	dst.ObjectMeta = src.ObjectMeta

	dst.Status.Node = src.Status.Node
	dst.Status.ServiceName = src.Status.ServiceName
	dst.Status.ServiceNamespace = src.Status.ServiceNamespace
	dst.Status.Peers = src.Status.Peers

	return nil
}
