// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/backups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/backupvaults"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetAppBackupVaultResource struct{}

var _ sdk.Resource = NetAppBackupVaultResource{}

func (r NetAppBackupVaultResource) ModelObject() interface{} {
	return &netAppModels.NetAppBackupVaultModel{}
}

func (r NetAppBackupVaultResource) ResourceType() string {
	return "azurerm_netapp_backup_vault"
}

func (r NetAppBackupVaultResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupvaults.ValidateBackupVaultID
}

func (r NetAppBackupVaultResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.VolumeQuotaRuleName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.AccountName,
		},

		"tags": commonschema.Tags(),
	}
}

func (r NetAppBackupVaultResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetAppBackupVaultResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BackupVaultsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model netAppModels.NetAppBackupVaultModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := backupvaults.NewBackupVaultID(subscriptionId, model.ResourceGroupName, model.AccountName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			parameters := backupvaults.BackupVault{
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
			}

			err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NetAppBackupVaultResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BackupVaultsClient

			id, err := backupvaults.ParseBackupVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppBackupVaultModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") {
				metadata.Logger.Infof("Updating %s", id)

				update := backupvaults.BackupVaultPatch{
					Tags: pointer.To(state.Tags),
				}

				if err := client.UpdateThenPoll(ctx, pointer.From(id), update); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r NetAppBackupVaultResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BackupVaultsClient

			id, err := backupvaults.ParseBackupVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppBackupVaultModel
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
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     5 * time.Second,
		MinTimeout:                5 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"404"},
		Refresh:                   netappBackupStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func netappBackupStateRefreshFunc(ctx context.Context, client *backups.BackupsClient, id backups.BackupId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving %s: %s", id, err)
			}
		}

		statusCode := "dropped connection"
		if res.HttpResponse != nil {
			statusCode = strconv.Itoa(res.HttpResponse.StatusCode)
		}

		return res, statusCode, nil
	}
}

func (r NetAppBackupVaultResource) Delete() sdk.ResourceFunc {
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
			for retries := 0; retries < 5; retries++ {
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

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for deletion of %s: %w", id, err)
	}

	return nil
}

func netappBackupVaultStateRefreshFunc(ctx context.Context, vaultClient *backupvaults.BackupVaultsClient, backupClient *backups.BackupsClient, id backupvaults.BackupVaultId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
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
	for attempt := 0; attempt < retryAttempts; attempt++ {
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
