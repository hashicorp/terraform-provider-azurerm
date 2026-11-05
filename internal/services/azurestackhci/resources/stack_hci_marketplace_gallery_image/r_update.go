// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_marketplace_gallery_image

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/marketplacegalleryimages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.MarketplaceGalleryImages

			id, err := marketplacegalleryimages.ParseMarketplaceGalleryImageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StackHCIMarketplaceGalleryImageResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := &marketplacegalleryimages.MarketplaceGalleryImagesUpdateRequest{}
			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = tags.Expand(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}
