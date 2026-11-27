// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_gateway_nat_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func vpnGatewayNatRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"vpn_gateway_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualwans.ValidateVpnGatewayID,
		},

		"external_mapping": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_space": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsCIDR,
					},

					"port_range": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"internal_mapping": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_space": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsCIDR,
					},

					"port_range": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"ip_configuration_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Instance0",
				"Instance1",
			}, false),
		},

		"mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(virtualwans.VpnNatRuleModeEgressSnat),
			ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForVpnNatRuleMode(), false),
		},

		"type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(virtualwans.VpnNatRuleTypeStatic),
			ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForVpnNatRuleType(), false),
		},
	}
}
