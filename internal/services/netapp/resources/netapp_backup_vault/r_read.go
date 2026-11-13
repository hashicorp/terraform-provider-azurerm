// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_backup_vault

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/backups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/backupvaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BackupVaultsClient

			id, err := backupvaults.ParseBackupVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state models.NetAppBackupVaultModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			state.AccountName = id.NetAppAccountName
			state.Name = id.BackupVaultName
			state.ResourceGroupName = id.ResourceGroupName

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func waitForBackupDeletion(ctx context.Context, client *backups.BackupsClient, id backups.BackupId) error {
	poller := custompollers.NewEventualConsistencyPoller(5, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, id)
		return resp.HttpResponse, err
	}, &custompollers.EventualConsistencyPollerOptions{
		Interval:         5 * time.Second,
		TargetStatusCode: pointer.To(http.StatusNotFound),
	})
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}
	return nil
}

func waitForBackupVaultDeletion(ctx context.Context, vaultClient *backupvaults.BackupVaultsClient, backupClient *backups.BackupsClient, id backupvaults.BackupVaultId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"404"},
		Refresh:                   netappBackupVaultStateRefreshFunc(ctx, vaultClient, backupClient, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %w", id, err)
	}

	return nil
}

func netappBackupVaultStateRefreshFunc(ctx context.Context, vaultClient *backupvaults.BackupVaultsClient, backupClient *backups.BackupsClient, id backupvaults.BackupVaultId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := vaultClient.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
			}
			return nil, "", fmt.Errorf("retrieving %s: %s", id, err)
		}

		// Handling a race condition where the backup started but the volume got deleted and it takes time for it to show up within the vault.
		// The vault deletion process will hang until the deadline is reached and never retry to delete the backup preventing vault to be deleted.
		// For this to work, need to scrub activity logs to see if there was a failed deletion operation due to backup just showing up in the vault
		// midway after the vault deletion process started.
		backupVaultID := backups.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.BackupVaultName)
		backupList, err := backupClient.ListByVault(ctx, backupVaultID, backups.ListByVaultOperationOptions{})
		if err != nil {
			return nil, "", fmt.Errorf("listing backups from vault %s: %w", id.ID(), err)
		}

		if backupList.Model != nil || len(pointer.From(backupList.Model)) > 0 {
			return nil, "409", fmt.Errorf("backups found on vault %s, forcing retry", id.ID()) // Forcing vault deletion to retry
		}

		return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
	}
}

func deleteBackupsFromVault(ctx context.Context, id *backupvaults.BackupVaultId, backupClient *backups.BackupsClient, shouldDestroyBackups bool) error {
	for {
		backupVaultID := backups.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.BackupVaultName)
		backupList, err := backupClient.ListByVault(ctx, backupVaultID, backups.ListByVaultOperationOptions{})
		if err != nil {
			return fmt.Errorf("listing backups from vault %s: %w", id.ID(), err)
		}

		if backupList.Model == nil || len(pointer.From(backupList.Model)) == 0 {
			return nil // No more backups to delete
		}

		for _, backup := range pointer.From(backupList.Model) {
			if backup.Name == nil {
				continue
			}

			if shouldDestroyBackups {
				backupID, err := backups.ParseBackupID(pointer.From(backup.Id))
				if err != nil {
					return fmt.Errorf("parsing backup ID %s: %w", pointer.From(backup.Id), err)
				}

				if err := retryBackupDelete(ctx, backupClient, pointer.From(backupID), 120, 30); err != nil {
					return fmt.Errorf("failed to delete backup %s: %w", backupID.ID(), err)
				}
			} else {
				return fmt.Errorf("cannot delete backups from backup vault due to missing DeleteBackupsOnBackupVaultDestroy feature set as true, backup vault id %s, DeleteBackupsOnBackupVaultDestroy setting is: %t", id.ID(), shouldDestroyBackups)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func retryBackupDelete(ctx context.Context, client *backups.BackupsClient, id backups.BackupId, retryAttempts, retryIntervalSec int) error {
	var lastErr error
	for range retryAttempts {
		if err := client.DeleteThenPoll(ctx, id); err == nil {
			if err := waitForBackupDeletion(ctx, client, id); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %w", id.ID(), err)
			}
			return nil
		} else if strings.Contains(err.Error(), "Please retry after backup transfer is complete") {
			lastErr = err
			time.Sleep(time.Duration(retryIntervalSec) * time.Second)
		} else {
			return fmt.Errorf("deleting backup %s: %+v", id.ID(), err)
		}
	}

	return fmt.Errorf("failed to delete backup after %d attempts: %v", retryAttempts, lastErr)
}
