// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_datastore_blobstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/datastore"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Datastore

			id, err := datastore.ParseDataStoreID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MachineLearningDataStoreBlobStorageModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			containerId, err := commonids.ParseStorageContainerID(state.StorageContainerID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: pointer.To(id.DataStoreName),
				Type: pointer.To(string(datastore.DatastoreTypeAzureBlob)),
			}

			props := &datastore.AzureBlobDatastore{
				AccountName:                   pointer.To(containerId.StorageAccountName),
				ContainerName:                 pointer.To(containerId.ContainerName),
				Description:                   pointer.To(state.Description),
				ServiceDataAccessAuthIdentity: pointer.ToEnum[datastore.ServiceDataAccessAuthIdentity](state.ServiceDataAuthIdentity),
				IsDefault:                     pointer.To(state.IsDefault),
				Tags:                          pointer.To(state.Tags),
			}

			accountKey := state.AccountKey
			if accountKey != "" {
				props.Credentials = datastore.AccountKeyDatastoreCredentials{
					Secrets: datastore.AccountKeyDatastoreSecrets{
						Key: pointer.To(accountKey),
					},
				}
			}

			sasToken := state.SharedAccessSignature
			if sasToken != "" {
				props.Credentials = datastore.SasDatastoreCredentials{
					Secrets: datastore.SasDatastoreSecrets{
						SasToken: pointer.To(sasToken),
					},
				}
			}

			// If `service_data_auth_identity` is set to `WorkspaceSystemAssignedIdentity` or `WorkspaceUserAssignedIdentity`,
			// explicit credentials such as `account_key` or `shared_access_signature` must not be provided.
			// Only when `service_data_auth_identity` is set to `None`, one of `account_key` or `shared_access_signature` must be specified.
			// In addition, when `service_data_auth_identity` is set to either `WorkspaceSystemAssignedIdentity` or `WorkspaceUserAssignedIdentity`,
			// the `Credentials` field must include "CredentialsType": "None". Omitting this will result in a validation error.
			if accountKey == "" && sasToken == "" {
				props.Credentials = datastore.BaseDatastoreCredentialsImpl{
					CredentialsType: datastore.CredentialsTypeNone,
				}
			}
			datastoreRaw.Properties = props

			if _, err = client.CreateOrUpdate(ctx, *id, datastoreRaw, datastore.CreateOrUpdateOperationOptions{SkipValidation: pointer.To(true)}); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
