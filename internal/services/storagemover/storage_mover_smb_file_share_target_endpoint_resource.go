// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name storage_mover_smb_file_share_target_endpoint -service-package-name storagemover -properties "name" -compare-values "subscription_id:storage_mover_id,resource_group_name:storage_mover_id,storage_mover_name:storage_mover_id"

type StorageMoverSmbFileShareTargetEndpointModel struct {
	Name             string `tfschema:"name"`
	StorageMoverId   string `tfschema:"storage_mover_id"`
	StorageAccountId string `tfschema:"storage_account_id"`
	FileShareName    string `tfschema:"file_share_name"`
	Description      string `tfschema:"description"`
}

type StorageMoverSmbFileShareTargetEndpointResource struct{}

var (
	_ sdk.ResourceWithIdentity = StorageMoverSmbFileShareTargetEndpointResource{}
	_ sdk.ResourceWithUpdate   = StorageMoverSmbFileShareTargetEndpointResource{}
)

func (r StorageMoverSmbFileShareTargetEndpointResource) ResourceType() string {
	return "azurerm_storage_mover_smb_file_share_target_endpoint"
}

func (r StorageMoverSmbFileShareTargetEndpointResource) ModelObject() interface{} {
	return &StorageMoverSmbFileShareTargetEndpointModel{}
}

func (r StorageMoverSmbFileShareTargetEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return endpoints.ValidateEndpointID
}

func (r StorageMoverSmbFileShareTargetEndpointResource) Identity() resourceids.ResourceId {
	return &endpoints.EndpointId{}
}

func (r StorageMoverSmbFileShareTargetEndpointResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9a-zA-Z][-_0-9a-zA-Z]{0,63}$`),
				`The name must be between 1 and 64 characters in length, begin with a letter or number, and may contain letters, numbers, dashes and underscore.`,
			),
		},

		"storage_mover_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: storagemovers.ValidateStorageMoverID,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"file_share_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageShareName,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 1024),
		},
	}
}

func (r StorageMoverSmbFileShareTargetEndpointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StorageMoverSmbFileShareTargetEndpointResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model StorageMoverSmbFileShareTargetEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.StorageMover.EndpointsClient
			storageMoverId, err := storagemovers.ParseStorageMoverID(model.StorageMoverId)
			if err != nil {
				return err
			}

			id := endpoints.NewEndpointID(storageMoverId.SubscriptionId, storageMoverId.ResourceGroupName, storageMoverId.StorageMoverName, model.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := endpoints.Endpoint{
				Properties: endpoints.AzureStorageSmbFileShareEndpointProperties{
					FileShareName:            model.FileShareName,
					StorageAccountResourceId: model.StorageAccountId,
				},
			}

			if model.Description != "" {
				if v, ok := properties.Properties.(endpoints.AzureStorageSmbFileShareEndpointProperties); ok {
					v.Description = pointer.To(model.Description)
					properties.Properties = v
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
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

func (r StorageMoverSmbFileShareTargetEndpointResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.EndpointsClient

			id, err := endpoints.ParseEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StorageMoverSmbFileShareTargetEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if metadata.ResourceData.HasChange("description") {
				if v, ok := properties.Properties.(endpoints.AzureStorageSmbFileShareEndpointProperties); ok {
					v.Description = pointer.To(model.Description)
					properties.Properties = v
				}
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r StorageMoverSmbFileShareTargetEndpointResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.EndpointsClient

			id, err := endpoints.ParseEndpointID(metadata.ResourceData.Id())
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

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r StorageMoverSmbFileShareTargetEndpointResource) flatten(metadata sdk.ResourceMetaData, id *endpoints.EndpointId, model *endpoints.Endpoint) error {
	state := StorageMoverSmbFileShareTargetEndpointModel{
		Name:           id.EndpointName,
		StorageMoverId: storagemovers.NewStorageMoverID(id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName).ID(),
	}

	if model != nil {
		if v, ok := model.Properties.(endpoints.AzureStorageSmbFileShareEndpointProperties); ok {
			state.FileShareName = v.FileShareName
			state.StorageAccountId = v.StorageAccountResourceId
			state.Description = pointer.From(v.Description)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r StorageMoverSmbFileShareTargetEndpointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.EndpointsClient

			id, err := endpoints.ParseEndpointID(metadata.ResourceData.Id())
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
