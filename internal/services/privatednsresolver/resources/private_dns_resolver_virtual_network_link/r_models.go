// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_virtual_network_link

type PrivateDNSResolverVirtualNetworkLinkModel struct {
	Name                   string            `tfschema:"name"`
	DnsForwardingRulesetId string            `tfschema:"dns_forwarding_ruleset_id"`
	Metadata               map[string]string `tfschema:"metadata"`
	VirtualNetworkId       string            `tfschema:"virtual_network_id"`
}

func (r Resource) ModelObject() any {
	return &PrivateDNSResolverVirtualNetworkLinkModel{}
}
