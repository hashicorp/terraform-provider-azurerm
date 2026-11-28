// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share_directory

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/directories"
)

func resourceStorageShareDirectoryDelete(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := directories.ParseDirectoryID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Directory %q (Share %q): %v", id.AccountId.AccountName, id.DirectoryPath, id.ShareName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
	}

	client, err := storageClient.FileShareDirectoriesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building File Share Client: %v", err)
	}

	if _, err = client.Delete(ctx, id.ShareName, id.DirectoryPath); err != nil {
		return fmt.Errorf("deleting %s: %v", id, err)
	}

	return nil
}
