// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_container_storage_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2025-08-01/protectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectionContainerStorageAccountRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.BackupProtectionContainersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protectioncontainers.ParseProtectionContainerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on backup protection container %s : %+v", id.String(), err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	if model := resp.Model; model != nil {
		if properties, ok := model.Properties.(protectioncontainers.AzureStorageContainer); ok {
			d.Set("storage_account_id", properties.SourceResourceId)
		}
	}

	return nil
}
