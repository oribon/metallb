// SPDX-License-Identifier:Apache-2.0

package config

import (
	"testing"
	"time"

	metallbv1 "go.universe.tf/metallb/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		config   ClusterResources
		mustFail bool
	}{
		{
			desc: "peer with bfd profile",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address:    "1.2.3.4",
							BFDProfile: "foo",
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "bfd profile set",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
						},
					},
				},
				BFDProfiles: []metallbv1.BFDProfile{
					{
						ObjectMeta: v1.ObjectMeta{Name: "foo"},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "v6 address",
			config: ClusterResources{
				Pools: []metallbv1.IPAddressPool{
					{
						Spec: metallbv1.IPAddressPoolSpec{
							Addresses: []string{"2001:db8::/64"},
						},
					},
				},
				BGPAdvs: []metallbv1.BGPAdvertisement{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "foo",
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "v6 address but pool not selected",
			config: ClusterResources{
				Pools: []metallbv1.IPAddressPool{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "foo",
						},
						Spec: metallbv1.IPAddressPoolSpec{
							Addresses: []string{"2001:db8::/64"},
						},
					},
				},
				BGPAdvs: []metallbv1.BGPAdvertisement{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "foo",
						},
						Spec: metallbv1.BGPAdvertisementSpec{
							IPAddressPools: []string{"bar"},
						},
					},
				},
			},
		},
		{
			desc: "v6 address and selected",
			config: ClusterResources{
				Pools: []metallbv1.IPAddressPool{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "foo",
						},
						Spec: metallbv1.IPAddressPoolSpec{
							Addresses: []string{"2001:db8::/64"},
						},
					},
				},
				BGPAdvs: []metallbv1.BGPAdvertisement{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "bar",
						},
						Spec: metallbv1.BGPAdvertisementSpec{
							IPAddressPools: []string{"foo1"},
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "bar",
						},
						Spec: metallbv1.BGPAdvertisementSpec{
							IPAddressPools: []string{"foo"},
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "v6 address and selected by labels",
			config: ClusterResources{
				Pools: []metallbv1.IPAddressPool{
					{
						ObjectMeta: v1.ObjectMeta{
							Name:   "foo",
							Labels: map[string]string{"key": "value"},
						},
						Spec: metallbv1.IPAddressPoolSpec{
							Addresses: []string{"2001:db8::/64"},
						},
					},
				},
				BGPAdvs: []metallbv1.BGPAdvertisement{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "bar",
						},
						Spec: metallbv1.BGPAdvertisementSpec{
							IPAddressPools: []string{"foo1"},
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "bar",
						},
						Spec: metallbv1.BGPAdvertisementSpec{
							IPAddressPoolSelectors: []v1.LabelSelector{
								{
									MatchLabels: map[string]string{"key": "value"},
								},
							},
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "enable BGP GracefulRestart",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							EnableGracefulRestart: true,
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "dynamic ASN",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							DynamicASN: metallbv1.InternalASNMode,
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "interface",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Interface: "eth0",
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "keepalive time",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							KeepaliveTime: &v1.Duration{Duration: time.Second},
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "connect time",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							ConnectTime: ptr.To(v1.Duration{Duration: time.Second}),
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "large BGP community inside BGP Advertisement",
			config: ClusterResources{
				BGPAdvs: []metallbv1.BGPAdvertisement{
					{
						Spec: metallbv1.BGPAdvertisementSpec{
							Communities: []string{"large:123:456:789"},
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "large BGP community inside Community CR",
			config: ClusterResources{
				Communities: []metallbv1.Community{
					{
						Spec: metallbv1.CommunitySpec{
							Communities: []metallbv1.CommunityAlias{
								{
									Value: "large:123:456:789",
								},
							},
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "should pass",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			err := DiscardFRROnly(test.config)
			if test.mustFail && err == nil {
				t.Fatalf("Expected error for %s", test.desc)
			}
			if !test.mustFail && err != nil {
				t.Fatalf("Not expected error %s for %s", err, test.desc)
			}
		})
	}
}

func TestValidateFRR(t *testing.T) {
	tests := []struct {
		desc     string
		config   ClusterResources
		mustFail bool
	}{
		{
			desc: "peer with routerid",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address:  "1.2.3.4",
							RouterID: "1.2.3.4",
						},
					},
				},
			},
		},
		{
			desc: "routerid set, one different",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address:  "1.2.3.4",
							RouterID: "1.2.3.4",
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address:  "1.2.3.5",
							RouterID: "1.2.3.4",
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address:  "1.2.3.6",
							RouterID: "1.2.3.5",
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "routerid set, one not set",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address:  "1.2.3.4",
							RouterID: "1.2.3.4",
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address:  "1.2.3.5",
							RouterID: "1.2.3.4",
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.6",
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "bfd profile set",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address:  "1.2.3.4",
							RouterID: "1.2.3.4",
						},
					},
				},
				BFDProfiles: []metallbv1.BFDProfile{
					{
						ObjectMeta: v1.ObjectMeta{
							Name: "foo",
						},
					},
				},
			},
		},
		{
			desc: "myAsn set, all equals",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
							MyASN:   123,
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.5",
							MyASN:   123,
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.6",
							MyASN:   123,
						},
					},
				},
			},
		},
		{
			desc: "myAsn set, one different",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
							MyASN:   123,
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.5",
							MyASN:   123,
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.6",
							MyASN:   124,
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "myAsn set, one different but with different vrf",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
							MyASN:   123,
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.5",
							MyASN:   123,
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.6",
							MyASN:   124,
							VRFName: "red",
						},
					},
				},
			},
		},
		{
			desc: "myAsn set, two different but with different vrf",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
							MyASN:   123,
							VRFName: "red",
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.5",
							MyASN:   123,
							VRFName: "red",
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.6",
							MyASN:   124,
						},
					},
				},
			},
		},
		{
			desc: "duplicate bgp address",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
						},
					},
				},
			},
			mustFail: true,
		}, {
			desc: "duplicate bgp address, different vrfs",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
							VRFName: "red",
						},
					}, {
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
							VRFName: "green",
						},
					},
				},
			},
			mustFail: false,
		}, {
			desc: "duplicate bgp address, same vrf",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
							VRFName: "red",
						},
					}, {
						Spec: metallbv1.BGPPeerSpec{
							Address: "1.2.3.4",
							VRFName: "red",
						},
					},
				},
			},
			mustFail: true,
		},
		{
			desc: "two peers with interface set different",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Interface: "eth0",
							MyASN:     123,
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Interface: "eth1",
							MyASN:     123,
						},
					},
				},
			},
			mustFail: false,
		},
		{
			desc: "two peers with interface set same",
			config: ClusterResources{
				Peers: []metallbv1.BGPPeer{
					{
						Spec: metallbv1.BGPPeerSpec{
							Interface: "eth0",
							MyASN:     123,
						},
					},
					{
						Spec: metallbv1.BGPPeerSpec{
							Interface: "eth0",
							MyASN:     123,
						},
					},
				},
			},
			mustFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			err := DiscardNativeOnly(test.config)
			if test.mustFail && err == nil {
				t.Fatalf("Expected error for %s", test.desc)
			}
			if !test.mustFail && err != nil {
				t.Fatalf("Not expected error %s for %s", err, test.desc)
			}
		})
	}
}
