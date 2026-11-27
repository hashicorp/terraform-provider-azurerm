// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package bastion_host

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/bastionhosts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func bastionHostSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Bastion Host.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BastionHostName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"copy_paste_enabled": {
			Description: "Whether Copy/Paste feature is enabled for the Bastion Host.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"file_copy_enabled": {
			Description: "Whether File Copy feature is enabled for the Bastion Host.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"ip_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the IP configuration.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.BastionIPConfigName,
					},
					"subnet_id": {
						Description:  "Reference to a subnet in which this Bastion Host has been created.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.BastionSubnetName,
					},
					"public_ip_address_id": {
						Description:  "Reference to a Public IP Address to associate with this Bastion Host.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidatePublicIPAddressID,
					},
				},
			},
		},

		"ip_connect_enabled": {
			Description: "Whether IP Connect feature is enabled for the Bastion Host.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"kerberos_enabled": {
			Description: "Whether Kerberos authentication feature is enabled for the Bastion Host.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"scale_units": {
			Description:  "The number of scale units with which to provision the Bastion Host.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      2,
			ValidateFunc: validation.IntBetween(2, 50),
		},

		"shareable_link_enabled": {
			Description: "Whether Shareable Link feature is enabled for the Bastion Host.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"sku": {
			Description:  "The SKU of the Bastion Host.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(bastionhosts.BastionHostSkuNameBasic),
			ValidateFunc: validation.StringInSlice(bastionhosts.PossibleValuesForBastionHostSkuName(), false),
		},

		"tunneling_enabled": {
			Description: "Whether Tunneling feature is enabled for the Bastion Host.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"session_recording_enabled": {
			Description: "Whether Session Recording feature is enabled for the Bastion Host.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"virtual_network_id": {
			Description:  "The ID of the Virtual Network for the Developer Bastion Host.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"dns_name": {
			Description: "The FQDN for the Bastion Host.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"private_only_enabled": {
			Description: "Whether Private-Only deployment is enabled for the Bastion Host.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),
	}
}
