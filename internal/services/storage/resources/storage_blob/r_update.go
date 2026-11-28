// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_blob

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

func resourceStorageBlobUpdate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := blobs.ParseBlobID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	var accountName string
	var containerName string
	var account *client.AccountDetails

	containerId, err := commonids.ParseStorageContainerID(d.Get("storage_container_id").(string))
	if err != nil {
		return err
	}
	accountName = containerId.StorageAccountName
	containerName = containerId.ContainerName

	account, err = storageClient.GetAccount(ctx, commonids.NewStorageAccountID(containerId.SubscriptionId, containerId.ResourceGroupName, containerId.StorageAccountName))
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Blob %q (Container %q): %v", accountName, id.BlobName, containerName, err)
	}

	if account == nil {
		return fmt.Errorf("locating Storage Account %q", accountName)
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

	if d.HasChanges("content_type", "cache_control") {
		input := blobs.SetPropertiesInput{
			ContentType:  pointer.To(d.Get("content_type").(string)),
			CacheControl: pointer.To(d.Get("cache_control").(string)),
		}

		// `content_md5` is `ForceNew` but must be included in the `SetPropertiesInput` update payload, or it will be zeroed on the blob.
		if contentMD5 := d.Get("content_md5").(string); contentMD5 != "" {
			data, err := helpers.ConvertHexToBase64Encoding(contentMD5)
			if err != nil {
				return fmt.Errorf("converting hex to base64 encoding for content_md5: %v", err)
			}
			input.ContentMD5 = pointer.To(data)
		}
		if _, err = blobsClient.SetProperties(ctx, id.ContainerName, id.BlobName, input); err != nil {
			return fmt.Errorf("updating Properties for %s: %v", id, err)
		}
	}

	if d.HasChange("metadata") {
		metaDataRaw := d.Get("metadata").(map[string]any)
		input := blobs.SetMetaDataInput{
			MetaData: helpers.ExpandMetaData(metaDataRaw),
		}
		// Encryption Scope must be specified when updating metadata
		if props.EncryptionScope != "" {
			input.EncryptionScope = pointer.To(props.EncryptionScope)
		}
		if _, err = blobsClient.SetMetaData(ctx, id.ContainerName, id.BlobName, input); err != nil {
			return fmt.Errorf("updating MetaData for %s: %v", id, err)
		}
	}

	if d.HasChange("access_tier") {
		// this is only applicable for Gen2/BlobStorage accounts
		accessTier := blobs.AccessTier(d.Get("access_tier").(string))

		if _, err := blobsClient.SetTier(ctx, id.ContainerName, id.BlobName, blobs.SetTierInput{Tier: accessTier}); err != nil {
			return fmt.Errorf("updating Access Tier for %s: %v", id, err)
		}
	}

	return resourceStorageBlobRead(d, meta)
}
