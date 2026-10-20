// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallName,
		},

		"azure_firewall_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"priority": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(100, 65000),
		},

		"action": {
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.NoZeroValues,
					},
					"source_addresses": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"source_ip_groups": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"description": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"fqdn_tags": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"target_fqdns": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"protocol": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(azurefirewalls.PossibleValuesForAzureFirewallApplicationRuleProtocolType(), false),
								},
								"port": {
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
