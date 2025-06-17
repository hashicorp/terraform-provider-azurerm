// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devboxdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/images"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = DevCenterDevBoxDefinitionResource{}
	_ sdk.ResourceWithUpdate = DevCenterDevBoxDefinitionResource{}
)

type DevCenterDevBoxDefinitionResource struct{}

func (r DevCenterDevBoxDefinitionResource) ModelObject() interface{} {
	return &DevCenterDevBoxDefinitionResourceModel{}
}

type DevCenterDevBoxDefinitionResourceModel struct {
	Name             string            `tfschema:"name"`
	Location         string            `tfschema:"location"`
	DevCenterId      string            `tfschema:"dev_center_id"`
	ImageReferenceId string            `tfschema:"image_reference_id"`
	SkuName          string            `tfschema:"sku_name"`
	Tags             map[string]string `tfschema:"tags"`
}

func (r DevCenterDevBoxDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return devboxdefinitions.ValidateDevCenterDevBoxDefinitionID
}

func (r DevCenterDevBoxDefinitionResource) ResourceType() string {
	return "azurerm_dev_center_dev_box_definition"
}

func (r DevCenterDevBoxDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevCenterDevBoxDefinitionName,
		},

		"location": commonschema.Location(),

		"dev_center_id": commonschema.ResourceIDReferenceRequiredForceNew(&devboxdefinitions.DevCenterId{}),

		"image_reference_id": commonschema.ResourceIDReferenceRequired(&images.GalleryImageId{}),

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r DevCenterDevBoxDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DevCenterDevBoxDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevBoxDefinitions
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DevCenterDevBoxDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterId, err := devboxdefinitions.ParseDevCenterID(model.DevCenterId)
			if err != nil {
				return err
			}

			id := devboxdefinitions.NewDevCenterDevBoxDefinitionID(subscriptionId, devCenterId.ResourceGroupName, devCenterId.DevCenterName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := devboxdefinitions.DevBoxDefinition{
				Location: location.Normalize(model.Location),
				Properties: &devboxdefinitions.DevBoxDefinitionProperties{
					ImageReference: &devboxdefinitions.ImageReference{
						Id: pointer.To(model.ImageReferenceId),
					},
					Sku: expandDevCenterDevBoxDefinitionSku(model.SkuName),
				},
				Tags: pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DevCenterDevBoxDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevBoxDefinitions

			id, err := devboxdefinitions.ParseDevCenterDevBoxDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := DevCenterDevBoxDefinitionResourceModel{
				Name:        id.DevBoxDefinitionName,
				DevCenterId: devboxdefinitions.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					if v := props.ImageReference; v != nil {
						state.ImageReferenceId = pointer.From(v.Id)
					}

					if v := props.Sku; v != nil {
						state.SkuName = flattenDevCenterDevBoxDefinition(props.Sku)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DevCenterDevBoxDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevBoxDefinitions

			id, err := devboxdefinitions.ParseDevCenterDevBoxDefinitionID(metadata.ResourceData.Id())
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

func (r DevCenterDevBoxDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevBoxDefinitions

			id, err := devboxdefinitions.ParseDevCenterDevBoxDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DevCenterDevBoxDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := devboxdefinitions.DevBoxDefinitionUpdate{
				Properties: &devboxdefinitions.DevBoxDefinitionUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("image_reference_id") {
				parameters.Properties.ImageReference = &devboxdefinitions.ImageReference{
					Id: pointer.To(model.ImageReferenceId),
				}
			}

			if metadata.ResourceData.HasChange("sku_name") {
				parameters.Properties.Sku = expandDevCenterDevBoxDefinitionSku(model.SkuName)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandDevCenterDevBoxDefinitionSku(input string) *devboxdefinitions.Sku {
	if input == "" {
		return nil
	}

	result := &devboxdefinitions.Sku{
		Name: input,
	}

	return result
}

func flattenDevCenterDevBoxDefinition(input *devboxdefinitions.Sku) string {
	var skuName string
	if input == nil {
		return skuName
	}

	skuName = pointer.From(input).Name

	return skuName
}
