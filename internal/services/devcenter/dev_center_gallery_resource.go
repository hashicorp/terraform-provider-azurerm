package devcenter

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/galleries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = DevCenterGalleryResource{}

type DevCenterGalleryResource struct{}

func (r DevCenterGalleryResource) ModelObject() interface{} {
	return &DevCenterGalleryResourceSchema{}
}

type DevCenterGalleryResourceSchema struct {
	DevCenterId       string `tfschema:"dev_center_id"`
	GalleryResourceId string `tfschema:"shared_gallery_id"`
	Name              string `tfschema:"name"`
}

func (r DevCenterGalleryResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return galleries.ValidateGalleryID
}

func (r DevCenterGalleryResource) ResourceType() string {
	return "azurerm_dev_center_gallery"
}

func (r DevCenterGalleryResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dev_center_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"shared_gallery_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r DevCenterGalleryResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DevCenterGalleryResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Galleries

			var config DevCenterGalleryResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			devCenterId, err := galleries.ParseDevCenterID(config.DevCenterId)
			if err != nil {
				return err
			}

			id := galleries.NewGalleryID(subscriptionId, devCenterId.ResourceGroupName, devCenterId.DevCenterName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := galleries.Gallery{
				Properties: &galleries.GalleryProperties{
					GalleryResourceId: config.GalleryResourceId,
				},
			}
			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DevCenterGalleryResource) Read() sdk.ResourceFunc {
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

func (r DevCenterGalleryResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Galleries

			id, err := galleries.ParseGalleryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
