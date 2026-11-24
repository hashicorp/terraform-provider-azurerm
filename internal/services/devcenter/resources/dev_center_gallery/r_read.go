// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dev_center_gallery

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/galleries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Galleries
			schema := DevCenterGalleryResourceSchema{}

			id, err := galleries.ParseGalleryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			devCenterId := galleries.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName)

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.DevCenterId = devCenterId.ID()
				schema.Name = id.GalleryName

				if props := model.Properties; props != nil {
					schema.GalleryResourceId = props.GalleryResourceId
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

type DevCenterGalleryResourceSchema struct {
	DevCenterId       string `tfschema:"dev_center_id"`
	GalleryResourceId string `tfschema:"shared_gallery_id"`
	Name              string `tfschema:"name"`
}
