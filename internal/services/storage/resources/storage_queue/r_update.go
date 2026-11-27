// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_queue

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storagequeues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageQueueUpdate(d *pluginsdk.ResourceData, meta any) error {
	queueClient := client.From(meta).ResourceManager.StorageQueues
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := storagequeues.ParseQueueID(d.Id())
	if err != nil {
		return err
	}

	existing, err := queueClient.QueueGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %q: %v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("unexpected null model after retrieving %v", id)
	}

	payload := storagequeues.StorageQueue{
		Properties: existing.Model.Properties,
	}

	if d.HasChange("metadata") {
		metaDataRaw := d.Get("metadata").(map[string]any)
		payload.Properties.Metadata = pointer.To(helpers.ExpandMetaData(metaDataRaw))
	}

	if _, err := queueClient.QueueCreate(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %v", id, err)
	}

	return resourceStorageQueueRead(d, meta)
}
