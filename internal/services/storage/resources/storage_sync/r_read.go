// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservicesresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageSyncRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Storage.SyncServiceClient
	registeredServerClient := meta.(*clients.Client).Storage.SyncRegisteredServerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := storagesyncservicesresource.ParseStorageSyncServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.StorageSyncServicesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceStorageSyncFlatten(ctx, d, id, resp.Model, registeredServerClient, true)
}
