// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_dns_forwarding_ruleset

type PrivateDNSResolverDnsForwardingRulesetModel struct {
	Name                         string            `tfschema:"name"`
	ResourceGroupName            string            `tfschema:"resource_group_name"`
	DnsResolverOutboundEndpoints []string          `tfschema:"private_dns_resolver_outbound_endpoint_ids"`
	Location                     string            `tfschema:"location"`
	Tags                         map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &PrivateDNSResolverDnsForwardingRulesetModel{}
}
