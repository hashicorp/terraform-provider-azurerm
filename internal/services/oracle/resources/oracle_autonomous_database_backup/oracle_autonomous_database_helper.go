// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_backup

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabasebackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
)

// getBackupFromOCI retrieves a backups by making a direct API call to OCI.
// It bypasses the standard client.Get() method because backup data is not
// stored within Azure's metadata resource provider (MetaRp).
func getBackupFromOCI(ctx context.Context, client *autonomousdatabasebackups.AutonomousDatabaseBackupsClient, adbId autonomousdatabases.AutonomousDatabaseId, backupId autonomousdatabasebackups.AutonomousDatabaseBackupId) (*autonomousdatabasebackups.AutonomousDatabaseBackup, error) {
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
