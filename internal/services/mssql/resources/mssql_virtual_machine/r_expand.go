// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandSqlVirtualMachineAutoBackupSettings(input []any) (*sqlvirtualmachines.AutoBackupSettings, error) {
	ret := sqlvirtualmachines.AutoBackupSettings{
		Enable: pointer.To(false),
	}

	if len(input) > 0 {
		config := input[0].(map[string]any)
		ret.Enable = pointer.To(true)

		if v, ok := config["retention_period_in_days"]; ok {
			ret.RetentionPeriod = pointer.To(int64(v.(int)))
		}
		if v, ok := config["storage_blob_endpoint"]; ok {
			ret.StorageAccountURL = pointer.To(v.(string))
		}
		if v, ok := config["storage_account_access_key"]; ok {
			ret.StorageAccessKey = pointer.To(v.(string))
		}

		if v, ok := config["encryption_password"]; ok && v.(string) != "" {
			ret.EnableEncryption = pointer.To(true)
			ret.Password = pointer.To(v.(string))
		} else {
			ret.EnableEncryption = pointer.To(false)
		}

		if v, ok := config["system_databases_backup_enabled"]; ok {
			ret.BackupSystemDbs = pointer.To(v.(bool))
		}

		ret.BackupScheduleType = pointer.To(sqlvirtualmachines.BackupScheduleTypeAutomated)
		if v, ok := config["manual_schedule"]; ok && len(v.([]any)) > 0 {
			manualSchedule := v.([]any)[0].(map[string]any)
			ret.BackupScheduleType = pointer.To(sqlvirtualmachines.BackupScheduleTypeManual)

			fullBackupFrequency := sqlvirtualmachines.FullBackupFrequencyType(manualSchedule["full_backup_frequency"].(string))

			daysOfWeek := manualSchedule["days_of_week"].(*pluginsdk.Set).List()
			if len(daysOfWeek) > 0 {
				if !strings.EqualFold(string(fullBackupFrequency), string(sqlvirtualmachines.FullBackupFrequencyTypeWeekly)) {
					return nil, fmt.Errorf("`manual_schedule.0.days_of_week` can only be specified when `manual_schedule.0.full_backup_frequency` is set to `Weekly`")
				}
				ret.DaysOfWeek = expandSqlVirtualMachineAutoBackupSettingsDaysOfWeek(daysOfWeek)
			}

			ret.FullBackupFrequency = &fullBackupFrequency
			ret.FullBackupStartTime = pointer.To(int64(manualSchedule["full_backup_start_hour"].(int)))
			ret.FullBackupWindowHours = pointer.To(int64(manualSchedule["full_backup_window_in_hours"].(int)))
			ret.LogBackupFrequency = pointer.To(int64(manualSchedule["log_backup_frequency_in_minutes"].(int)))
		}
	}

	return &ret, nil
}

func expandSqlVirtualMachineAutoBackupSettingsDaysOfWeek(input []any) *[]sqlvirtualmachines.AutoBackupDaysOfWeek {
	result := make([]sqlvirtualmachines.AutoBackupDaysOfWeek, 0)
	for _, item := range input {
		result = append(result, sqlvirtualmachines.AutoBackupDaysOfWeek(item.(string)))
	}
	return &result
}

func expandSqlVirtualMachineAutoPatchingSettings(input []any) *sqlvirtualmachines.AutoPatchingSettings {
	if len(input) == 0 {
		return &sqlvirtualmachines.AutoPatchingSettings{
			Enable: pointer.To(false),
		}
	}

	autoPatchingSetting := input[0].(map[string]any)

	return &sqlvirtualmachines.AutoPatchingSettings{
		Enable:                        pointer.To(true),
		MaintenanceWindowDuration:     pointer.To(int64(autoPatchingSetting["maintenance_window_duration_in_minutes"].(int))),
		MaintenanceWindowStartingHour: pointer.To(int64(autoPatchingSetting["maintenance_window_starting_hour"].(int))),
		DayOfWeek:                     pointer.ToEnum[sqlvirtualmachines.DayOfWeek](autoPatchingSetting["day_of_week"].(string)),
	}
}

func expandSqlVirtualMachineAssessmentSettings(input []any) *sqlvirtualmachines.AssessmentSettings {
	if len(input) == 0 {
		return nil
	}
	assessmentSetting := input[0].(map[string]any)

	return &sqlvirtualmachines.AssessmentSettings{
		Enable:         pointer.To(true),
		RunImmediately: pointer.To(assessmentSetting["run_immediately"].(bool)),
		Schedule:       expandSqlVirtualMachineAssessmentSettingsSchedule(assessmentSetting["schedule"].([]any)),
	}
}

func expandSqlVirtualMachineAssessmentSettingsSchedule(input []any) *sqlvirtualmachines.Schedule {
	if len(input) == 0 {
		return &sqlvirtualmachines.Schedule{}
	}

	scheduleConfig := input[0].(map[string]any)

	schedule := &sqlvirtualmachines.Schedule{
		Enable:    pointer.To(true),
		DayOfWeek: pointer.ToEnum[sqlvirtualmachines.AssessmentDayOfWeek](scheduleConfig["day_of_week"].(string)),
		StartTime: pointer.To(scheduleConfig["start_time"].(string)),
	}

	if weeklyInterval := scheduleConfig["weekly_interval"].(int); weeklyInterval != 0 {
		schedule.WeeklyInterval = pointer.To(int64(weeklyInterval))
	}

	if monthlyOccurrence := scheduleConfig["monthly_occurrence"].(int); monthlyOccurrence != 0 {
		schedule.MonthlyOccurrence = pointer.To(int64(monthlyOccurrence))
	}

	return schedule
}

