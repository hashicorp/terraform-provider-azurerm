// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_cosmosdb_account

type BackupPolicyCosmosdbAccountModel struct {
	Name                        string                                            `tfschema:"name"`
	DataProtectionBackupVaultId string                                            `tfschema:"data_protection_backup_vault_id"`
	BackupFrequency             string                                            `tfschema:"backup_frequency"`
	IncrementalBackupSchedules  []string                                          `tfschema:"incremental_backup_schedules"`
	DefaultRetentionDuration    string                                            `tfschema:"default_retention_duration"`
	BackupSchedule              string                                            `tfschema:"backup_schedule"`
	WeeklyRetentionRule         []BackupPolicyCosmosdbAccountWeeklyRetentionRule  `tfschema:"weekly_retention_rule"`
	MonthlyRetentionRule        []BackupPolicyCosmosdbAccountMonthlyRetentionRule `tfschema:"monthly_retention_rule"`
	YearlyRetentionRule         []BackupPolicyCosmosdbAccountYearlyRetentionRule  `tfschema:"yearly_retention_rule"`
}

type BackupPolicyCosmosdbAccountWeeklyRetentionRule struct {
	Duration   string   `tfschema:"duration"`
	DaysOfWeek []string `tfschema:"days_of_week"`
}

type BackupPolicyCosmosdbAccountMonthlyRetentionRule struct {
	Duration     string   `tfschema:"duration"`
	DaysOfWeek   []string `tfschema:"days_of_week"`
	WeeksOfMonth []string `tfschema:"weeks_of_month"`
}

type BackupPolicyCosmosdbAccountYearlyRetentionRule struct {
	Duration     string   `tfschema:"duration"`
	DaysOfWeek   []string `tfschema:"days_of_week"`
	WeeksOfMonth []string `tfschema:"weeks_of_month"`
	MonthsOfYear []string `tfschema:"months_of_year"`
}

type BackupPolicyCosmosdbAccountRetentionRule struct {
	Name             string
	Duration         string
	BackupOccurrence string
	DaysOfWeek       []string
	MonthsOfYear     []string
	WeeksOfMonth     []string
}

func (r Resource) ModelObject() any {
	return &BackupPolicyCosmosdbAccountModel{}
}
