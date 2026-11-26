// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall_application_rule_collection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func firewallApplicationRuleCollectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Application Rule Collection which must be unique within the Firewall.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallName,
		},

		"azure_firewall_name": {
			Description:  "The name of the Firewall in which the Application Rule Collection should be created.",
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
						ValidateFunc: validation.NoZeroValues,
					},
					"source_addresses": {
						Description: "A list of source IP addresses and/or IP ranges.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"source_ip_groups": {
						Description: "A list of source IP Group IDs for the rule.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"description": {
						Description: "A description for the rule.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"fqdn_tags": {
						Description: "A list of FQDN tags.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"target_fqdns": {
						Description: "A list of FQDNs.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"protocol": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Description:  "The type of connection.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(azurefirewalls.PossibleValuesForAzureFirewallApplicationRuleProtocolType(), false),
								},
								"port": {
									Description:  "A port for the connection.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IsPortNumber,
								},
							},
						},
					},
				},
			},
		},
	}
}
