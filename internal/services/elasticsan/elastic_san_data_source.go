// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elasticsan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ElasticSANDataSource struct{}

var _ sdk.DataSource = ElasticSANDataSource{}

type ElasticSANDataSourceModel struct {
	BaseSizeInTiB        int64                        `tfschema:"base_size_in_tib"`
	ExtendedSizeInTiB    int64                        `tfschema:"extended_size_in_tib"`
	Location             string                       `tfschema:"location"`
	Name                 string                       `tfschema:"name"`
	ResourceGroupName    string                       `tfschema:"resource_group_name"`
	Sku                  []ElasticSANResourceSkuModel `tfschema:"sku"`
	Tags                 map[string]interface{}       `tfschema:"tags"`
	TotalIops            int64                        `tfschema:"total_iops"`
	TotalMBps            int64                        `tfschema:"total_mbps"`
	TotalSizeInTiB       int64                        `tfschema:"total_size_in_tib"`
	TotalVolumeSizeInGiB int64                        `tfschema:"total_volume_size_in_gib"`
	VolumeGroupCount     int64                        `tfschema:"volume_group_count"`
	Zones                []string                     `tfschema:"zones"`
}

func (r ElasticSANDataSource) ResourceType() string {
	return "azurerm_elastic_san"
}

func (r ElasticSANDataSource) ModelObject() interface{} {
	return &ElasticSANDataSourceModel{}
}

func (r ElasticSANDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ElasticSanName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r ElasticSANDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"zones": commonschema.ZonesMultipleComputed(),

		"base_size_in_tib": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"extended_size_in_tib": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"sku": {
			Computed: true,
			Type:     pluginsdk.TypeList,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},

					"tier": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
				},
			},
		},

		"total_iops": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"total_mbps": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"total_size_in_tib": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"total_volume_size_in_gib": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"volume_group_count": {
			Computed: true,
			Type:     pluginsdk.TypeInt,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r ElasticSANDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.ElasticSans
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ElasticSANDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := elasticsans.NewElasticSanID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s does not exist", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)

				prop := model.Properties
				state.Sku = FlattenSku(prop.Sku)
				state.Zones = zones.Flatten(prop.AvailabilityZones)
				state.BaseSizeInTiB = prop.BaseSizeTiB
				state.ExtendedSizeInTiB = prop.ExtendedCapacitySizeTiB
				state.TotalIops = pointer.From(prop.TotalIops)
				state.TotalMBps = pointer.From(prop.TotalMBps)
				state.TotalSizeInTiB = pointer.From(prop.TotalSizeTiB)
				state.TotalVolumeSizeInGiB = pointer.From(prop.TotalVolumeSizeGiB)
				state.VolumeGroupCount = pointer.From(prop.VolumeGroupCount)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
