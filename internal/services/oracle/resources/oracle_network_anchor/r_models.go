// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_network_anchor

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

type NetworkAnchorResourceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	Location          string `tfschema:"location"`

	ResourceAnchorId string       `tfschema:"resource_anchor_id"`
	SubnetId         string       `tfschema:"subnet_id"`
	Zones            zones.Schema `tfschema:"zones"`

	BackupSubnetCidr                 string                   `tfschema:"backup_subnet_cidr"`
	DnsForwardingRule                []DnsForwardingRuleModel `tfschema:"dns_forwarding_rule"`
	DnsListeningEndpointAllowedCidrs []string                 `tfschema:"dns_listening_endpoint_allowed_cidrs"`
	DnsPrivateZoneReplicationEnabled bool                     `tfschema:"dns_private_zone_replication_enabled"`

	DnsForwardingEndpointIpAddress  string            `tfschema:"dns_forwarding_endpoint_ip_address"`
	DnsForwardingEndpointNsgRuleUrl string            `tfschema:"dns_forwarding_endpoint_nsg_rule_url"`
	DnsForwardingRuleUrl            string            `tfschema:"dns_forwarding_rule_url"`
	DnsListeningEndpointIpAddress   string            `tfschema:"dns_listening_endpoint_ip_address"`
	DnsListeningEndpointNsgRuleUrl  string            `tfschema:"dns_listening_endpoint_nsg_rule_url"`
	OciVcnDnsLabel                  string            `tfschema:"oci_vcn_dns_label"`
	Tags                            map[string]string `tfschema:"tags"`
}

type DnsForwardingRuleModel struct {
	DomainNames         []string `tfschema:"domain_names"`
	ForwardingIPAddress string   `tfschema:"forwarding_ip_address"`
}

func (Resource) ModelObject() any {
	return &NetworkAnchorResourceModel{}
}
