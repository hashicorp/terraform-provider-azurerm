// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall_nat_rule_collection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func firewallNatRuleCollectionSchema() map[string]*pluginsdk.Schema {
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
			ValidateFunc: validation.StringInSlice(azurefirewalls.PossibleValuesForAzureFirewallNatRCActionType(), false),
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
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"description": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"translated_address": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"translated_port": {
						Type:     pluginsdk.TypeString,
						Required: true,
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
					"destination_addresses": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"destination_ports": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"protocols": {
						Type:     pluginsdk.TypeList,
						Required: true,
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
