// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share_directory

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/directories"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
)

func resourceStorageShareDirectoryRead(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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
		log.Printf("[WARN] Unable to determine Resource Group for Storage Share Directory %q (Share %s, Account %s) - assuming removed & removing from state", id.DirectoryPath, id.ShareName, id.AccountId.AccountName)
		d.SetId("")
		return nil
	}

	// Determine the file endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeFile)
	if err != nil {
		return fmt.Errorf("determining File endpoint: %v", err)
	}

	// Parse the file endpoint as a data plane account ID
	accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	storageShareId := shares.NewShareID(*accountId, id.ShareName)

	client, err := storageClient.FileShareDirectoriesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building File Share Client: %v", err)
	}

	props, err := client.Get(ctx, id.ShareName, id.DirectoryPath)
	if err != nil {
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("name", id.DirectoryPath)
	d.Set("storage_share_url", storageShareId.ID())

	if err = d.Set("metadata", helpers.FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("setting `metadata`: %v", err)
	}

	return nil
}
