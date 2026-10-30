// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_outbound_endpoint

type PrivateDNSResolverOutboundEndpointModel struct {
	Name                 string            `tfschema:"name"`
	PrivateDNSResolverId string            `tfschema:"private_dns_resolver_id"`
	Location             string            `tfschema:"location"`
	SubnetId             string            `tfschema:"subnet_id"`
	Tags                 map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &PrivateDNSResolverOutboundEndpointModel{}
}
