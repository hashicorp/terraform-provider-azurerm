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

			if existing.Model.Properties == nil || existing.Model.Properties.AzureFile == nil {
				return fmt.Errorf("could not update %s: existing resource is missing `AzureFile` properties", *id)
			}

			// This *must* be sent, and is currently the only updatable property on the resource.
			existing.Model.Properties.AzureFile.AccountKey = pointer.To(metadata.ResourceData.Get("access_key").(string))

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
