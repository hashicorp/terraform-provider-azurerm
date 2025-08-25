// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabasebackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
)

func FlattenLongTermBackUpScheduleDetails(longTermBackUpScheduleDetails *autonomousdatabases.LongTermBackUpScheduleDetails) []LongTermBackUpScheduleDetails {
	output := make([]LongTermBackUpScheduleDetails, 0)
	if longTermBackUpScheduleDetails != nil {
		return append(output, LongTermBackUpScheduleDetails{
			RepeatCadence:         string(pointer.From(longTermBackUpScheduleDetails.RepeatCadence)),
			TimeOfBackup:          pointer.From(longTermBackUpScheduleDetails.TimeOfBackup),
			RetentionPeriodInDays: pointer.From(longTermBackUpScheduleDetails.RetentionPeriodInDays),
			Enabled:               !pointer.From(longTermBackUpScheduleDetails.IsDisabled),
		})
	}
	return output
}

func findBackupByName(ctx context.Context, client *autonomousdatabasebackups.AutonomousDatabaseBackupsClient, adbId autonomousdatabases.AutonomousDatabaseId, backupId autonomousdatabasebackups.AutonomousDatabaseBackupId) (*autonomousdatabasebackups.AutonomousDatabaseBackup, error) {
	resp, err := client.ListByParent(ctx, autonomousdatabasebackups.AutonomousDatabaseId(adbId))
	if err != nil {
		return nil, fmt.Errorf("listing backups for %s: %+v", adbId.ID(), err)
	}

	id := backupId.ID()

	if model := resp.Model; model != nil {
		for _, backup := range *model {
			if backup.Id != nil && strings.EqualFold(*backup.Id, id) {
				return &backup, nil
			}
		}
	}

	return nil, nil
}
