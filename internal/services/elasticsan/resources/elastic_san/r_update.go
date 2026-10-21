// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic_san

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.ElasticSans

			id, err := elasticsans.ParseElasticSanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ElasticSANResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := elasticsans.ElasticSanUpdate{}

			if metadata.ResourceData.HasChange("base_size_in_tib") {
				if payload.Properties == nil {
					payload.Properties = &elasticsans.ElasticSanUpdateProperties{}
				}
				payload.Properties.BaseSizeTiB = pointer.To(config.BaseSizeInTiB)
			}

			if metadata.ResourceData.HasChange("extended_size_in_tib") {
				if payload.Properties == nil {
					payload.Properties = &elasticsans.ElasticSanUpdateProperties{}
				}
				payload.Properties.ExtendedCapacitySizeTiB = pointer.To(config.ExtendedSizeInTiB)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(config.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
