// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_data_lake_gen2_filesystem

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/filesystems"
)

func resourceStorageDataLakeGen2FileSystemDelete(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := filesystems.ParseFileSystemID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	// Retrieve the storage account properties
	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Data Lake Gen2 Filesystem %q: %v", id.AccountId.AccountName, id.FileSystemName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
	}

	// Build the data plane client
	dataPlaneFilesystemsClient, err := storageClient.DataLakeFilesystemsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Lake Gen2 Filesystems Client: %v", err)
	}

	resp, err := dataPlaneFilesystemsClient.Delete(ctx, id.FileSystemName)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %v", id, err)
		}
	}

	return nil
}
