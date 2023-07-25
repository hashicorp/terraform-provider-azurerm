// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databoxedge

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/devices"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type EdgeDeviceDataSource struct{}

var _ sdk.DataSource = EdgeDeviceDataSource{}

func (d EdgeDeviceDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataboxEdgeName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (d EdgeDeviceDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"location": commonschema.LocationComputed(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

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

		"tags": commonschema.TagsDataSource(),
	}
}

func (d EdgeDeviceDataSource) ModelObject() interface{} {
	return &EdgeDeviceModel{}
}

func (d EdgeDeviceDataSource) ResourceType() string {
	return "azurerm_databox_edge_device"
}

func (d EdgeDeviceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.DataboxEdge.DeviceClient

			var metaModel EdgeDeviceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := devices.NewDataBoxEdgeDeviceID(subscriptionId, metaModel.ResourceGroupName, metaModel.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := EdgeDeviceModel{
				Name:              id.DataBoxEdgeDeviceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if props := model.Properties; props != nil {
					state.DeviceProperties = flattenDeviceProperties(props)
				}
				state.SkuName = flattenDeviceSku(model.Sku)
				state.Tags = *model.Tags
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
