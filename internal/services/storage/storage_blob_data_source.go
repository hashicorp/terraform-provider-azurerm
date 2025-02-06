// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

func dataSourceStorageBlob() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStorageBlobRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				// TODO: add validation
			},

			"storage_account_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"storage_container_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"access_tier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"content_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"content_md5": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"encryption_scope": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"metadata": MetaDataComputedSchema(),
		},
	}
}

func dataSourceStorageBlobRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	accountName := d.Get("storage_account_name").(string)
	containerName := d.Get("storage_container_name").(string)
	name := d.Get("name").(string)

	account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Blob %q (Container %q): %v", accountName, name, containerName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", accountName)
	}

	blobsClient, err := storageClient.BlobsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Blobs Client: %v", err)
	}

	blobEndpoint, err := account.DataPlaneEndpoint(client.EndpointTypeBlob)
	if err != nil {
		return fmt.Errorf("retrieving the blob data plane endpoint: %v", err)
	}

	accountId, err := accounts.ParseAccountID(*blobEndpoint, storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	id := blobs.NewBlobID(*accountId, containerName, name)

	log.Printf("[INFO] Retrieving %s", id)
	input := blobs.GetPropertiesInput{}
	props, err := blobsClient.GetProperties(ctx, containerName, name, input)
	if err != nil {
		if response.WasNotFound(props.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving properties for %s: %v", id, err)
	}

	d.Set("name", name)
	d.Set("storage_container_name", containerName)
	d.Set("storage_account_name", accountName)

	d.Set("access_tier", string(props.AccessTier))
	d.Set("content_type", props.ContentType)

	// Set the ContentMD5 value to md5 hash in hex
	contentMD5 := ""
	if props.ContentMD5 != "" {
		contentMD5, err = convertBase64ToHexEncoding(props.ContentMD5)
		if err != nil {
			return fmt.Errorf("in converting hex to base64 encoding for content_md5: %s", err)
		}
	}
	d.Set("content_md5", contentMD5)

	d.Set("encryption_scope", props.EncryptionScope)

	d.Set("type", strings.TrimSuffix(string(props.BlobType), "Blob"))

	d.SetId(id.ID())

	d.Set("url", id.ID())

	if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("setting `metadata`: %+v", err)
	}

	return nil
}
