// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_mysql_flexible_server

type BackupPolicyMySQLFlexibleServerModel struct {
	Name                         string                                                `tfschema:"name"`
	VaultId                      string                                                `tfschema:"vault_id"`
	BackupRepeatingTimeIntervals []string                                              `tfschema:"backup_repeating_time_intervals"`
	DefaultRetentionRule         []BackupPolicyMySQLFlexibleServerDefaultRetentionRule `tfschema:"default_retention_rule"`
	RetentionRules               []BackupPolicyMySQLFlexibleServerRetentionRule        `tfschema:"retention_rule"`
	TimeZone                     string                                                `tfschema:"time_zone"`
}

func (r Resource) ModelObject() any {
	return &BackupPolicyMySQLFlexibleServerModel{}
}
