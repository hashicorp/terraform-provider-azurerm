// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_network_anchor

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/resourceanchors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 24),
				validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]+$`), "Name may include letters, numbers, or hyphens only"),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"resource_anchor_id": commonschema.ResourceIDReferenceRequiredForceNew(&resourceanchors.ResourceAnchorId{}),

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"zones": commonschema.ZonesMultipleRequiredForceNew(),

		"backup_subnet_cidr": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsCIDR,
		},

		"dns_forwarding_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"domain_names": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						ForceNew: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.DomainName,
						},
					},
					"forwarding_ip_address": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsIPAddress,
					},
				},
			},
		},

		"dns_listening_endpoint_allowed_cidrs": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsCIDR,
			},
		},

		"dns_private_zone_replication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dns_forwarding_endpoint_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_forwarding_endpoint_nsg_rule_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_forwarding_rule_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_listening_endpoint_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"dns_listening_endpoint_nsg_rule_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"oci_vcn_dns_label": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
