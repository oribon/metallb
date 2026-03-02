// SPDX-License-Identifier:Apache-2.0

package v1beta1

import (
	metallbv1 "go.universe.tf/metallb/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this BFDProfile to the Hub version (v1).
func (src *BFDProfile) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metallbv1.BFDProfile)
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.ReceiveInterval = src.Spec.ReceiveInterval
	dst.Spec.TransmitInterval = src.Spec.TransmitInterval
	dst.Spec.DetectMultiplier = src.Spec.DetectMultiplier
	dst.Spec.EchoInterval = src.Spec.EchoInterval
	dst.Spec.EchoMode = src.Spec.EchoMode
	dst.Spec.PassiveMode = src.Spec.PassiveMode
	dst.Spec.MinimumTTL = src.Spec.MinimumTTL

	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *BFDProfile) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metallbv1.BFDProfile)
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.ReceiveInterval = src.Spec.ReceiveInterval
	dst.Spec.TransmitInterval = src.Spec.TransmitInterval
	dst.Spec.DetectMultiplier = src.Spec.DetectMultiplier
	dst.Spec.EchoInterval = src.Spec.EchoInterval
	dst.Spec.EchoMode = src.Spec.EchoMode
	dst.Spec.PassiveMode = src.Spec.PassiveMode
	dst.Spec.MinimumTTL = src.Spec.MinimumTTL

	return nil
}
