// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	metallbv1 "go.universe.tf/metallb/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this ServiceL2Status to the Hub version (v1).
func (src *ServiceL2Status) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metallbv1.ServiceL2Status)
	dst.ObjectMeta = src.ObjectMeta

	if src.Status.Interfaces != nil {
		dst.Status.Interfaces = make([]metallbv1.InterfaceInfo, len(src.Status.Interfaces))
		for i, iface := range src.Status.Interfaces {
			dst.Status.Interfaces[i] = metallbv1.InterfaceInfo{
				Name: iface.Name,
			}
		}
	}
	dst.Status.Node = src.Status.Node
	dst.Status.ServiceName = src.Status.ServiceName
	dst.Status.ServiceNamespace = src.Status.ServiceNamespace

	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *ServiceL2Status) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metallbv1.ServiceL2Status)
	dst.ObjectMeta = src.ObjectMeta

	if src.Status.Interfaces != nil {
		dst.Status.Interfaces = make([]InterfaceInfo, len(src.Status.Interfaces))
		for i, iface := range src.Status.Interfaces {
			dst.Status.Interfaces[i] = InterfaceInfo{
				Name: iface.Name,
			}
		}
	}
	dst.Status.Node = src.Status.Node
	dst.Status.ServiceName = src.Status.ServiceName
	dst.Status.ServiceNamespace = src.Status.ServiceNamespace

	return nil
}
