// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_backup_vault

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/backupvaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			vaultClient := metadata.Client.NetApp.BackupVaultsClient
			backupClient := metadata.Client.NetApp.BackupClient

			id, err := backupvaults.ParseBackupVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Attempt to delete backup vault with retries
			for range 5 {
				// Delete backups
				if err := deleteBackupsFromVault(ctx, id, backupClient, metadata.Client.Features.NetApp.DeleteBackupsOnBackupVaultDestroy); err != nil {
					return err
				}

				// DeleteThenPoll cannot be used due to potential race condition where the backup started when the volume got deleted but it takes time for it to show up within the vault
				// This will be handled by waitForBackupVaultDeletion and operation will be retried if needed
				if _, err := vaultClient.Delete(ctx, pointer.From(id)); err != nil {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}

				// Wait for deletion to complete
				err = waitForBackupVaultDeletion(ctx, vaultClient, backupClient, pointer.From(id))
				if err == nil {
					return nil // Successful deletion
				}

				if strings.Contains(err.Error(), "backups found on vault") {
					// Backup may not show up in the vault through a GET so we will wait for a bit before retrying
					time.Sleep(30 * time.Second)
					continue
				}

				return err // If it's a different error, return it immediately
			}

			return fmt.Errorf("failed to delete backup vault after 5 attempts")
		},
	}
}
