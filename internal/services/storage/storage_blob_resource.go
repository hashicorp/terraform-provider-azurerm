// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

func resourceStorageBlob() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageBlobCreate,
		Read:   resourceStorageBlobRead,
		Update: resourceStorageBlobUpdate,
		Delete: resourceStorageBlobDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.BlobV0ToV1{},
		}),

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			_, err := blobs.ParseBlobID(id, storageDomainSuffix)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountName,
			},

			"storage_container_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageContainerName,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Append",
					"Block",
					"Page",
				}, false),
			},

			"size": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      0,
				ValidateFunc: validation.IntDivisibleBy(512),
			},

			"access_tier": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(blobs.Archive),
					string(blobs.Cool),
					string(blobs.Hot),
				}, false),
			},

			"content_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "application/octet-stream",
			},

			"cache_control": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"encryption_scope": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageEncryptionScopeName,
			},

			"source": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_uri", "source_content"},
			},

			"source_content": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source", "source_uri"},
			},

			"source_uri": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source", "source_content"},
			},

			"content_md5": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_uri"},
			},

			"url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"parallelism": {
				// TODO: @tombuildsstuff - a note this only works for Page blobs
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      8,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"metadata": MetaDataComputedSchema(),
		},

		CustomizeDiff: func(ctx context.Context, diff *pluginsdk.ResourceDiff, i interface{}) error {
			if content := diff.Get("source_content"); content != "" && diff.Get("type") == "Page" {
				if len(content.(string))%512 != 0 {
					return fmt.Errorf(`"source" must be aligned to 512-byte boundary for "type" set to "Page"`)
				}
			}
			return nil
		},
	}
}

func resourceStorageBlobCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	accountName := d.Get("storage_account_name").(string)
	containerName := d.Get("storage_container_name").(string)
	name := d.Get("name").(string)

	account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
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
	if d.IsNewResource() {
		input := blobs.GetPropertiesInput{}
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
		contentMD5, err = convertHexToBase64Encoding(d.Get("content_md5").(string))
		if err != nil {
			return fmt.Errorf("failed to base64 encode `content_md5` value: %s", err)
		}
	}

	log.Printf("[DEBUG] Creating %s..", id)
	metaDataRaw := d.Get("metadata").(map[string]interface{})
	blobInput := BlobUpload{
		AccountName:   accountName,
		ContainerName: containerName,
		BlobName:      name,
		Client:        blobsClient,

		BlobType:      d.Get("type").(string),
		CacheControl:  d.Get("cache_control").(string),
		ContentType:   d.Get("content_type").(string),
		ContentMD5:    contentMD5,
		MetaData:      ExpandMetaData(metaDataRaw),
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
	log.Printf("[DEBUG] Created %s.", id)

	d.SetId(id.ID())

	return resourceStorageBlobUpdate(d, meta)
}

func resourceStorageBlobUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := blobs.ParseBlobID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing %q: %v", d.Id(), err)
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Blob %q (Container %q): %v", id.AccountId.AccountName, id.BlobName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
	}

	blobsClient, err := storageClient.BlobsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Blobs Client: %v", err)
	}

	log.Printf("[INFO] Retrieving %s", id)
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

	if d.HasChange("content_type") || d.HasChange("cache_control") {
		log.Printf("[DEBUG] Updating Properties for %s...", id)
		input := blobs.SetPropertiesInput{
			ContentType:  pointer.To(d.Get("content_type").(string)),
			CacheControl: pointer.To(d.Get("cache_control").(string)),
		}

		// `content_md5` is `ForceNew` but must be included in the `SetPropertiesInput` update payload, or it will be zeroed on the blob.
		if contentMD5 := d.Get("content_md5").(string); contentMD5 != "" {
			data, err := convertHexToBase64Encoding(contentMD5)
			if err != nil {
				return fmt.Errorf("converting hex to base64 encoding for content_md5: %v", err)
			}
			input.ContentMD5 = pointer.To(data)
		}
		if _, err = blobsClient.SetProperties(ctx, id.ContainerName, id.BlobName, input); err != nil {
			return fmt.Errorf("updating Properties for %s: %v", id, err)
		}
		log.Printf("[DEBUG] Updated Properties for %s", id)
	}

	if d.HasChange("metadata") {
		log.Printf("[DEBUG] Updating MetaData for %s...", id)
		metaDataRaw := d.Get("metadata").(map[string]interface{})
		input := blobs.SetMetaDataInput{
			MetaData: ExpandMetaData(metaDataRaw),
		}
		// Encryption Scope must be specified when updating metadata
		if props.EncryptionScope != "" {
			input.EncryptionScope = pointer.To(props.EncryptionScope)
		}
		if _, err = blobsClient.SetMetaData(ctx, id.ContainerName, id.BlobName, input); err != nil {
			return fmt.Errorf("updating MetaData for %s: %v", id, err)
		}
		log.Printf("[DEBUG] Updated MetaData for %s", id)
	}

	if d.HasChange("access_tier") {
		// this is only applicable for Gen2/BlobStorage accounts
		log.Printf("[DEBUG] Updating Access Tier for %s...", id)
		accessTier := blobs.AccessTier(d.Get("access_tier").(string))

		if _, err := blobsClient.SetTier(ctx, id.ContainerName, id.BlobName, blobs.SetTierInput{Tier: accessTier}); err != nil {
			return fmt.Errorf("updating Access Tier for %s: %v", id, err)
		}

		log.Printf("[DEBUG] Updated Access Tier for %s", id)
	}

	return resourceStorageBlobRead(d, meta)
}

func resourceStorageBlobRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := blobs.ParseBlobID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing %q: %v", d.Id(), err)
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Blob %q (Container %q): %v", id.AccountId.AccountName, id.BlobName, id.ContainerName, err)
	}
	if account == nil {
		log.Printf("[DEBUG] Unable to locate Account %q for Blob %q (Container %q) - assuming removed & removing from state!", id.AccountId.AccountName, id.BlobName, id.ContainerName)
		d.SetId("")
		return nil
	}

	blobsClient, err := storageClient.BlobsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Blobs Client: %v", err)
	}

	log.Printf("[INFO] Retrieving %s", id)
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
	d.Set("storage_container_name", id.ContainerName)
	d.Set("storage_account_name", id.AccountId.AccountName)

	d.Set("access_tier", string(props.AccessTier))
	d.Set("content_type", props.ContentType)
	d.Set("cache_control", props.CacheControl)

	// Set the ContentMD5 value to md5 hash in hex
	contentMD5 := ""
	if props.ContentMD5 != "" {
		contentMD5, err = convertBase64ToHexEncoding(props.ContentMD5)
		if err != nil {
			return fmt.Errorf("converting hex to base64 encoding for content_md5: %v", err)
		}
	}
	d.Set("content_md5", contentMD5)

	d.Set("encryption_scope", props.EncryptionScope)

	d.Set("type", strings.TrimSuffix(string(props.BlobType), "Blob"))
	d.Set("url", d.Id())

	if err = d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("setting `metadata`: %v", err)
	}
	// The CopySource is only returned if the blob hasn't been modified (e.g. metadata configured etc)
	// as such, we need to conditionally set this to ensure it's trackable if possible
	if props.CopySource != "" {
		d.Set("source_uri", props.CopySource)
	}

	return nil
}

func resourceStorageBlobDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := blobs.ParseBlobID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing %q: %v", d.Id(), err)
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Blob %q (Container %q): %s", id.AccountId.AccountName, id.BlobName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
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
