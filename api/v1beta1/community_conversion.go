// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	metallbv1 "go.universe.tf/metallb/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this Community to the Hub version (v1).
func (src *Community) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metallbv1.Community)
	dst.ObjectMeta = src.ObjectMeta

	if src.Spec.Communities != nil {
		dst.Spec.Communities = make([]metallbv1.CommunityAlias, len(src.Spec.Communities))
		for i, c := range src.Spec.Communities {
			dst.Spec.Communities[i] = metallbv1.CommunityAlias{
				Name:  c.Name,
				Value: c.Value,
			}
		}
	}

	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *Community) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metallbv1.Community)
	dst.ObjectMeta = src.ObjectMeta

	if src.Spec.Communities != nil {
		dst.Spec.Communities = make([]CommunityAlias, len(src.Spec.Communities))
		for i, c := range src.Spec.Communities {
			dst.Spec.Communities[i] = CommunityAlias{
				Name:  c.Name,
				Value: c.Value,
			}
		}
	}

	return nil
}
