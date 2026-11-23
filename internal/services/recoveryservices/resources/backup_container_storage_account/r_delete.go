// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_container_storage_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2025-08-01/protectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectionContainerStorageAccountDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.BackupProtectionContainersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protectioncontainers.ParseProtectionContainerID(d.Id())
	if err != nil {
		return err
	}

	if err = client.UnregisterThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("unregistering %s: %+v", id, err)
	}

	return nil
}
