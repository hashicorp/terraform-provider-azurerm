// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_security_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func networkSecurityGroupSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Network Security Group.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"security_rule": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Set:        hashNetworkSecurityRule,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name of the security rule.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},

					"description": {
						Description:  "A description for this rule.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(0, 140),
					},

					"protocol": {
						Description:  "Network protocol this rule applies to.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(networksecuritygroups.PossibleValuesForSecurityRuleProtocol(), false),
					},

					"source_port_range": {
						Description: "Source Port or Range. Integer or range between `0` and `65535` or `*` to match any.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"source_port_ranges": {
						Description: "List of source ports or port ranges.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         pluginsdk.HashString,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},

					"destination_port_range": {
						Description: "Destination Port or Range. Integer or range between `0` and `65535` or `*` to match any.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"destination_port_ranges": {
						Description: "List of destination ports or port ranges.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         pluginsdk.HashString,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},

					"source_address_prefix": {
						Description: "CIDR or source IP range or * to match any IP. Tags such as `VirtualNetwork`, `AzureLoadBalancer` and `Internet` can also be used.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"source_address_prefixes": {
						Description: "List of source address prefixes. Tags may not be used.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         pluginsdk.HashString,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},

					"destination_address_prefix": {
						Description: "CIDR or destination IP range or * to match any IP. Tags such as `VirtualNetwork`, `AzureLoadBalancer` and `Internet` can also be used.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"destination_address_prefixes": {
						Description: "List of destination address prefixes. Tags may not be used.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         pluginsdk.HashString,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},

					"destination_application_security_group_ids": {
						Description: "A List of destination Application Security Group IDs.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         pluginsdk.HashStringInsensitively,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},

					"source_application_security_group_ids": {
						Description: "A List of source Application Security Group IDs.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         pluginsdk.HashStringInsensitively,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},

					"access": {
						Description:  "Whether network traffic is allowed or denied.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(networksecuritygroups.PossibleValuesForSecurityRuleAccess(), false),
					},

					"priority": {
						Description:  "The priority of the rule. The priority number must be unique for each rule in the collection. The lower the priority number, the higher the priority of the rule.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(100, 4096),
					},

					"direction": {
						Description:  "The direction specifies if rule will be evaluated on incoming or outgoing traffic.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(networksecuritygroups.PossibleValuesForSecurityRuleDirection(), false),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
