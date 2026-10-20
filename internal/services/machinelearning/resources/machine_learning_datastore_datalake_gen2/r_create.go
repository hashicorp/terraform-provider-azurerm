// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_datastore_datalake_gen2

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
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

			var model MachineLearningDataStoreDataLakeGen2Model
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
					return tf.ImportAsExistsError("azurerm_machine_learning_datastore_datalake_gen2", id.ID())
				}
			}

			containerId, err := commonids.ParseStorageContainerID(model.StorageContainerID)
			if err != nil {
				return err
			}

			datastoreRaw := datastore.DatastoreResource{
				Name: pointer.To(model.Name),
				Type: pointer.To(string(datastore.DatastoreTypeAzureDataLakeGenTwo)),
			}

			props := &datastore.AzureDataLakeGen2Datastore{
				SubscriptionId:                &containerId.SubscriptionId,
				ResourceGroup:                 &containerId.ResourceGroupName,
				AccountName:                   containerId.StorageAccountName,
				Filesystem:                    containerId.ContainerName,
				Endpoint:                      pointer.To(metadata.Client.Storage.StorageDomainSuffix),
				Description:                   pointer.To(model.Description),
				ServiceDataAccessAuthIdentity: pointer.ToEnum[datastore.ServiceDataAccessAuthIdentity](model.ServiceDataIdentity),
				Tags:                          pointer.To(model.Tags),
			}

			var creds datastore.DatastoreCredentials = datastore.NoneDatastoreCredentials{}

			if len(model.TenantID) != 0 && len(model.ClientID) != 0 && len(model.ClientSecret) != 0 {
				resourceId, ok := metadata.Client.Account.Environment.DataLake.ResourceIdentifier()
				if !ok {
					return fmt.Errorf("could not determine resource identifier for DataLake in the %q cloud environment", metadata.Client.Account.Environment.Name)
				}
				creds = datastore.ServicePrincipalDatastoreCredentials{
					AuthorityURL: pointer.To(model.AuthorityURL),
					ResourceURL:  resourceId,
					TenantId:     model.TenantID,
					ClientId:     model.ClientID,
					Secrets: datastore.ServicePrincipalDatastoreSecrets{
						ClientSecret: pointer.To(model.ClientSecret),
					},
				}
			}
			props.Credentials = creds
			datastoreRaw.Properties = props

			if _, err = client.CreateOrUpdate(ctx, id, datastoreRaw, datastore.CreateOrUpdateOperationOptions{SkipValidation: pointer.To(true)}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
