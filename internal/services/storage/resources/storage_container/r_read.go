// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_container

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageContainerRead(d *pluginsdk.ResourceData, meta any) error {
	containerClient := meta.(*clients.Client).Storage.ResourceManager.BlobContainers
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageContainerID(d.Id())
	if err != nil {
		return err
	}

	existing, err := containerClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			log.Printf("[DEBUG] %q was not found, removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return resourceStorageContainerFlatten(ctx, d, id, existing.Model, meta, true)
}

var containerAccessTypeConversionMap = map[string]string{
	"blob":      "Blob",
	"container": "Container",
	"private":   "None",
	"Blob":      "blob",
	"Container": "container",
	"None":      "private",
	"":          "private",
}