func expandSqlVirtualMachineKeyVaultCredential(input []any) *sqlvirtualmachines.KeyVaultCredentialSettings {
	if len(input) == 0 {
		return nil
	}
	keyVaultCredentialSetting := input[0].(map[string]any)

	return &sqlvirtualmachines.KeyVaultCredentialSettings{
		Enable:                 pointer.To(true),
		CredentialName:         pointer.To(keyVaultCredentialSetting["name"].(string)),
		AzureKeyVaultURL:       pointer.To(keyVaultCredentialSetting["key_vault_url"].(string)),
		ServicePrincipalName:   pointer.To(keyVaultCredentialSetting["service_principal_name"].(string)),
		ServicePrincipalSecret: pointer.To(keyVaultCredentialSetting["service_principal_secret"].(string)),
	}
}

func expandSqlVirtualMachineStorageConfigurationSettings(input []any) *sqlvirtualmachines.StorageConfigurationSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	storageSettings := input[0].(map[string]any)

	return &sqlvirtualmachines.StorageConfigurationSettings{
		DiskConfigurationType: pointer.ToEnum[sqlvirtualmachines.DiskConfigurationType](storageSettings["disk_type"].(string)),
		StorageWorkloadType:   pointer.ToEnum[sqlvirtualmachines.StorageWorkloadType](storageSettings["storage_workload_type"].(string)),
		SqlSystemDbOnDataDisk: pointer.To(storageSettings["system_db_on_data_disk_enabled"].(bool)),
		SqlDataSettings:       expandSqlVirtualMachineDataStorageSettings(storageSettings["data_settings"].([]any)),
		SqlLogSettings:        expandSqlVirtualMachineDataStorageSettings(storageSettings["log_settings"].([]any)),
		SqlTempDbSettings:     expandSqlVirtualMachineTempDbSettings(storageSettings["temp_db_settings"].([]any)),
	}
}

func expandSqlVirtualMachineDataStorageSettings(input []any) *sqlvirtualmachines.SQLStorageSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	dataStorageSettings := input[0].(map[string]any)

	return &sqlvirtualmachines.SQLStorageSettings{
		Luns:            expandSqlVirtualMachineStorageSettingsLuns(dataStorageSettings["luns"].([]any)),
		DefaultFilePath: pointer.To(dataStorageSettings["default_file_path"].(string)),
	}
}

func expandSqlVirtualMachineStorageSettingsLuns(input []any) *[]int64 {
	expandedLuns := make([]int64, 0)
	for i := range input {
		if input[i] != nil {
			expandedLuns = append(expandedLuns, int64(input[i].(int)))
		}
	}

	return &expandedLuns
}

func expandSqlVirtualMachineTempDbSettings(input []any) *sqlvirtualmachines.SQLTempDbSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	tempDbSettings := input[0].(map[string]any)

	return &sqlvirtualmachines.SQLTempDbSettings{
		Luns:            expandSqlVirtualMachineStorageSettingsLuns(tempDbSettings["luns"].([]any)),
		DefaultFilePath: pointer.To(tempDbSettings["default_file_path"].(string)),
		DataFileCount:   pointer.To(int64(tempDbSettings["data_file_count"].(int))),
		DataFileSize:    pointer.To(int64(tempDbSettings["data_file_size_mb"].(int))),
		DataGrowth:      pointer.To(int64(tempDbSettings["data_file_growth_in_mb"].(int))),
		LogFileSize:     pointer.To(int64(tempDbSettings["log_file_size_mb"].(int))),
		LogGrowth:       pointer.To(int64(tempDbSettings["log_file_growth_mb"].(int))),
	}
}

func expandSqlVirtualMachineSQLInstance(input []any) (*sqlvirtualmachines.SQLInstanceSettings, error) {
	if len(input) == 0 || input[0] == nil {
		return &sqlvirtualmachines.SQLInstanceSettings{}, nil
	}

	settings := input[0].(map[string]any)
	maxServerMemoryMB := settings["max_server_memory_mb"].(int)
	minServerMemoryMB := settings["min_server_memory_mb"].(int)

	if maxServerMemoryMB < minServerMemoryMB {
		return nil, fmt.Errorf("`max_server_memory_mb` must be greater than or equal to `min_server_memory_mb`")
	}

	result := sqlvirtualmachines.SQLInstanceSettings{
		Collation:                          pointer.To(settings["collation"].(string)),
		IsIfiEnabled:                       pointer.To(settings["instant_file_initialization_enabled"].(bool)),
		IsLpimEnabled:                      pointer.To(settings["lock_pages_in_memory_enabled"].(bool)),
		IsOptimizeForAdHocWorkloadsEnabled: pointer.To(settings["adhoc_workloads_optimization_enabled"].(bool)),
		MaxDop:                             pointer.To(int64(settings["max_dop"].(int))),
		MaxServerMemoryMB:                  pointer.To(int64(maxServerMemoryMB)),
		MinServerMemoryMB:                  pointer.To(int64(minServerMemoryMB)),
	}

	return &result, nil
}

func expandSqlVirtualMachineWsfcDomainCredentials(input []any) *sqlvirtualmachines.WsfcDomainCredentials {
	if len(input) == 0 {
		return nil
	}
	wsfcDomainCredentials := input[0].(map[string]any)

	return &sqlvirtualmachines.WsfcDomainCredentials{
		ClusterBootstrapAccountPassword: pointer.To(wsfcDomainCredentials["cluster_bootstrap_account_password"].(string)),
		ClusterOperatorAccountPassword:  pointer.To(wsfcDomainCredentials["cluster_operator_account_password"].(string)),
		SqlServiceAccountPassword:       pointer.To(wsfcDomainCredentials["sql_service_account_password"].(string)),
	}
}
