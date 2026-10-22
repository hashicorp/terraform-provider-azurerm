// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			var model AutonomousDatabaseRegularResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			if _, err = client.Get(ctx, *id); err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			// Check what needs to be updated
			needsGeneralUpdate := r.hasGeneralUpdates(metadata)
			needsBackupScheduleUpdate := metadata.ResourceData.HasChange("long_term_backup_schedule")
			needsPasswordUpdate := metadata.ResourceData.HasChange("admin_password")
			needsBackupRetentionDaysUpdate := metadata.ResourceData.HasChange("backup_retention_period_in_days")

			// Step 1: Handle general updates (everything except backup schedule)
			if needsGeneralUpdate {
				generalUpdate := autonomousdatabases.AutonomousDatabaseUpdate{
					Properties: &autonomousdatabases.AutonomousDatabaseUpdateProperties{},
				}

				if metadata.ResourceData.HasChange("tags") {
					generalUpdate.Tags = pointer.To(model.Tags)
				}
				if metadata.ResourceData.HasChange("data_storage_size_in_tbs") {
					generalUpdate.Properties.DataStorageSizeInTbs = pointer.To(model.DataStorageSizeInTbs)
				}
				if metadata.ResourceData.HasChange("compute_count") {
					generalUpdate.Properties.ComputeCount = pointer.To(model.ComputeCount)
				}
				if metadata.ResourceData.HasChange("auto_scaling_enabled") {
					generalUpdate.Properties.IsAutoScalingEnabled = pointer.To(model.AutoScalingEnabled)
				}
				if metadata.ResourceData.HasChange("auto_scaling_for_storage_enabled") {
					generalUpdate.Properties.IsAutoScalingForStorageEnabled = pointer.To(model.AutoScalingForStorageEnabled)
				}
				if metadata.ResourceData.HasChange("allowed_ips") {
					generalUpdate.Properties.WhitelistedIPs = pointer.To(model.AllowedIps)
				}

				if err := client.UpdateThenPoll(ctx, *id, generalUpdate); err != nil {
					return fmt.Errorf("updating general properties for %s: %+v", *id, err)
				}
			}

			// update password
			if needsPasswordUpdate {
				passwordUpdate := autonomousdatabases.AutonomousDatabaseUpdate{
					Properties: &autonomousdatabases.AutonomousDatabaseUpdateProperties{
						AdminPassword: pointer.To(model.AdminPassword),
					},
				}

				if err := client.UpdateThenPoll(ctx, *id, passwordUpdate); err != nil {
					return fmt.Errorf("updating Admin password for %s: %+v", *id, err)
				}
			}

			if needsBackupRetentionDaysUpdate {
				backupRetentionDaysUpdate := autonomousdatabases.AutonomousDatabaseUpdate{
					Properties: &autonomousdatabases.AutonomousDatabaseUpdateProperties{
						BackupRetentionPeriodInDays: pointer.To(model.BackupRetentionPeriodInDays),
					},
				}

				if err := client.UpdateThenPoll(ctx, *id, backupRetentionDaysUpdate); err != nil {
					return fmt.Errorf("updating `backup_retention_period_in_days` for %s: %+v", *id, err)
				}
			}

			// Step 2: Handle backup schedule update separately
			if needsBackupScheduleUpdate {
				backupUpdate := autonomousdatabases.AutonomousDatabaseUpdate{
					Properties: &autonomousdatabases.AutonomousDatabaseUpdateProperties{
						LongTermBackupSchedule: expandLongTermBackupSchedule(model.LongTermBackUpSchedule),
					},
				}

				if err := client.UpdateThenPoll(ctx, *id, backupUpdate); err != nil {
					return fmt.Errorf("updating backup schedule for %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
