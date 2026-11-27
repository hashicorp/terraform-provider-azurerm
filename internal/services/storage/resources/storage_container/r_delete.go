// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_container

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageContainerDelete(d *pluginsdk.ResourceData, meta any) error {
	containerClient := client.From(meta).ResourceManager.BlobContainers
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseStorageContainerID(d.Id())
	if err != nil {
		return err
	}

	if _, err := containerClient.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %v", d.Id(), err)
	}

	return nil
}
