// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_static_website

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

func (a Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			accountDetails, err := storageClient.GetAccount(ctx, *id)
			if err != nil {
				//nolint:nilerr // If we don't find the account we can safely assume we don't need to remove the website since it must already be deleted
				return nil
			}

			properties := accounts.StorageServiceProperties{
				StaticWebsite: &accounts.StaticWebsite{},
			}

			client, err := storageClient.AccountsDataPlaneClient(ctx, *accountDetails, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Accounts Data Plane Client: %s", err)
			}

			if _, err = client.SetServiceProperties(ctx, id.StorageAccountName, properties); err != nil {
				return fmt.Errorf("deleting static website for %s: %+v", id, err)
			}

			return nil
		},
	}
}
