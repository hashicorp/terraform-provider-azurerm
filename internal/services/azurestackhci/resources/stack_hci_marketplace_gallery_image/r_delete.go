// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_marketplace_gallery_image

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/marketplacegalleryimages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 1 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.MarketplaceGalleryImages

			id, err := marketplacegalleryimages.ParseMarketplaceGalleryImageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
