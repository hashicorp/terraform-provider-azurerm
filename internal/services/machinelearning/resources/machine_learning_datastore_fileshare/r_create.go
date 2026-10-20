// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_datastore_fileshare

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/datastore"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Datastore
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MachineLearningDataStoreFileShareModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			workspaceId, err := workspaces.ParseWorkspaceID(model.WorkSpaceID)
			if err != nil {
				return err
			}

			id := datastore.NewDataStoreID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_machine_learning_datastore_fileshare", id.ID())
				}
			}

			shareId, err := parseStorageFileShareID(model.StorageFileShareID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: pointer.To(model.Name),
				Type: pointer.To(string(datastore.DatastoreTypeAzureFile)),
			}

			props := &datastore.AzureFileDatastore{
				AccountName:                   shareId.StorageAccountName,
				Endpoint:                      pointer.To(metadata.Client.Storage.StorageDomainSuffix),
				FileShareName:                 shareId.ShareName,
				Description:                   pointer.To(model.Description),
				ServiceDataAccessAuthIdentity: pointer.ToEnum[datastore.ServiceDataAccessAuthIdentity](model.ServiceDataIdentity),
				Tags:                          pointer.To(model.Tags),
			}

			accountKey := model.AccountKey
			if accountKey != "" {
				props.Credentials = datastore.AccountKeyDatastoreCredentials{
					Secrets: datastore.AccountKeyDatastoreSecrets{
						Key: pointer.To(accountKey),
					},
				}
			}

			sasToken := model.SharedAccessSignature
			if sasToken != "" {
				props.Credentials = datastore.SasDatastoreCredentials{
					Secrets: datastore.SasDatastoreSecrets{
						SasToken: pointer.To(sasToken),
					},
				}
			}
			datastoreRaw.Properties = props

			if _, err = client.CreateOrUpdate(ctx, id, datastoreRaw, datastore.CreateOrUpdateOperationOptions{SkipValidation: pointer.To(true)}); err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
