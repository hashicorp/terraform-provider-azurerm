// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dev_center_gallery

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/galleries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	devcenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := devcenterClient.From(metadata.Client).V20250201.Galleries

			var config DevCenterGalleryResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := account.From(metadata.Client).SubscriptionId

			devCenterId, err := galleries.ParseDevCenterID(config.DevCenterId)
			if err != nil {
				return err
			}

			id := galleries.NewGalleryID(subscriptionId, devCenterId.ResourceGroupName, devCenterId.DevCenterName, config.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := galleries.Gallery{
				Properties: &galleries.GalleryProperties{
					GalleryResourceId: config.GalleryResourceId,
				},
			}
			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
