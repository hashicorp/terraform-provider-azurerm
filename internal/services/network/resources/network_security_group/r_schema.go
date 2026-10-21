// Copyright IBM Corp. 2014, 2025
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"security_rule": {
			Type:       pluginsdk.TypeSet,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			Set:        hashNetworkSecurityRule,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"description": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(0, 140),
					},

					"protocol": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(networksecuritygroups.PossibleValuesForSecurityRuleProtocol(), false),
					},

					"source_port_range": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"source_port_ranges": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						Set:      pluginsdk.HashString,
					},

					"destination_port_range": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"destination_port_ranges": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						Set:      pluginsdk.HashString,
					},

					"source_address_prefix": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"source_address_prefixes": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						Set:      pluginsdk.HashString,
					},

					"destination_address_prefix": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"destination_address_prefixes": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						Set:      pluginsdk.HashString,
					},

					"destination_application_security_group_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
						},
						Set: pluginsdk.HashStringInsensitively,
					},

					"source_application_security_group_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
						},
						Set: pluginsdk.HashStringInsensitively,
					},

					"access": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(networksecuritygroups.PossibleValuesForSecurityRuleAccess(), false),
					},

					"priority": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(100, 4096),
					},

					"direction": {
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
