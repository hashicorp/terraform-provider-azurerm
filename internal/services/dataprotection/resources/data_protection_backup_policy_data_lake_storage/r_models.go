// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_data_lake_storage

type BackupPolicyDataLakeStorageModel struct {
	Name                        string                                     `tfschema:"name"`
	DataProtectionBackupVaultId string                                     `tfschema:"data_protection_backup_vault_id"`
	BackupSchedule              []string                                   `tfschema:"backup_schedule"`
	DefaultRetentionDuration    string                                     `tfschema:"default_retention_duration"`
	RetentionRules              []BackupPolicyDataLakeStorageRetentionRule `tfschema:"retention_rule"`
	TimeZone                    string                                     `tfschema:"time_zone"`
}

func (r Resource) ModelObject() any {
	return &BackupPolicyDataLakeStorageModel{}
}
