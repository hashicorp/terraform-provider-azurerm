// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ElasticSan.Volumes

			id, err := volumes.ParseVolumeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ElasticSANVolumeResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := volumes.VolumeUpdate{}
			if metadata.ResourceData.HasChange("size_in_gib") {
				if payload.Properties == nil {
					payload.Properties = &volumes.VolumeUpdateProperties{}
				}
				payload.Properties.SizeGiB = pointer.To(config.SizeInGiB)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
