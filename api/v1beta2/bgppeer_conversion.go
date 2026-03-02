// SPDX-License-Identifier:Apache-2.0

package v1beta2

import (
	metallbv1 "go.universe.tf/metallb/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this BGPPeer to the Hub version (v1).
func (src *BGPPeer) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*metallbv1.BGPPeer)
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.MyASN = src.Spec.MyASN
	dst.Spec.ASN = src.Spec.ASN
	dst.Spec.DynamicASN = metallbv1.DynamicASNMode(src.Spec.DynamicASN)
	dst.Spec.Address = src.Spec.Address
	dst.Spec.Interface = src.Spec.Interface
	dst.Spec.SrcAddress = src.Spec.SrcAddress
	dst.Spec.Port = src.Spec.Port
	dst.Spec.HoldTime = src.Spec.HoldTime
	dst.Spec.KeepaliveTime = src.Spec.KeepaliveTime
	dst.Spec.ConnectTime = src.Spec.ConnectTime
	dst.Spec.RouterID = src.Spec.RouterID
	dst.Spec.NodeSelectors = src.Spec.NodeSelectors
	dst.Spec.Password = src.Spec.Password
	dst.Spec.PasswordSecret = src.Spec.PasswordSecret
	dst.Spec.BFDProfile = src.Spec.BFDProfile
	dst.Spec.EnableGracefulRestart = src.Spec.EnableGracefulRestart
	dst.Spec.EBGPMultiHop = src.Spec.EBGPMultiHop
	dst.Spec.VRFName = src.Spec.VRFName
	dst.Spec.DualStackAddressFamily = src.Spec.DualStackAddressFamily
	// DisableMP is deprecated and a no-op; intentionally not converted.

	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *BGPPeer) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*metallbv1.BGPPeer)
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.MyASN = src.Spec.MyASN
	dst.Spec.ASN = src.Spec.ASN
	dst.Spec.DynamicASN = DynamicASNMode(src.Spec.DynamicASN)
	dst.Spec.Address = src.Spec.Address
	dst.Spec.Interface = src.Spec.Interface
	dst.Spec.SrcAddress = src.Spec.SrcAddress
	dst.Spec.Port = src.Spec.Port
	dst.Spec.HoldTime = src.Spec.HoldTime
	dst.Spec.KeepaliveTime = src.Spec.KeepaliveTime
	dst.Spec.ConnectTime = src.Spec.ConnectTime
	dst.Spec.RouterID = src.Spec.RouterID
	dst.Spec.NodeSelectors = src.Spec.NodeSelectors
	dst.Spec.Password = src.Spec.Password
	dst.Spec.PasswordSecret = src.Spec.PasswordSecret
	dst.Spec.BFDProfile = src.Spec.BFDProfile
	dst.Spec.EnableGracefulRestart = src.Spec.EnableGracefulRestart
	dst.Spec.EBGPMultiHop = src.Spec.EBGPMultiHop
	dst.Spec.VRFName = src.Spec.VRFName
	dst.Spec.DualStackAddressFamily = src.Spec.DualStackAddressFamily
	// DisableMP is deprecated and a no-op; defaults to false.

	return nil
}
