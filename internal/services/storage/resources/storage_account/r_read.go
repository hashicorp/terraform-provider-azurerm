// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageAccountRead(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage.ResourceManager
	client := storageClient.StorageAccounts
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := resourceStorageAccountFlatten(ctx, d, *id, resp.Model, meta, true); err != nil {
		return fmt.Errorf("encoding %s: %+v", *id, err)
	}

	return nil
}
