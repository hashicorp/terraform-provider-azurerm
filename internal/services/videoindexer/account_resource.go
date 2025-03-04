// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package videoindexer

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/videoindexer/2024-01-01/accounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AccountResource struct{}

var _ sdk.ResourceWithUpdate = AccountResource{}

type AccountModel struct {
	Name          string                                     `tfschema:"name"`
	ResourceGroup string                                     `tfschema:"resource_group_name"`
	Location      string                                     `tfschema:"location"`
	Storage       []StorageModel                             `tfschema:"storage"`
	Identity      []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Tags          map[string]string                          `tfschema:"tags"`
}

type StorageModel struct {
	StorageAccountId       string `tfschema:"storage_account_id"`
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
}

func (r AccountResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ServicePlanName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"storage": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"storage_account_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},
					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"tags": tags.Schema(),
	}
}

func (r AccountResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AccountResource) ModelObject() interface{} {
	return &AccountModel{}
}

func (r AccountResource) ResourceType() string {
	return "azurerm_video_indexer_account"
}

func (r AccountResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accounts.ValidateAccountID
}

func (r AccountResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var account AccountModel
			if err := metadata.Decode(&account); err != nil {
				return err
			}

			client := metadata.Client.VideoIndexer.AccountClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := accounts.NewAccountID(subscriptionId, account.ResourceGroup, account.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(account.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			payload := accounts.Account{
				Identity: expandedIdentity,
				Location: location.Normalize(account.Location),
				Tags:     pointer.To(account.Tags),
				Properties: &accounts.AccountPropertiesForPutRequest{
					StorageServices: expandStorageForCreate(account.Storage),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r AccountResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.VideoIndexer.AccountClient
			id, err := accounts.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			account, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(account.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := AccountModel{
				Name:          id.AccountName,
				ResourceGroup: id.ResourceGroupName,
			}

			if model := account.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = flattenedIdentity

				if props := model.Properties; props != nil {
					state.Storage, err = flattenStorage(props.StorageServices)
					if err != nil {
						return fmt.Errorf("flattening `storage`: %+v", err)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AccountResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := accounts.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.VideoIndexer.AccountClient

			var account AccountModel
			if err := metadata.Decode(&account); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := accounts.AccountPatch{}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(account.Tags)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(account.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("storage") {
				payload.Properties = &accounts.AccountPropertiesForPatchRequest{
					StorageServices: expandStorageForUpdate(account.Storage),
				}
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (r AccountResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.VideoIndexer.AccountClient

			id, err := accounts.ParseAccountID(metadata.ResourceData.Id())
			metadata.Logger.Infof("deleting %s", id)

			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}

			return nil
		},
	}
}

func expandStorageForCreate(input []StorageModel) *accounts.StorageServicesForPutRequest {
	if len(input) == 0 {
		return &accounts.StorageServicesForPutRequest{}
	}

	return &accounts.StorageServicesForPutRequest{
		ResourceId:           pointer.To(input[0].StorageAccountId),
		UserAssignedIdentity: pointer.To(input[0].UserAssignedIdentityId),
	}
}

func expandStorageForUpdate(input []StorageModel) *accounts.StorageServicesForPatchRequest {
	if len(input) == 0 {
		return &accounts.StorageServicesForPatchRequest{}
	}

	return &accounts.StorageServicesForPatchRequest{
		UserAssignedIdentity: pointer.To(input[0].UserAssignedIdentityId),
	}
}

func flattenStorage(input *accounts.StorageServicesForPutRequest) ([]StorageModel, error) {
	if input == nil {
		return []StorageModel{}, nil
	}

	storage := StorageModel{}
	if v := pointer.From(input.ResourceId); v != "" {
		id, err := commonids.ParseStorageAccountID(v)
		if err != nil {
			return []StorageModel{}, err
		}
		storage.StorageAccountId = id.ID()
	}

	if v := pointer.From(input.UserAssignedIdentity); v != "" {
		id, err := commonids.ParseUserAssignedIdentityID(v)
		if err != nil {
			return []StorageModel{}, err
		}
		storage.UserAssignedIdentityId = id.ID()
	}

	return []StorageModel{storage}, nil
}
