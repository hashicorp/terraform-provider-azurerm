// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_marketplace_gallery_image

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/marketplacegalleryimages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.MarketplaceGalleryImages

			var config StackHCIMarketplaceGalleryImageResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := marketplacegalleryimages.NewMarketplaceGalleryImageID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := marketplacegalleryimages.MarketplaceGalleryImages{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				ExtendedLocation: &marketplacegalleryimages.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(marketplacegalleryimages.ExtendedLocationTypesCustomLocation),
				},
				Properties: &marketplacegalleryimages.MarketplaceGalleryImageProperties{
					Identifier:       expandStackHCIMarketplaceGalleryImageIdentifier(config.Identifier),
					OsType:           marketplacegalleryimages.OperatingSystemTypes(config.OsType),
					HyperVGeneration: pointer.ToEnum[marketplacegalleryimages.HyperVGeneration](config.HypervGeneration),
					Version: &marketplacegalleryimages.GalleryImageVersion{
						Name: pointer.To(config.Version),
					},
				},
			}

			if config.StoragePathId != "" {
				payload.Properties.ContainerId = pointer.To(config.StoragePathId)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("performing create %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func expandStackHCIMarketplaceGalleryImageIdentifier(input []StackHCIMarketplaceGalleryImageIdentifier) *marketplacegalleryimages.GalleryImageIdentifier {
	if len(input) == 0 {
		return nil
	}

	v := input[0]

	return &marketplacegalleryimages.GalleryImageIdentifier{
		Offer:     v.Offer,
		Publisher: v.Publisher,
		Sku:       v.Sku,
	}
}
