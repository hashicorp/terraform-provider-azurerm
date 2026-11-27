// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share_file

import (
	"fmt"

	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	storageClient2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/files"
)

func resourceStorageShareFileDelete(d *pluginsdk.ResourceData, meta any) error {
	storageClient := storageClient2.From(meta)
	subscriptionId := clientsAccount.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForDelete(clientsAccount.StopContext(meta), d)
	defer cancel()

	id, err := files.ParseFileID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for File %q (Share %q): %v", id.AccountId.AccountName, id.FileName, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
	}

	client, err := storageClient.FileShareFilesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building File Share File Client for %s: %v", account.StorageAccountId, err)
	}

	if _, err = client.Delete(ctx, id.ShareName, id.DirectoryPath, id.FileName); err != nil {
		return fmt.Errorf("deleting %s: %v", id, err)
	}

	return nil
}
