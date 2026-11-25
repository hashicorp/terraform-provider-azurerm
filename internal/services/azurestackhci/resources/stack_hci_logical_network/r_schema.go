// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stack_hci_logical_network

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Azure Stack HCI Logical Network.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][\-\.\_a-zA-Z0-9]{0,62}[a-zA-Z0-9]$`),
				"name must begin and end with an alphanumeric character, be between 2 and 64 characters in length and can only contain alphanumeric characters, hyphens, periods or underscores.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": {
			Description:  "The ID of Custom Location where the Azure Stack HCI Logical Network should exist.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: customlocations.ValidateCustomLocationID,
		},

		"virtual_switch_name": {
			Description:  "The name of the virtual switch on the cluster used to associate with the Azure Stack HCI Logical Network. Possible switch names can be retrieved by following this [Azure guide](https://learn.microsoft.com/azure-stack/hci/manage/create-logical-networks?tabs=azurecli#prerequisites).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dns_servers": {
			Description: "A list of IPv4 addresses of DNS servers available to VMs deployed in the Logical Networks.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsIPv4Address,
			},
		},

		"subnet": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_allocation_method": {
						Description:  "The IP address allocation method for the subnet.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(logicalnetworks.PossibleValuesForIPAllocationMethodEnum(), false),
					},

					"address_prefix": {
						Description:  "The address prefix in CIDR notation.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsCIDR,
					},

					"ip_pool": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"start": {
									Description:  "The IPv4 address of the start of the IP address pool.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},
								"end": {
									Description:  "The IPv4 address of the end of the IP address pool.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},
							},
						},
					},

					"route": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"address_prefix": {
									Description:  "The Address in CIDR notation.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsCIDR,
								},

								"next_hop_ip_address": {
									Description:  "The IPv4 address of the next hop.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},

								"name": {
									Description: "The name of the route.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
									ForceNew:    true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^[a-zA-Z0-9][\-\.\_a-zA-Z0-9]{0,78}[a-zA-Z0-9]$`),
										"name must be between 2 and 80 characters and can only contain alphanumeric characters, hyphen, dot and underline",
									),
								},
							},
						},
					},

					"vlan_id": {
						Description:  "The VLAN ID for the Logical Network.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
