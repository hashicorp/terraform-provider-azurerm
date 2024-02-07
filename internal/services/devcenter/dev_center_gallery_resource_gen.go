package devcenter

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/galleries"
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
			client := metadata.Client.DevCenter.V20230401.Galleries

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

			var payload galleries.Gallery
			if err := r.mapDevCenterGalleryResourceSchemaToGallery(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
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
			client := metadata.Client.DevCenter.V20230401.Galleries
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
				if err := r.mapGalleryToDevCenterGalleryResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
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
			client := metadata.Client.DevCenter.V20230401.Galleries

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

func (r DevCenterGalleryResource) mapDevCenterGalleryResourceSchemaToGalleryProperties(input DevCenterGalleryResourceSchema, output *galleries.GalleryProperties) error {
	output.GalleryResourceId = input.GalleryResourceId
	return nil
}

func (r DevCenterGalleryResource) mapGalleryPropertiesToDevCenterGalleryResourceSchema(input galleries.GalleryProperties, output *DevCenterGalleryResourceSchema) error {
	output.GalleryResourceId = input.GalleryResourceId
	return nil
}

func (r DevCenterGalleryResource) mapDevCenterGalleryResourceSchemaToGallery(input DevCenterGalleryResourceSchema, output *galleries.Gallery) error {

	if output.Properties == nil {
		output.Properties = &galleries.GalleryProperties{}
	}
	if err := r.mapDevCenterGalleryResourceSchemaToGalleryProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "GalleryProperties", "Properties", err)
	}

	return nil
}

func (r DevCenterGalleryResource) mapGalleryToDevCenterGalleryResourceSchema(input galleries.Gallery, output *DevCenterGalleryResourceSchema) error {

	if input.Properties == nil {
		input.Properties = &galleries.GalleryProperties{}
	}
	if err := r.mapGalleryPropertiesToDevCenterGalleryResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "GalleryProperties", "Properties", err)
	}

	return nil
}
