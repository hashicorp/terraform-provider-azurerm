// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic_san

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.ElasticSans

			id, err := elasticsans.ParseElasticSanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ElasticSANResourceModel{
				Name:              id.ElasticSanName,
				ResourceGroupName: id.ResourceGroupName,
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

			return metadata.Encode(&state)
		},
	}
}

func ExpandSku(input []ElasticSANResourceSkuModel) elasticsans.Sku {
	if len(input) == 0 {
		return elasticsans.Sku{}
	}

	output := elasticsans.Sku{
		Name: elasticsans.SkuName(input[0].Name),
	}

	if input[0].Tier != "" {
		output.Tier = pointer.ToEnum[elasticsans.SkuTier](input[0].Tier)
	}

	return output
}

func FlattenSku(input elasticsans.Sku) []ElasticSANResourceSkuModel {
	return []ElasticSANResourceSkuModel{
		{
			Name: string(input.Name),
			Tier: pointer.FromEnum(input.Tier),
		},
	}
}
