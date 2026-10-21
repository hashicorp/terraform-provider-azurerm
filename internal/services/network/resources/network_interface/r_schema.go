// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func networkInterfaceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"ip_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"subnet_id": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     commonids.ValidateSubnetID,
					},

					"private_ip_address": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// Note: O+C because Azure assigns a private IP from the subnet when not specified
						Computed: true,
					},

					"private_ip_address_version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(networkinterfaces.IPVersionIPvFour),
						ValidateFunc: validation.StringInSlice(networkinterfaces.PossibleValuesForIPVersion(), false),
					},

					"private_ip_address_allocation": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(networkinterfaces.PossibleValuesForIPAllocationMethod(), false),
					},

					"public_ip_address_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidatePublicIPAddressID,
					},

					"primary": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
					},

					"gateway_load_balancer_frontend_ip_configuration_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.AsGeneratedID(loadbalancers.ParseFrontendIPConfigurationIDInsensitively),
					},
				},
			},
		},

		// Optional
		"auxiliary_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(networkinterfaces.NetworkInterfaceAuxiliaryModeAcceleratedConnections),
				string(networkinterfaces.NetworkInterfaceAuxiliaryModeFloating),
				string(networkinterfaces.NetworkInterfaceAuxiliaryModeMaxConnections),
			}, false),
			RequiredWith: []string{"auxiliary_sku"},
		},

		"auxiliary_sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(networkinterfaces.NetworkInterfaceAuxiliarySkuAEight),
				string(networkinterfaces.NetworkInterfaceAuxiliarySkuAFour),
				string(networkinterfaces.NetworkInterfaceAuxiliarySkuAOne),
				string(networkinterfaces.NetworkInterfaceAuxiliarySkuATwo),
			}, false),
			RequiredWith: []string{"auxiliary_mode"},
		},

		"dns_servers": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"accelerated_networking_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"ip_forwarding_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"internal_dns_name_label": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"internal_domain_name_suffix": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),

		// Computed
		"applied_dns_servers": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"mac_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"virtual_machine_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
