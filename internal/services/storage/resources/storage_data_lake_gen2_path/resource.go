// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_data_lake_gen2_path

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/paths"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageDataLakeGen2PathCreate,
		Read:   resourceStorageDataLakeGen2PathRead,
		Update: resourceStorageDataLakeGen2PathUpdate,
		Delete: resourceStorageDataLakeGen2PathDelete,

		Importer: helpers.ImporterValidatingStorageResourceIdThen(func(id, storageDomainSuffix string) error {
			_, err := paths.ParsePathID(id, storageDomainSuffix)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
			subscriptionId := meta.(*clients.Client).Account.SubscriptionId
			ctx, cancel := context.WithTimeout(ctx, d.Timeout(pluginsdk.TimeoutRead))
			defer cancel()

			storageClient := meta.(*clients.Client).Storage

			id, err := paths.ParsePathID(d.Id(), storageClient.StorageDomainSuffix)
			if err != nil {
				return []*pluginsdk.ResourceData{d}, err
			}

			// Retrieve the storage account
			account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving Account %q for Data Lake Gen2 Path %q in File System %q: %s", id.AccountId.AccountName, id.Path, id.FileSystemName, err)
			}
			if account == nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("unable to locate Storage Account: %q", id.AccountId.AccountName)
			}

			// Build the data plane client
			dataPlaneFilesystemsClient, err := storageClient.DataLakeFilesystemsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("building Data Lake Gen2 Filesystems Client: %v", err)
			}

			if _, err = dataPlaneFilesystemsClient.GetProperties(ctx, id.FileSystemName); err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving File System %q for Data Lake Gen2 Path %q in Account %q: %s", id.FileSystemName, id.Path, id.AccountId.AccountName, err)
			}

			d.Set("storage_account_id", account.StorageAccountId.ID())
			d.Set("filesystem_name", id.FileSystemName)

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: storageDataLakeGen2PathSchema(),
	}
}
