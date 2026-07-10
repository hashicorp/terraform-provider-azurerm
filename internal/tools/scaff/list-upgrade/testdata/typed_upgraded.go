package testdata

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

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name storage_mover -service-package-name storagemover -properties "name,resource_group_name"

type StorageMoverModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Description       string            `tfschema:"description"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

type StorageMoverResource struct{}

var (
	_ sdk.ResourceWithIdentity = StorageMoverResource{}
	_ sdk.ResourceWithUpdate   = StorageMoverResource{}
)

func (r StorageMoverResource) Identity() resourceids.ResourceId {
	return &storagemovers.StorageMoverId{}
}

func (r StorageMoverResource) ResourceType() string {
	return "azurerm_storage_mover"
}

func (r StorageMoverResource) ModelObject() interface{} {
	return &StorageMoverModel{}
}

func (r StorageMoverResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return storagemovers.ValidateStorageMoverID
}

func (r StorageMoverResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r StorageMoverResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StorageMoverResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model StorageMoverModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.StorageMover.StorageMoversClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := storagemovers.NewStorageMoverID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &storagemovers.StorageMover{
				Location:   location.Normalize(model.Location),
				Properties: &storagemovers.StorageMoverProperties{},
				Tags:       &model.Tags,
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

func (r StorageMoverResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.StorageMoversClient

			id, err := storagemovers.ParseStorageMoverID(metadata.ResourceData.Id())
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

func (r StorageMoverResource) flatten(metadata sdk.ResourceMetaData, id *storagemovers.StorageMoverId, model *storagemovers.StorageMover) error {
	state := StorageMoverModel{
		Name:              id.StorageMoverName,
		ResourceGroupName: id.ResourceGroupName,
		Location:          location.Normalize(model.Location),
	}

	if model.Tags != nil {
		state.Tags = *model.Tags
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r StorageMoverResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.StorageMoversClient

			id, err := storagemovers.ParseStorageMoverID(metadata.ResourceData.Id())
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
