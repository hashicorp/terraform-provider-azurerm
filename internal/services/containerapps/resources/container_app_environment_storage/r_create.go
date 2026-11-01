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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.StorageClient

			var storage ContainerAppEnvironmentStorageModel

			if err := metadata.Decode(&storage); err != nil {
				return err
			}

			containerAppEnvironmentId, err := managedenvironmentsstorages.ParseManagedEnvironmentID(storage.ContainerAppEnvironmentId)
			if err != nil {
				return err
			}

			id := managedenvironmentsstorages.NewStorageID(metadata.Client.Account.SubscriptionId, containerAppEnvironmentId.ResourceGroupName, containerAppEnvironmentId.ManagedEnvironmentName, storage.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			accessMode := managedenvironmentsstorages.AccessMode(storage.AccessMode)
			managedEnvironmentStorage := managedenvironmentsstorages.ManagedEnvironmentStorage{}

			if storage.NfsServer != "" {
				managedEnvironmentStorage.Properties = &managedenvironmentsstorages.ManagedEnvironmentStorageProperties{
					NfsAzureFile: &managedenvironmentsstorages.NfsAzureFileProperties{
						AccessMode: &accessMode,
						Server:     pointer.To(storage.NfsServer),
						ShareName:  pointer.To(storage.ShareName),
					},
				}
			} else {
				managedEnvironmentStorage.Properties = &managedenvironmentsstorages.ManagedEnvironmentStorageProperties{
					AzureFile: &managedenvironmentsstorages.AzureFileProperties{
						AccessMode:  &accessMode,
						AccountKey:  pointer.To(storage.AccessKey),
						AccountName: pointer.To(storage.AccountName),
						ShareName:   pointer.To(storage.ShareName),
					},
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id, managedEnvironmentStorage); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
