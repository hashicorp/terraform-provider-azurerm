// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironmentsstorages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.StorageClient

			id, err := managedenvironmentsstorages.ParseStorageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			var state ContainerAppEnvironmentStorageModel

			state.Name = id.StorageName
			state.ContainerAppEnvironmentId = managedenvironmentsstorages.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName).ID()

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					if azureFile := props.AzureFile; azureFile != nil {
						state.AccountName = pointer.From(azureFile.AccountName)
						if azureFile.AccessMode != nil {
							state.AccessMode = string(*azureFile.AccessMode)
						}
						state.ShareName = pointer.From(azureFile.ShareName)
					} else if nfsAzureFile := props.NfsAzureFile; nfsAzureFile != nil {
						state.NfsServer = pointer.From(nfsAzureFile.Server)
						if nfsAzureFile.AccessMode != nil {
							state.AccessMode = string(*nfsAzureFile.AccessMode)
						}
						state.ShareName = pointer.From(nfsAzureFile.ShareName)
					}
				}
			}
			if keyFromConfig, ok := metadata.ResourceData.GetOk("access_key"); ok {
				state.AccessKey = keyFromConfig.(string)
			}

			return metadata.Encode(&state)
		},
	}
}
