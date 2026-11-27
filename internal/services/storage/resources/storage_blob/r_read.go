// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_blob

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

func resourceStorageBlobRead(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := blobs.ParseBlobID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing %q: %v", d.Id(), err)
	}

	var accountName string
	var containerName string
	var account *client.AccountDetails

	containerIdStr := d.Get("storage_container_id").(string)
	if containerIdStr == "" { // this may be missing at import time
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
		log.Printf("[DEBUG] Unable to locate Account %q for Blob %q (Container %q) - assuming removed & removing from state!", accountName, id.BlobName, containerName)
		d.SetId("")
		return nil
	}

	blobsClient, err := storageClient.BlobsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Blobs Client: %v", err)
	}

	input := blobs.GetPropertiesInput{}
	props, err := blobsClient.GetProperties(ctx, id.ContainerName, id.BlobName, input)
	if err != nil {
		if response.WasNotFound(props.HttpResponse) {
			log.Printf("[INFO] Blob %q was not found in Container %q / Account %q - assuming removed & removing from state...", id.BlobName, id.ContainerName, id.AccountId.AccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving properties for %s: %v", id, err)
	}

	d.Set("name", id.BlobName)
	d.Set("storage_container_id", commonids.NewStorageContainerID(subscriptionId, account.StorageAccountId.ResourceGroupName, accountName, containerName).ID())

	d.Set("access_tier", string(props.AccessTier))
	d.Set("content_type", props.ContentType)
	d.Set("cache_control", props.CacheControl)

	// Set the ContentMD5 value to md5 hash in hex
	contentMD5 := ""
	if props.ContentMD5 != "" {
		contentMD5, err = helpers.ConvertBase64ToHexEncoding(props.ContentMD5)
		if err != nil {
			return fmt.Errorf("converting hex to base64 encoding for content_md5: %v", err)
		}
	}
	d.Set("content_md5", contentMD5)

	d.Set("encryption_scope", props.EncryptionScope)

	d.Set("type", strings.TrimSuffix(string(props.BlobType), "Blob"))
	d.Set("url", d.Id())

	if err = d.Set("metadata", helpers.FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("setting `metadata`: %v", err)
	}
	// The CopySource is only returned if the blob hasn't been modified (e.g. metadata configured etc)
	// as such, we need to conditionally set this to ensure it's trackable if possible
	if props.CopySource != "" {
		d.Set("source_uri", props.CopySource)
	}

	return nil
}
