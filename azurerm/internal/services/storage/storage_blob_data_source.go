package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/blobs"
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
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	accountName := d.Get("storage_account_name").(string)
	containerName := d.Get("storage_container_name").(string)
	name := d.Get("name").(string)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Blob %q (Container %q): %s", accountName, name, containerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
	}

	blobsClient, err := storageClient.BlobsClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Blobs Client: %s", err)
	}

	id := blobsClient.GetResourceID(accountName, containerName, name)

	log.Printf("[INFO] Retrieving Storage Blob %q (Container %q / Account %q).", name, containerName, accountName)
	input := blobs.GetPropertiesInput{}
	props, err := blobsClient.GetProperties(ctx, accountName, containerName, name, input)
	if err != nil {
		if utils.ResponseWasNotFound(props.Response) {
			log.Printf("[INFO] Blob %q was not found in Container %q / Account %q - assuming removed & removing from state...", name, containerName, accountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving properties for Blob %q (Container %q / Account %q): %s", name, containerName, accountName, err)
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
			return fmt.Errorf("Error in converting hex to base64 encoding for content_md5: %s", err)
		}
	}
	d.Set("content_md5", contentMD5)

	d.Set("type", strings.TrimSuffix(string(props.BlobType), "Blob"))

	d.SetId(id)

	d.Set("url", id)

	if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("Error setting `metadata`: %+v", err)
	}

	return nil
}
