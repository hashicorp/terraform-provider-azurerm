// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservicesresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageSyncUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := storageClient.From(meta).SyncServiceClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := storagesyncservicesresource.ParseStorageSyncServiceID(d.Id())
	if err != nil {
		return err
	}

	update := storagesyncservicesresource.StorageSyncServiceUpdateParameters{}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("incoming_traffic_policy") {
		update.Properties = &storagesyncservicesresource.StorageSyncServiceUpdateProperties{
			IncomingTrafficPolicy: pointer.ToEnum[storagesyncservicesresource.IncomingTrafficPolicy](d.Get("incoming_traffic_policy").(string)),
		}
	}

	if err = client.StorageSyncServicesUpdateThenPoll(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceStorageSyncRead(d, meta)
}
