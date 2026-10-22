// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_backup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabasebackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if !metadata.ResourceData.HasChange("retention_period_in_days") {
				return nil
			}
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabaseBackups

			id, err := autonomousdatabasebackups.ParseAutonomousDatabaseBackupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			adbId := autonomousdatabases.NewAutonomousDatabaseID(
				id.SubscriptionId,
				id.ResourceGroupName,
				id.AutonomousDatabaseName,
			)
			backupId := autonomousdatabasebackups.NewAutonomousDatabaseBackupID(id.SubscriptionId, id.ResourceGroupName, id.AutonomousDatabaseName, id.AutonomousDatabaseBackupName)

			var model AutonomousDatabaseBackupResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("retrieving %s: %+v", backupId, err)
			}

			if _, err = getBackupFromOCI(ctx, client, adbId, backupId); err != nil {
				return fmt.Errorf("retrieving %s: %+v", backupId, err)
			}

			update := autonomousdatabasebackups.AutonomousDatabaseBackupUpdate{
				Properties: &autonomousdatabasebackups.AutonomousDatabaseBackupUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("retention_period_in_days") {
				update.Properties.RetentionPeriodInDays = &model.RetentionPeriodInDays
			}

			if err := client.UpdateThenPoll(ctx, *id, update); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}
