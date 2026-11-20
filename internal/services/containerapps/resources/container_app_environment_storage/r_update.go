// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironmentsstorages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.StorageClient

			id, err := managedenvironmentsstorages.ParseStorageID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var storage ContainerAppEnvironmentStorageModel
			if err := metadata.Decode(&storage); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s for update: %+v", *id, err)
			}

			if existing.Model == nil || existing.Model.Properties == nil || existing.Model.Properties.AzureFile == nil {
				return fmt.Errorf("could not update %s: existing resource is missing `AzureFile` properties", *id)
			}

			// One of the account key or Key Vault reference must be sent on every update.
			azureFile := existing.Model.Properties.AzureFile
			if len(storage.AccessKeyVaultSecret) > 0 {
				azureFile.AccountKey = nil
				azureFile.AccountKeyVaultProperties = expandContainerAppEnvironmentStorageKeyVaultSecret(storage.AccessKeyVaultSecret)
			} else {
				azureFile.AccountKey = pointer.To(storage.AccessKey)
				azureFile.AccountKeyVaultProperties = nil
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
