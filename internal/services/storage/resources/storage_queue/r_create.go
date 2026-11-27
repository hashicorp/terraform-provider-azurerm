// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_queue

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storagequeues"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageQueueCreate(d *pluginsdk.ResourceData, meta any) error {
	queueClient := client.From(meta).ResourceManager.StorageQueues
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	queueName := d.Get("name").(string)

	metaDataRaw := d.Get("metadata").(map[string]any)
	metaData := helpers.ExpandMetaData(metaDataRaw)

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := storagequeues.NewQueueID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, queueName)

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := queueClient.QueueGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %q: %v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_storage_queue", id.ID())
		}
	}

	payload := storagequeues.StorageQueue{
		Properties: &storagequeues.QueueProperties{
			Metadata: &metaData,
		},
	}

	if _, err := queueClient.QueueCreate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageQueueRead(d, meta)
}
