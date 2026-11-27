// Copyright IBM Corp. 2014, 2026
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
			Description: "The name of the Network Interface.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"ip_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "A name used for this IP Configuration.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"subnet_id": {
						Description:      "The ID of the Subnet where this Network Interface should be located in.",
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ValidateFunc:     commonids.ValidateSubnetID,
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"private_ip_address": {
						Description: "The Static IP Address which should be used.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						// Note: O+C because Azure assigns a private IP from the subnet when not specified
						Computed: true,
					},

					"private_ip_address_version": {
						Description:  "The IP Version to use.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(networkinterfaces.IPVersionIPvFour),
						ValidateFunc: validation.StringInSlice(networkinterfaces.PossibleValuesForIPVersion(), false),
					},

					"private_ip_address_allocation": {
						Description:  "The allocation method used for the Private IP Address.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(networkinterfaces.PossibleValuesForIPAllocationMethod(), false),
					},

					"public_ip_address_id": {
						Description:  "Reference to a Public IP Address to associate with this NIC.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidatePublicIPAddressID,
					},

					"primary": {
						Description: "Whether this is the primary IP configuration.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},

					"gateway_load_balancer_frontend_ip_configuration_id": {
						Description:  "The Frontend IP Configuration ID of a Gateway SKU Load Balancer.",
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
			Description: "The auxiliary mode used to enable network high-performance feature on Network Virtual Appliances (NVAs). This feature offers competitive performance in Connections Per Second (CPS) optimization, along with improvements to handling large amounts of simultaneous connections.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(networkinterfaces.NetworkInterfaceAuxiliaryModeAcceleratedConnections),
				string(networkinterfaces.NetworkInterfaceAuxiliaryModeFloating),
				string(networkinterfaces.NetworkInterfaceAuxiliaryModeMaxConnections),
			}, false),
			RequiredWith: []string{"auxiliary_sku"},
		},

		"auxiliary_sku": {
			Description: "The SKU used for the network high-performance feature on Network Virtual Appliances (NVAs).",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(networkinterfaces.NetworkInterfaceAuxiliarySkuAEight),
				string(networkinterfaces.NetworkInterfaceAuxiliarySkuAFour),
				string(networkinterfaces.NetworkInterfaceAuxiliarySkuAOne),
				string(networkinterfaces.NetworkInterfaceAuxiliarySkuATwo),
			}, false),
			RequiredWith: []string{"auxiliary_mode"},
		},

		"dns_servers": {
			Description: "A list of IP Addresses defining the DNS Servers which should be used for this Network Interface.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"accelerated_networking_enabled": {
			Description: "Whether accelerated networking is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"ip_forwarding_enabled": {
			Description: "Whether IP Forwarding is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"internal_dns_name_label": {
			Description:  "The (relative) DNS Name used for internal communications between Virtual Machines in the same Virtual Network.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"internal_domain_name_suffix": {
			Description: "Even if `internal_dns_name_label` is not specified, a DNS entry is created for the primary NIC of the VM. This DNS name can be constructed by concatenating the VM name with the value of `internal_domain_name_suffix`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),

		// Computed
		"applied_dns_servers": {
			Description: "If the Virtual Machine using this Network Interface is part of an Availability Set, then this list will have the union of all DNS servers from all Network Interfaces that are part of the Availability Set.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"mac_address": {
			Description: "The Media Access Control (MAC) Address of the Network Interface.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"private_ip_address": {
			Description: "The first private IP address of the network interface.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"private_ip_addresses": {
			Description: "The private IP addresses of the network interface.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"virtual_machine_id": {
			Description: "The ID of the Virtual Machine which this Network Interface is connected to.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
