// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_forwarding_rule

type PrivateDNSResolverForwardingRuleModel struct {
	Name                   string                 `tfschema:"name"`
	DnsForwardingRulesetId string                 `tfschema:"dns_forwarding_ruleset_id"`
	DomainName             string                 `tfschema:"domain_name"`
	Enabled                bool                   `tfschema:"enabled"`
	Metadata               map[string]string      `tfschema:"metadata"`
	TargetDnsServers       []TargetDnsServerModel `tfschema:"target_dns_servers"`
}

type TargetDnsServerModel struct {
	IPAddress string `tfschema:"ip_address"`
	Port      int64  `tfschema:"port"`
}

func (r Resource) ModelObject() any {
	return &PrivateDNSResolverForwardingRuleModel{}
}
