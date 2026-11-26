// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall_network_rule_collection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func firewallNetworkRuleCollectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Network Rule Collection which must be unique within the Firewall.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallName,
		},

		"azure_firewall_name": {
			Description:  "The name of the Firewall in which the Network Rule Collection should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"priority": {
			Description:  "The priority of the rule collection.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(100, 65000),
		},

		"action": {
			Description:  "The action the rule will apply to matching traffic.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(azurefirewalls.PossibleValuesForAzureFirewallRCActionType(), false),
		},

		"rule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the rule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"description": {
						Description: "A description for the rule.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"source_addresses": {
						Description: "A list of source IP addresses and/or IP ranges.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"source_ip_groups": {
						Description: "A list of IP Group IDs for the rule.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"destination_addresses": {
						Description: "Either a list of destination IP addresses and/or IP ranges, or a list of destination [Service Tags](https://learn.microsoft.com/azure/virtual-network/service-tags-overview#available-service-tags).",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"destination_ip_groups": {
						Description: "A list of destination IP Group IDs for the rule.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"destination_ports": {
						Description: "A list of destination ports.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"destination_fqdns": {
						Description: "A list of destination FQDNS for the rule.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"protocols": {
						Description: "A list of protocols.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(azurefirewalls.PossibleValuesForAzureFirewallNetworkRuleProtocol(), false),
						},
					},
				},
			},
		},
	}
}
