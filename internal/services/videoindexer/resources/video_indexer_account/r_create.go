// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package video_indexer_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/videoindexer/2025-04-01/accounts"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	videoindexerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/videoindexer/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var account AccountModel
			if err := metadata.Decode(&account); err != nil {
				return err
			}

			client := videoindexerClient.From(metadata.Client).AccountClient
			subscriptionId := clientsAccount.From(metadata.Client).SubscriptionId

			id := accounts.NewAccountID(subscriptionId, account.ResourceGroup, account.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("retrieving %s: %v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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
					PublicNetworkAccess: pointer.ToEnum[accounts.PublicNetworkAccess](account.PublicNetworkAccess),
					StorageServices:     expandStorageForCreate(account.Storage),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
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
