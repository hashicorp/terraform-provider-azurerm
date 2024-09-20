// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/blob/blobs"
)

func dataSourceStorageBlobContent() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStorageBlobContentRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"storage_account_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"storage_container_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"content": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStorageBlobContentRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	input := blobs.GetInput{}
	content, err := blobsClient.Get(ctx, containerName, name, input)
	if err != nil {
		if response.WasNotFound(content.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving content for %s: %v", id, err)
	}

	d.Set("content", string(*content.Contents))

	d.Set("name", name)
	d.Set("storage_container_name", containerName)
	d.Set("storage_account_name", accountName)

	d.SetId(id.ID())

	return nil
}
