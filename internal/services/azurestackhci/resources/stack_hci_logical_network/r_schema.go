// Copyright IBM Corp. 2014, 2025
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][\-\.\_a-zA-Z0-9]{0,62}[a-zA-Z0-9]$`),
				"name must begin and end with an alphanumeric character, be between 2 and 64 characters in length and can only contain alphanumeric characters, hyphens, periods or underscores.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: customlocations.ValidateCustomLocationID,
		},

		"virtual_switch_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dns_servers": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(logicalnetworks.PossibleValuesForIPAllocationMethodEnum(), false),
					},

					"address_prefix": {
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},
								"end": {
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsCIDR,
								},

								"next_hop_ip_address": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.IsIPv4Address,
								},

								"name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^[a-zA-Z0-9][\-\.\_a-zA-Z0-9]{0,78}[a-zA-Z0-9]$`),
										"name must be between 2 and 80 characters and can only contain alphanumeric characters, hyphen, dot and underline",
									),
								},
							},
						},
					},

					"vlan_id": {
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
