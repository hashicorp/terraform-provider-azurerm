// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name fake_storage_mover -service-package-name storagemover -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary"

type FakeStorageMoverModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Description       string            `tfschema:"description"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

type FakeStorageMoverResource struct{}

var (
	_ sdk.ResourceWithIdentity = FakeStorageMoverResource{}
	_ sdk.ResourceWithUpdate   = FakeStorageMoverResource{}
)

func (r FakeStorageMoverResource) Identity() resourceids.ResourceId {
	return &storagemovers.FakeStorageMoverId{}
}

func (r FakeStorageMoverResource) ResourceType() string {
	return "azurerm_fake_storage_mover"
}

func (r FakeStorageMoverResource) ModelObject() interface{} {
	return &FakeStorageMoverModel{}
}

func (r FakeStorageMoverResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return storagemovers.ValidateFakeStorageMoverID
}

func (r FakeStorageMoverResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r FakeStorageMoverResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r FakeStorageMoverResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model FakeStorageMoverModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.FakeStorageMover.FakeStorageMoversClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := storagemovers.NewFakeStorageMoverID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &storagemovers.FakeStorageMover{
				Location:   location.Normalize(model.Location),
				Properties: &storagemovers.FakeStorageMoverProperties{},
				Tags:       &model.Tags,
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r FakeStorageMoverResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.FakeStorageMover.FakeStorageMoversClient

			id, err := storagemovers.ParseFakeStorageMoverID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model FakeStorageMoverModel
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

func (r FakeStorageMoverResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.FakeStorageMover.FakeStorageMoversClient

			id, err := storagemovers.ParseFakeStorageMoverID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			return r.flatten(metadata, id, model)
		},
	}
}

func (r FakeStorageMoverResource) flatten(metadata sdk.ResourceMetaData, id *storagemovers.FakeStorageMoverId, model *storagemovers.FakeStorageMover) error {
	state := FakeStorageMoverModel{
		Name:              id.FakeStorageMoverName,
		ResourceGroupName: id.ResourceGroupName,
		Location:          location.Normalize(model.Location),
	}

	description := ""
	if properties := model.Properties; properties != nil {
		if properties.Description != nil {
			description = *properties.Description
		}
	}
	state.Description = description

	if model.Tags != nil {
		state.Tags = *model.Tags
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r FakeStorageMoverResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.FakeStorageMover.FakeStorageMoversClient

			id, err := storagemovers.ParseFakeStorageMoverID(metadata.ResourceData.Id())
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
