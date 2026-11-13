// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver

type PrivateDNSResolverDnsResolverModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
	VirtualNetworkId  string            `tfschema:"virtual_network_id"`
}

func (r Resource) ModelObject() any {
	return &PrivateDNSResolverDnsResolverModel{}
}
