// SPDX-License-Identifier:Apache-2.0

package e2e

// Proto holds the protocol we are speaking.
type Proto string

// MetalLB supported protocols.
const (
	BGP    Proto = "bgp"
	Layer2 Proto = "layer2"
)

// configFile is the configuration as parsed out of the ConfigMap,
// without validation or useful high level types.
type configFile struct {
	Peers          []peer            `yaml:"peers,omitempty"`
	BGPCommunities map[string]string `yaml:"bgp-communities,omitempty"`
	Pools          []addressPool     `yaml:"address-pools,omitempty"`
}

type peer struct {
	MyASN         uint32         `yaml:"my-asn,omitempty"`
	ASN           uint32         `yaml:"peer-asn,omitempty"`
	Addr          string         `yaml:"peer-address,omitempty"`
	SrcAddr       string         `yaml:"source-address,omitempty"`
	Port          uint16         `yaml:"peer-port,omitempty"`
	HoldTime      string         `yaml:"hold-time,omitempty"`
	RouterID      string         `yaml:"router-id,omitempty"`
	NodeSelectors []nodeSelector `yaml:"node-selectors,omitempty"`
	Password      string         `yaml:"password,omitempty"`
}

type nodeSelector struct {
	MatchLabels      map[string]string      `yaml:"match-labels,omitempty"`
	MatchExpressions []selectorRequirements `yaml:"match-expressions,omitempty"`
}

type selectorRequirements struct {
	Key      string   `yaml:"key,omitempty"`
	Operator string   `yaml:"operator,omitempty"`
	Values   []string `yaml:"values,omitempty"`
}

type addressPool struct {
	Protocol          Proto              `yaml:"protocol,omitempty"`
	Name              string             `yaml:"name,omitempty"`
	Addresses         []string           `yaml:"addresses,omitempty"`
	AvoidBuggyIPs     bool               `yaml:"avoid-buggy-ips,omitempty"`
	AutoAssign        *bool              `yaml:"auto-assign,omitempty"`
	BGPAdvertisements []bgpAdvertisement `yaml:"bgp-advertisements,omitempty"`
}

type bgpAdvertisement struct {
	AggregationLength *int     `yaml:"aggregation-length,omitempty"`
	LocalPref         *uint32  `yaml:"localpref,omitempty"`
	Communities       []string `yaml:"communities,omitempty"`
}
