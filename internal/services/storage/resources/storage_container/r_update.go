// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_container

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/blobcontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageContainerUpdate(d *pluginsdk.ResourceData, meta any) error {
	containerClient := meta.(*clients.Client).Storage.ResourceManager.BlobContainers
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageContainerID(d.Id())
	if err != nil {
		return err
	}

	update := blobcontainers.BlobContainer{
		Properties: &blobcontainers.ContainerProperties{},
	}

	if d.HasChange("container_access_type") {
		accessLevelRaw := d.Get("container_access_type").(string)
		update.Properties.PublicAccess = pointer.ToEnum[blobcontainers.PublicAccess](containerAccessTypeConversionMap[accessLevelRaw])
	}

	if d.HasChange("metadata") {
		metaDataRaw := d.Get("metadata").(map[string]any)
		update.Properties.Metadata = pointer.To(helpers.ExpandMetaData(metaDataRaw))
	}

	if _, err := containerClient.Update(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %v", id, err)
	}

	return resourceStorageContainerRead(d, meta)
}
