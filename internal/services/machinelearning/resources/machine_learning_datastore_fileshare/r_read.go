// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_datastore_fileshare

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/datastore"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileshares"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	storageAccountHelper "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Datastore
			storageClient := metadata.Client.Storage
			subscriptionId := metadata.Client.Account.SubscriptionId

			id, err := datastore.ParseDataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			workspaceId := workspaces.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)
			model := MachineLearningDataStoreFileShareModel{
				Name:        *resp.Model.Name,
				WorkSpaceID: workspaceId.ID(),
			}

			data := resp.Model.Properties.(datastore.AzureFileDatastore)
			serviceDataIdentity := ""
			if v := data.ServiceDataAccessAuthIdentity; v != nil {
				serviceDataIdentity = string(*v)
			}
			model.ServiceDataIdentity = serviceDataIdentity

			var storageAccount *storageAccountHelper.AccountDetails
			if fileShareIdStr := metadata.ResourceData.Get("storage_fileshare_id").(string); fileShareIdStr != "" {
				shareId, err := parseStorageFileShareID(fileShareIdStr)
				if err != nil {
					return err
				}
				storageAccount, err = storageClient.GetAccount(ctx, commonids.NewStorageAccountID(shareId.SubscriptionId, shareId.ResourceGroupName, shareId.StorageAccountName))
				if err != nil {
					return fmt.Errorf("retrieving Account %q for Share %q: %s", data.AccountName, data.FileShareName, err)
				}
			} else {
				// In the case of import, we cannot rely on having a value for `storage_container_id` so we need to fallback on listing the accounts to search.
				storageAccount, err = storageClient.FindAccount(ctx, subscriptionId, data.AccountName)
				if err != nil {
					return fmt.Errorf("retrieving Account %q for Share %q: %s", data.AccountName, data.FileShareName, err)
				}
			}

			if storageAccount == nil {
				return fmt.Errorf("unable to locate Storage Account %q", data.AccountName)
			}

			model.StorageFileShareID = fileshares.NewShareID(storageAccount.StorageAccountId.SubscriptionId, storageAccount.StorageAccountId.ResourceGroupName, data.AccountName, data.FileShareName).ID()

			model.IsDefault = *data.IsDefault

			if v, ok := metadata.ResourceData.GetOk("account_key"); ok {
				if v.(string) != "" {
					model.AccountKey = v.(string)
				}
			}

			if v, ok := metadata.ResourceData.GetOk("shared_access_signature"); ok {
				if v.(string) != "" {
					model.SharedAccessSignature = v.(string)
				}
			}

			model.Description = pointer.From(data.Description)

			if data.Tags != nil {
				model.Tags = *data.Tags
			}

			return metadata.Encode(&model)
		},
	}
}

// parses the `storage_fileshare_id` into a fileshares.ShareId needed for 5.0 where the id changes
// from  "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/%s/fileshares/%s"
// to "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/fileServices/default/shares/%s"
// since we only use id.StorageAccountName and id.FileshareName parsing into a fileshares.ShareId works for both 4.0 and 5.0
func parseStorageFileShareID(input string) (*fileshares.ShareId, error) {
	return fileshares.ParseShareID(input)
}
