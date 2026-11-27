// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_data_lake_gen2_path

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/paths"
	"github.com/jackofallops/giovanni/storage/accesscontrol"
)

func resourceStorageDataLakeGen2PathRead(d *pluginsdk.ResourceData, meta any) error {
	storageClient := client.From(meta)
	subscriptionId := clientsAccount.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(clientsAccount.StopContext(meta), d)
	defer cancel()

	id, err := paths.ParsePathID(d.Id(), storageClient.StorageDomainSuffix)
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
	dataPlanePathsClient, err := storageClient.DataLakePathsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Lake Gen2 Paths Client: %v", err)
	}

	resp, err := dataPlanePathsClient.GetProperties(ctx, id.FileSystemName, id.Path, paths.GetPropertiesInput{Action: paths.GetPropertiesActionGetStatus})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Path %q does not exist in File System %q in Storage Account %q - removing from state...", id.Path, id.FileSystemName, id.AccountId.AccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("path", id.Path)
	d.Set("resource", resp.ResourceType)
	d.Set("owner", resp.Owner)
	d.Set("group", resp.Group)

	// The above `getStatus` API request doesn't return the ACLs
	// Have to make a `getAccessControl` request, but that doesn't return all fields either!
	resp, err = dataPlanePathsClient.GetProperties(ctx, id.FileSystemName, id.Path, paths.GetPropertiesInput{Action: paths.GetPropertiesActionGetAccessControl})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Path %q does not exist in File System %q in Storage Account %q - removing from state...", id.Path, id.FileSystemName, id.AccountId.AccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving ACLs for %s: %v", id, err)
	}

	acl, err := accesscontrol.ParseACL(resp.ACL)
	if err != nil {
		return fmt.Errorf("parsing response ACL %q: %v", resp.ACL, err)
	}
	d.Set("ace", helpers.FlattenDataLakeGen2AceList(d, acl))

	return nil
}
