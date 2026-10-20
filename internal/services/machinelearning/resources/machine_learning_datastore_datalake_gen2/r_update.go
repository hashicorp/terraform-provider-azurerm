// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_datastore_datalake_gen2

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

			var state MachineLearningDataStoreDataLakeGen2Model
			if err := metadata.Decode(&state); err != nil {
				return err
			}
			containerId, err := commonids.ParseStorageContainerID(state.StorageContainerID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: pointer.To(id.DataStoreName),
				Type: pointer.To(string(datastore.DatastoreTypeAzureDataLakeGenTwo)),
			}

			props := &datastore.AzureDataLakeGen2Datastore{
				SubscriptionId:                &containerId.SubscriptionId,
				ResourceGroup:                 &containerId.ResourceGroupName,
				AccountName:                   containerId.StorageAccountName,
				Filesystem:                    containerId.ContainerName,
				Description:                   pointer.To(state.Description),
				ServiceDataAccessAuthIdentity: pointer.ToEnum[datastore.ServiceDataAccessAuthIdentity](state.ServiceDataIdentity),
				Tags:                          pointer.To(state.Tags),
			}

			var creds datastore.DatastoreCredentials = datastore.NoneDatastoreCredentials{}

			if len(state.TenantID) != 0 && len(state.ClientID) != 0 && len(state.ClientSecret) != 0 {
				resourceId, ok := metadata.Client.Account.Environment.DataLake.ResourceIdentifier()
				if !ok {
					return fmt.Errorf("could not determine resource identifier for DataLake in the %q cloud environment", metadata.Client.Account.Environment.Name)
				}
				creds = datastore.ServicePrincipalDatastoreCredentials{
					AuthorityURL: pointer.To(state.AuthorityURL),
					ResourceURL:  resourceId,
					TenantId:     state.TenantID,
					ClientId:     state.ClientID,
					Secrets: datastore.ServicePrincipalDatastoreSecrets{
						ClientSecret: pointer.To(state.ClientSecret),
					},
				}
			}
			props.Credentials = creds
			datastoreRaw.Properties = props

			if _, err = client.CreateOrUpdate(ctx, *id, datastoreRaw, datastore.CreateOrUpdateOperationOptions{SkipValidation: pointer.To(true)}); err != nil {
				return fmt.Errorf("creating/updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
