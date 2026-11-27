// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.StorageMoversClient

			id, err := storagemovers.ParseStorageMoverID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StorageMoverModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Properties.Description = &model.Description
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
