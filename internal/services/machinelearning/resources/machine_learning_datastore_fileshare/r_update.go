// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_datastore_fileshare

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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

			var state MachineLearningDataStoreFileShareModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			shareId, err := parseStorageFileShareID(state.StorageFileShareID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: pointer.To(id.DataStoreName),
				Type: pointer.To(string(datastore.DatastoreTypeAzureFile)),
			}

			props := &datastore.AzureFileDatastore{
				AccountName:                   shareId.StorageAccountName,
				FileShareName:                 shareId.ShareName,
				Description:                   pointer.To(state.Description),
				ServiceDataAccessAuthIdentity: pointer.ToEnum[datastore.ServiceDataAccessAuthIdentity](state.ServiceDataIdentity),
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
			datastoreRaw.Properties = props

			if _, err = client.CreateOrUpdate(ctx, *id, datastoreRaw, datastore.CreateOrUpdateOperationOptions{SkipValidation: pointer.To(true)}); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
