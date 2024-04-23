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

func (r ElasticSANDataSource) ResourceType() string {
	return "azurerm_elastic_san"
}

func (r ElasticSANDataSource) ModelObject() interface{} {
	return &ElasticSANResourceModel{}
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

			var model ElasticSANResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := elasticsans.NewElasticSanID(subscriptionId, model.ResourceGroupName, model.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s does not exist", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ElasticSANResourceModel{
				Name:              model.Name,
				ResourceGroupName: model.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)

				prop := model.Properties
				state.Sku = FlattenSku(prop.Sku)
				state.Zones = zones.Flatten(prop.AvailabilityZones)
				state.BaseSizeInTiB = int(prop.BaseSizeTiB)
				state.ExtendedSizeInTiB = int(prop.ExtendedCapacitySizeTiB)
				state.TotalIops = int(pointer.From(prop.TotalIops))
				state.TotalMBps = int(pointer.From(prop.TotalMBps))
				state.TotalSizeInTiB = int(pointer.From(prop.TotalSizeTiB))
				state.TotalVolumeSizeInGiB = int(pointer.From(prop.TotalVolumeSizeGiB))
				state.VolumeGroupCount = int(pointer.From(prop.VolumeGroupCount))
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
