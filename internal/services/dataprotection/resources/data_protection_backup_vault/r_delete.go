// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataProtectionBackupVaultDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupvaultresources.ParseBackupVaultID(d.Id())
	if err != nil {
		return err
	}

	if err := client.BackupVaultsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting DataProtection BackupVault (%q): %+v", id, err)
	}

	// API has bug, which appears API returns before the resource is fully deleted. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/38944
	pollerType := custompollers.NewDataProtectionBackupVaultPoller(client, *id)
	poller := pollers.NewPoller(pollerType, 30*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	return nil
}
