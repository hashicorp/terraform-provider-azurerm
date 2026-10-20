// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_backup_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/backuppolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BackupPolicyClient

			id, err := backuppolicies.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state models.NetAppBackupPolicyModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			update := backuppolicies.BackupPolicyPatch{
				Properties: &backuppolicies.BackupPolicyProperties{},
			}

			// Checking properties with changes
			if metadata.ResourceData.HasChange("tags") {
				update.Tags = pointer.To(state.Tags)
			}

			if metadata.ResourceData.HasChange("daily_backups_to_keep") {
				update.Properties.DailyBackupsToKeep = pointer.To(state.DailyBackupsToKeep)
			}

			if metadata.ResourceData.HasChange("weekly_backups_to_keep") {
				update.Properties.WeeklyBackupsToKeep = pointer.To(state.WeeklyBackupsToKeep)
			}

			if metadata.ResourceData.HasChange("monthly_backups_to_keep") {
				update.Properties.MonthlyBackupsToKeep = pointer.To(state.MonthlyBackupsToKeep)
			}

			if metadata.ResourceData.HasChange("enabled") {
				update.Properties.Enabled = pointer.To(state.Enabled)
			}

			if err := client.UpdateThenPoll(ctx, pointer.From(id), update); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
