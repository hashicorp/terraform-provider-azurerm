// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_postgresql_flexible_server

type BackupPolicyPostgreSQLFlexibleServerModel struct {
	Name                         string                                                     `tfschema:"name"`
	VaultId                      string                                                     `tfschema:"vault_id"`
	BackupRepeatingTimeIntervals []string                                                   `tfschema:"backup_repeating_time_intervals"`
	DefaultRetentionRule         []BackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule `tfschema:"default_retention_rule"`
	RetentionRules               []BackupPolicyPostgreSQLFlexibleServerRetentionRule        `tfschema:"retention_rule"`
	TimeZone                     string                                                     `tfschema:"time_zone"`
}

func (r Resource) ModelObject() any {
	return &BackupPolicyPostgreSQLFlexibleServerModel{}
}
