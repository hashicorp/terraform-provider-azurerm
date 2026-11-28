// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_blob

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

func resourceStorageBlobCreate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	var accountName string
	var containerName string
	var account *client.AccountDetails
	var err error

	containerIdStr := d.Get("storage_container_id").(string)
	containerId, err := commonids.ParseStorageContainerID(containerIdStr)
	if err != nil {
		return err
	}
	accountName = containerId.StorageAccountName
	containerName = containerId.ContainerName
	account, err = storageClient.GetAccount(ctx, commonids.NewStorageAccountID(containerId.SubscriptionId, containerId.ResourceGroupName, containerId.StorageAccountName))
	if err != nil {
		return fmt.Errorf("retrieving Storage Account %q for Blob %q (Container %q): %v", accountName, name, containerName, err)
	}

	if account == nil {
		return fmt.Errorf("locating Storage Account %q", accountName)
	}

	blobsClient, err := storageClient.BlobsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Blobs Client: %v", err)
	}

	accountId := accounts.AccountId{
		AccountName:   accountName,
		DomainSuffix:  storageClient.StorageDomainSuffix,
		SubDomainType: accounts.BlobSubDomainType,
	}

	id := blobs.NewBlobID(accountId, containerName, name)
	input := blobs.GetPropertiesInput{}
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		props, err := blobsClient.GetProperties(ctx, containerName, name, input)
		if err != nil {
			if !response.WasNotFound(props.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %v", id, err)
			}
		}
		if !response.WasNotFound(props.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_storage_blob", id.ID())
		}
	}

	contentMD5Raw := d.Get("content_md5").(string)
	contentMD5 := ""
	if contentMD5Raw != "" {
		// Azure uses a Base64 encoded representation of the standard MD5 sum of the file
		contentMD5, err = helpers.ConvertHexToBase64Encoding(d.Get("content_md5").(string))
		if err != nil {
			return fmt.Errorf("failed to base64 encode `content_md5` value: %s", err)
		}
	}

	metaDataRaw := d.Get("metadata").(map[string]any)
	blobInput := BlobUpload{
		AccountName:   accountName,
		ContainerName: containerName,
		BlobName:      name,
		Client:        blobsClient,
		BlobType:      d.Get("type").(string),
		CacheControl:  d.Get("cache_control").(string),
		ContentType:   d.Get("content_type").(string),
		ContentMD5:    contentMD5,
		MetaData:      helpers.ExpandMetaData(metaDataRaw),
		Parallelism:   d.Get("parallelism").(int),
		Size:          d.Get("size").(int),
		Source:        d.Get("source").(string),
		SourceContent: d.Get("source_content").(string),
		SourceUri:     d.Get("source_uri").(string),
	}

	if encryptionScope := d.Get("encryption_scope"); encryptionScope.(string) != "" {
		blobInput.EncryptionScope = encryptionScope.(string)
	}

	if err = blobInput.Create(ctx); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageBlobUpdate(d, meta)
}
