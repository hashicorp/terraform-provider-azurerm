// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_backup

type AutonomousDatabaseBackupResourceModel struct {
	AutonomousDatabaseId string `tfschema:"autonomous_database_id"`
	Name                 string `tfschema:"name"`

	// Required
	Type                  string `tfschema:"type"`
	RetentionPeriodInDays int64  `tfschema:"retention_period_in_days"`
}

func (r Resource) ModelObject() any {
	return &AutonomousDatabaseBackupResourceModel{}
}
