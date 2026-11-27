// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_blob

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

func resourceStorageBlobDelete(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := blobs.ParseBlobID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing %q: %v", d.Id(), err)
	}

	var accountName string
	var containerName string
	var account *client.AccountDetails

	containerIdStr := d.Get("storage_container_id").(string) // lintignore: gradually-deprecated
	if containerIdStr == "" {
		accountName = id.AccountId.AccountName
		containerName = id.ContainerName

		account, err = storageClient.FindAccount(ctx, subscriptionId, accountName)
		if err != nil {
			return fmt.Errorf("retrieving Account %q for Blob %q (Container %q): %v", accountName, id.BlobName, containerName, err)
		}
	} else {
		containerId, err := commonids.ParseStorageContainerID(containerIdStr)
		if err != nil {
			return err
		}
		accountName = containerId.StorageAccountName
		containerName = containerId.ContainerName

		account, err = storageClient.GetAccount(ctx, commonids.NewStorageAccountID(containerId.SubscriptionId, containerId.ResourceGroupName, containerId.StorageAccountName))
		if err != nil {
			return fmt.Errorf("retrieving Account %q for Blob %q (Container %q): %v", accountName, id.BlobName, containerName, err)
		}
	}

	if account == nil {
		return fmt.Errorf("locating Storage Account %q", accountName)
	}

	blobsClient, err := storageClient.BlobsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Blobs Client: %v", err)
	}

	input := blobs.DeleteInput{
		DeleteSnapshots: true,
	}
	if resp, err := blobsClient.Delete(ctx, id.ContainerName, id.BlobName, input); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting %s: %v", id, err)
	}

	return nil
}
