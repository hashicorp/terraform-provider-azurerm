// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_data_lake_gen2_filesystem

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/filesystems"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/paths"
	"github.com/jackofallops/giovanni/storage/accesscontrol"
)

func resourceStorageDataLakeGen2FileSystemRead(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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

	// Build the data plane clients
	dataPlaneFilesystemsClient, err := storageClient.DataLakeFilesystemsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Lake Gen2 Filesystems Client: %v", err)
	}
	dataPlanePathsClient, err := storageClient.DataLakePathsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Lake Gen2 Paths Client: %v", err)
	}

	resp, err := dataPlaneFilesystemsClient.GetProperties(ctx, id.FileSystemName)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] File System %q does not exist in Storage Account %q - removing from state...", id.FileSystemName, id.AccountId.AccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("name", id.FileSystemName)
	d.Set("default_encryption_scope", resp.DefaultEncryptionScope)

	if err = d.Set("properties", resp.Properties); err != nil {
		return fmt.Errorf("setting `properties`: %v", err)
	}

	var ace []any
	var owner, group string
	// acl is only enabled when `IsHnsEnabled` is true otherwise the rest api will report error
	if account.IsHnsEnabled {
		// The above `getStatus` API request doesn't return the ACLs
		// Have to make a `getAccessControl` request, but that doesn't return all fields either!
		payload := paths.GetPropertiesInput{
			Action: paths.GetPropertiesActionGetAccessControl,
		}
		pathResponse, err := dataPlanePathsClient.GetProperties(ctx, id.FileSystemName, "/", payload)
		if err == nil {
			acl, err := accesscontrol.ParseACL(pathResponse.ACL)
			if err != nil {
				return fmt.Errorf("parsing response ACL %q: %s", pathResponse.ACL, err)
			}
			ace = helpers.FlattenDataLakeGen2AceList(d, acl)
			owner = pathResponse.Owner
			group = pathResponse.Group
		}
	}
	d.Set("ace", ace)
	d.Set("owner", owner)
	d.Set("group", group)

	return nil
}
