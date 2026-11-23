// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_private_cloud

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func vmwarePrivateCloudSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"av20",
				"av36",
				"av36t",
				"av36p",
				"av36pt",
				"av48",
				"av48t",
				"av52",
				"av52t",
				"av64",
			}, false),
		},

		"management_cluster": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"size": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(3, 16),
					},

					"hosts": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"id": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"network_subnet_cidr": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsCIDR,
		},

		"internet_connection_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"nsxt_password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"vcenter_password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"circuit": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"express_route_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"express_route_private_peering_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"primary_subnet_cidr": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"secondary_subnet_cidr": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"hcx_cloud_manager_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"management_subnet_cidr": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"nsxt_certificate_thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"nsxt_manager_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"provisioning_subnet_cidr": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"vcenter_certificate_thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"vcsa_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"vmotion_subnet_cidr": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
