// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataboxEdgeName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
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
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"culture": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"hcs_version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"capacity": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"model": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"status": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"software_version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"node_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"serial_number": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"time_zone": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}
