// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_backup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabasebackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
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

			backup, err := getBackupFromOCI(ctx, client, adbId, backupId)
			if err != nil {
				return fmt.Errorf("retrieving backup: %+v", err)
			}

			if backup == nil {
				if err := metadata.MarkAsGone(id); err != nil {
					return err
				}
				return nil
			}

			state := AutonomousDatabaseBackupResourceModel{
				Name:                 id.AutonomousDatabaseBackupName,
				AutonomousDatabaseId: adbId.ID(),
			}

			if props := backup.Properties; props != nil {
				state.RetentionPeriodInDays = pointer.From(props.RetentionPeriodInDays)
				state.Type = pointer.FromEnum(props.BackupType)
			}

			return metadata.Encode(&state)
		},
	}
}
