// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stack_hci_network_interface

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
			Description: "The name which should be used for this Azure Stack HCI Network Interface.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][\-\.\_a-zA-Z0-9]{0,78}[a-zA-Z0-9]$`),
				"name must begin and end with an alphanumeric character, be between 2 and 80 characters in length and can only contain alphanumeric characters, hyphens, periods or underscores.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": commonschema.ResourceIDReferenceRequiredForceNew(&customlocations.CustomLocationId{}),

		"ip_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": commonschema.ResourceIDReferenceRequiredForceNew(&logicalnetworks.LogicalNetworkId{}),

					"private_ip_address": {
						Description:  "The IPv4 address of the IP configuration.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsIPv4Address,
					},

					"gateway": {
						Description: "The IPv4 address of the gateway for the Network Interface.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"prefix_length": {
						Description: "The prefix length for the address of the Network Interface.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"dns_servers": {
			Description: "A list of IPv4 addresses of DNS servers available to VMs deployed in the Network Interface.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			ForceNew:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsIPv4Address,
			},
		},

		"mac_address": {
			Description:  "The MAC address of the Network Interface.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
