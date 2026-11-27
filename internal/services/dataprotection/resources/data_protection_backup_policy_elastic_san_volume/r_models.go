// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_elastic_san_volume

type BackupPolicyElasticSanVolumeModel struct {
	Name                        string `tfschema:"name"`
	DataProtectionBackupVaultId string `tfschema:"data_protection_backup_vault_id"`
	DailyBackupSchedule         string `tfschema:"daily_backup_schedule"`
	DefaultRetentionDuration    string `tfschema:"default_retention_duration"`
	DailyRetentionDuration      string `tfschema:"daily_retention_duration"`
	MonthlyRetentionDuration    string `tfschema:"monthly_retention_duration"`
	WeeklyBackupSchedule        string `tfschema:"weekly_backup_schedule"`
	WeeklyRetentionDuration     string `tfschema:"weekly_retention_duration"`
}

func (r Resource) ModelObject() any {
	return &BackupPolicyElasticSanVolumeModel{}
}
