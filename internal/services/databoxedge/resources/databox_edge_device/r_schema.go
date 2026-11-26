// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databox_edge_device

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description:  "The name which should be used for this Databox Edge Device.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataboxEdgeName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Description:  "The `sku_name` is comprised of two segments separated by a hyphen (e.g. `TEA_1Node_UPS_Heater-Standard`). The first segment of the `sku_name` defines the `name` of the SKU. The second segment defines the `tier` of the `sku_name`. For more information see the [product documentation](https://learn.microsoft.com/azure/templates/microsoft.databoxedge/databoxedgedevices).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataboxEdgeDeviceSkuName,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"device_properties": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"configured_role_types": {
						Description: "Type of compute roles configured.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"culture": {
						Description: "The Data Box Edge/Gateway device culture.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"hcs_version": {
						Description: "The device software version number of the device (e.g. 1.2.18105.6).",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"capacity": {
						Description: "The Data Box Edge/Gateway device local capacity in MB.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"model": {
						Description: "The Data Box Edge/Gateway device model.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"status": {
						Description: "The status of the Data Box Edge/Gateway device.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"software_version": {
						Description: "The Data Box Edge/Gateway device software version.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"type": {
						Description: "The type of the Data Box Edge/Gateway device.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"node_count": {
						Description: "The number of nodes in the cluster.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},

					"serial_number": {
						Description: "The Serial Number of Data Box Edge/Gateway device.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"time_zone": {
						Description: "The Data Box Edge/Gateway device timezone.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}
