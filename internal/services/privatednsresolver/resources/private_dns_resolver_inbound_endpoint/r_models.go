// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_resolver_inbound_endpoint

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/inboundendpoints"
)

type PrivateDNSResolverInboundEndpointModel struct {
	Name                 string                 `tfschema:"name"`
	PrivateDNSResolverId string                 `tfschema:"private_dns_resolver_id"`
	IPConfigurations     []IPConfigurationModel `tfschema:"ip_configurations"`
	Location             string                 `tfschema:"location"`
	Tags                 map[string]string      `tfschema:"tags"`
}

type IPConfigurationModel struct {
	PrivateIPAddress          string                              `tfschema:"private_ip_address"`
	PrivateIPAllocationMethod inboundendpoints.IPAllocationMethod `tfschema:"private_ip_allocation_method"`
	SubnetId                  string                              `tfschema:"subnet_id"`
}

func (r Resource) ModelObject() any {
	return &PrivateDNSResolverInboundEndpointModel{}
}
