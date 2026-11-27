// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_table

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageTableUpdate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rmId, err := tables.ParseTableID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, subscriptionId, rmId.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", rmId.StorageAccountName, rmId.TableName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", rmId.StorageAccountName)
	}

	if d.HasChange("acl") {
		aclsRaw := d.Get("acl").(*pluginsdk.Set).List()
		acls := expandStorageTableACLs(aclsRaw)

		aclClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Tables Client: %v", err)
		}

		if err = aclClient.UpdateACLs(ctx, rmId.TableName, acls); err != nil {
			return fmt.Errorf("updating ACLs for table %q: %v", rmId.TableName, err)
		}
	}

	return resourceStorageTableRead(d, meta)
}
