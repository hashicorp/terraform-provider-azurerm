// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileshares"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
)

func resourceStorageShareUpdate(d *pluginsdk.ResourceData, meta any) error {
	sharesClient := client.From(meta).ResourceManager.FileShares
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := fileshares.ParseShareID(d.Id())
	if err != nil {
		return err
	}

	update := fileshares.FileShare{
		Properties: &fileshares.FileShareProperties{},
	}

	if d.HasChange("quota") {
		quota := d.Get("quota").(int)
		update.Properties.ShareQuota = pointer.To(int64(quota))
	}

	if d.HasChange("metadata") {
		metaDataRaw := d.Get("metadata").(map[string]any)
		metaData := helpers.ExpandMetaData(metaDataRaw)

		update.Properties.Metadata = pointer.To(metaData)
	}

	if d.HasChange("acl") {
		update.Properties.SignedIdentifiers = expandStorageShareACLs(d.Get("acl").(*pluginsdk.Set).List())
	}

	if d.HasChange("access_tier") {
		tier := shares.AccessTier(d.Get("access_tier").(string))
		update.Properties.AccessTier = pointer.ToEnum[fileshares.ShareAccessTier](string(tier))
	}

	if _, err = sharesClient.Update(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %v", id, err)
	}

	return resourceStorageShareRead(d, meta)
}
