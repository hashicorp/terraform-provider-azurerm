// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_nat_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkgateway/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualNetworkGatewayNatRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Virtual Network Gateway NAT Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"virtual_network_gateway_id": {
			Description:  "The ID of the Virtual Network Gateway that this Virtual Network Gateway NAT Rule belongs to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualnetworkgateways.ValidateVirtualNetworkGatewayID,
		},

		"external_mapping": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_space": {
						Description:  "The string CIDR representing the address space for the Virtual Network Gateway NAT Rule external mapping.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsCIDR,
					},

					"port_range": {
						Description:  "The single port range for the Virtual Network Gateway NAT Rule external mapping.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"internal_mapping": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_space": {
						Description:  "The string CIDR representing the address space for the Virtual Network Gateway NAT Rule internal mapping.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsCIDR,
					},

					"port_range": {
						Description:  "The single port range for the Virtual Network Gateway NAT Rule internal mapping.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"mode": {
			Description:  "The source NAT direction of the Virtual Network Gateway NAT.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(networkgateways.VpnNatRuleModeEgressSnat),
			ValidateFunc: validation.StringInSlice(networkgateways.PossibleValuesForVpnNatRuleMode(), false),
		},

		"type": {
			Description:  "The type of the Virtual Network Gateway NAT Rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(networkgateways.VpnNatRuleTypeStatic),
			ValidateFunc: validation.StringInSlice(networkgateways.PossibleValuesForVpnNatRuleType(), false),
		},

		"ip_configuration_id": {
			Description:  "The ID of the IP Configuration this Virtual Network Gateway NAT Rule applies to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.VirtualNetworkGatewayIpConfigurationID,
		},
	}
}
