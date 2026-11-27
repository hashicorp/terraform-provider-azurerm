// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_data_lake_gen2_filesystem

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/filesystems"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageDataLakeGen2FileSystemCreate,
		Read:   resourceStorageDataLakeGen2FileSystemRead,
		Update: resourceStorageDataLakeGen2FileSystemUpdate,
		Delete: resourceStorageDataLakeGen2FileSystemDelete,

		Importer: helpers.ImporterValidatingStorageResourceIdThen(func(id, storageDomainSuffix string) error {
			_, err := filesystems.ParseFileSystemID(id, storageDomainSuffix)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
			storageClient := meta.(*clients.Client).Storage
			subscriptionId := meta.(*clients.Client).Account.SubscriptionId
			ctx, cancel := context.WithTimeout(ctx, d.Timeout(pluginsdk.TimeoutRead))
			defer cancel()

			id, err := filesystems.ParseFileSystemID(d.Id(), storageClient.StorageDomainSuffix)
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("parsing ID %q for import of Data Lake Gen2 File System: %v", d.Id(), err)
			}

			// we then need to look up the Storage Account ID
			account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving Account %q for Data Lake Gen2 File System %q: %s", id.AccountId.AccountName, id.FileSystemName, err)
			}
			if account == nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("unable to locate Storage Account: %q", id.AccountId.AccountName)
			}

			d.Set("storage_account_id", account.StorageAccountId.ID())

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: storageDataLakeGen2FilesystemSchema(),
	}
}
