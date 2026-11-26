// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceMssqlVirtualMachineSetFlatten(d *pluginsdk.ResourceData, id *sqlvirtualmachines.SqlVirtualMachineId, model *sqlvirtualmachines.SqlVirtualMachine) error {
	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("virtual_machine_id", props.VirtualMachineResourceId)
			sqlLicenseType := ""
			if licenceType := props.SqlServerLicenseType; licenceType != nil {
				sqlLicenseType = string(*licenceType)
			}
			d.Set("sql_license_type", sqlLicenseType)
			if err := d.Set("auto_backup", flattenSqlVirtualMachineAutoBackup(props.AutoBackupSettings, d)); err != nil {
				return fmt.Errorf("setting `auto_backup`: %+v", err)
			}

			if err := d.Set("auto_patching", flattenSqlVirtualMachineAutoPatching(props.AutoPatchingSettings)); err != nil {
				return fmt.Errorf("setting `auto_patching`: %+v", err)
			}

			if err := d.Set("assessment", flattenSqlVirtualMachineAssessmentSettings(props.AssessmentSettings)); err != nil {
				return fmt.Errorf("setting `assessment`: %+v", err)
			}

			if err := d.Set("key_vault_credential", flattenSqlVirtualMachineKeyVaultCredential(props.KeyVaultCredentialSettings, d)); err != nil {
				return fmt.Errorf("setting `key_vault_credential`: %+v", err)
			}

			if mgmtSettings := props.ServerConfigurationsManagementSettings; mgmtSettings != nil {
				if cfgs := mgmtSettings.AdditionalFeaturesServerConfigurations; cfgs != nil {
					d.Set("r_services_enabled", mgmtSettings.AdditionalFeaturesServerConfigurations.IsRServicesEnabled)
				}
				if scus := mgmtSettings.SqlConnectivityUpdateSettings; scus != nil {
					d.Set("sql_connectivity_port", mgmtSettings.SqlConnectivityUpdateSettings.Port)
					d.Set("sql_connectivity_type", pointer.From(mgmtSettings.SqlConnectivityUpdateSettings.ConnectivityType))
				}

				d.Set("sql_instance", flattenSqlVirtualMachineSQLInstance(mgmtSettings.SqlInstanceSettings))
			}

			// `storage_configuration.0.storage_workload_type` is in a different spot than the rest of the `storage_configuration`
			// so we'll grab that here and pass it along
			storageWorkloadType := ""
			if props.ServerConfigurationsManagementSettings != nil && props.ServerConfigurationsManagementSettings.SqlWorkloadTypeUpdateSettings != nil && props.ServerConfigurationsManagementSettings.SqlWorkloadTypeUpdateSettings.SqlWorkloadType != nil {
				storageWorkloadType = string(*props.ServerConfigurationsManagementSettings.SqlWorkloadTypeUpdateSettings.SqlWorkloadType)
			}

			if err := d.Set("storage_configuration", flattenSqlVirtualMachineStorageConfigurationSettings(props.StorageConfigurationSettings, storageWorkloadType)); err != nil {
				return fmt.Errorf("setting `storage_configuration`: %+v", err)
			}

			sqlVirtualMachineGroupId := ""
			if props.SqlVirtualMachineGroupResourceId != nil {
				parsedId, err := sqlvirtualmachinegroups.ParseSqlVirtualMachineGroupIDInsensitively(*props.SqlVirtualMachineGroupResourceId)
				if err != nil {
					return err
				}

				// get correct casing for subscription in id due to https://github.com/Azure/azure-rest-api-specs/issues/25211
				sqlVirtualMachineGroupId = sqlvirtualmachinegroups.NewSqlVirtualMachineGroupID(id.SubscriptionId, parsedId.ResourceGroupName, parsedId.SqlVirtualMachineGroupName).ID()
			}
			d.Set("sql_virtual_machine_group_id", sqlVirtualMachineGroupId)

			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return err
			}
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenSqlVirtualMachineAutoBackup(autoBackup *sqlvirtualmachines.AutoBackupSettings, d *pluginsdk.ResourceData) []any {
	if autoBackup == nil || autoBackup.Enable == nil || !*autoBackup.Enable {
		return []any{}
	}

	manualSchedule := make([]any, 0)
	if autoBackup.BackupScheduleType != nil && strings.EqualFold(string(*autoBackup.BackupScheduleType), string(sqlvirtualmachines.BackupScheduleTypeManual)) {
		var fullBackupStartHour int
		if autoBackup.FullBackupStartTime != nil {
			fullBackupStartHour = int(*autoBackup.FullBackupStartTime)
		}

		var fullBackupWindowHours int
		if autoBackup.FullBackupWindowHours != nil {
			fullBackupWindowHours = int(*autoBackup.FullBackupWindowHours)
		}

		var logBackupFrequency int
		if autoBackup.LogBackupFrequency != nil {
			logBackupFrequency = int(*autoBackup.LogBackupFrequency)
			// API returns 60 minutes as zero
			if logBackupFrequency == 0 {
				logBackupFrequency = 60
			}
		}

		var fullBackupFrequency string
		if autoBackup.FullBackupFrequency != nil {
			fullBackupFrequency = string(*autoBackup.FullBackupFrequency)
		}

		manualSchedule = []any{
			map[string]any{
				"full_backup_frequency":           fullBackupFrequency,
				"full_backup_start_hour":          fullBackupStartHour,
				"full_backup_window_in_hours":     fullBackupWindowHours,
				"log_backup_frequency_in_minutes": logBackupFrequency,
				"days_of_week":                    flattenSqlVirtualMachineAutoBackupDaysOfWeek(autoBackup.DaysOfWeek),
			},
		}
	}

	var retentionPeriod int
	if autoBackup.RetentionPeriod != nil {
		retentionPeriod = int(*autoBackup.RetentionPeriod)
	}

	// Password, StorageAccessKey, StorageAccountURL are not returned, so we try to copy them
	// from existing config as a best effort.
	storageKey := d.Get("auto_backup.0.storage_account_access_key").(string)
	blobEndpoint := d.Get("auto_backup.0.storage_blob_endpoint").(string)
	encryptionPassword := ""

	// Copy password from config only if encryption is enabled in Azure
	if pointer.From(autoBackup.EnableEncryption) {
		encryptionPassword = d.Get("auto_backup.0.encryption_password").(string)
	}

	return []any{
		map[string]any{
			"encryption_password":             encryptionPassword,
			"manual_schedule":                 manualSchedule,
			"retention_period_in_days":        retentionPeriod,
			"storage_account_access_key":      storageKey,
			"storage_blob_endpoint":           blobEndpoint,
			"system_databases_backup_enabled": autoBackup.BackupSystemDbs != nil && *autoBackup.BackupSystemDbs,
		},
	}
}

func flattenSqlVirtualMachineAutoBackupDaysOfWeek(daysOfWeek *[]sqlvirtualmachines.AutoBackupDaysOfWeek) []any {
	output := make([]any, 0)

	if daysOfWeek != nil {
		for _, v := range *daysOfWeek {
			output = append(output, string(v))
		}
	}

	return output
}

func flattenSqlVirtualMachineAutoPatching(autoPatching *sqlvirtualmachines.AutoPatchingSettings) []any {
	if autoPatching == nil || autoPatching.Enable == nil || !*autoPatching.Enable {
		return []any{}
	}

	var startHour int
	if autoPatching.MaintenanceWindowStartingHour != nil {
		startHour = int(*autoPatching.MaintenanceWindowStartingHour)
	}

	var duration int
	if autoPatching.MaintenanceWindowDuration != nil {
		duration = int(*autoPatching.MaintenanceWindowDuration)
	}

	var dayOfWeek string
	if autoPatching.DayOfWeek != nil {
		dayOfWeek = string(*autoPatching.DayOfWeek)
	}

	return []any{
		map[string]any{
			"day_of_week":                            dayOfWeek,
			"maintenance_window_starting_hour":       startHour,
			"maintenance_window_duration_in_minutes": duration,
		},
	}
}

func flattenSqlVirtualMachineAssessmentSettings(assessmentSettings *sqlvirtualmachines.AssessmentSettings) []any {
	if assessmentSettings == nil || assessmentSettings.Enable == nil || !*assessmentSettings.Enable {
		return []any{}
	}

	var (
		runImmediately bool
		enabled        bool
	)
	if assessmentSettings.RunImmediately != nil {
		runImmediately = *assessmentSettings.RunImmediately
	}

	if assessmentSettings.Enable != nil {
		enabled = *assessmentSettings.Enable
	}

	var attr map[string]any
	if schedule := assessmentSettings.Schedule; schedule != nil {
		var (
			weeklyInterval    int64
			monthlyOccurrence int64
			dayOfWeek         string
			startTime         string
		)

		if schedule.WeeklyInterval != nil {
			weeklyInterval = *schedule.WeeklyInterval
		}
		if schedule.MonthlyOccurrence != nil {
			monthlyOccurrence = *schedule.MonthlyOccurrence
		}
		if schedule.DayOfWeek != nil {
			dayOfWeek = string(*schedule.DayOfWeek)
		}
		if schedule.StartTime != nil {
			startTime = *schedule.StartTime
		}

		attr = map[string]any{
			"weekly_interval":    weeklyInterval,
			"monthly_occurrence": monthlyOccurrence,
			"day_of_week":        dayOfWeek,
			"start_time":         startTime,
		}
	}

	return []any{
		map[string]any{
			"run_immediately": runImmediately,
			"enabled":         enabled,
			"schedule":        []any{attr},
		},
	}
}

func flattenSqlVirtualMachineKeyVaultCredential(keyVault *sqlvirtualmachines.KeyVaultCredentialSettings, d *pluginsdk.ResourceData) []any {
	if keyVault == nil || keyVault.Enable == nil || !*keyVault.Enable {
		return []any{}
	}

	keyVaultUrl := ""
	if v, ok := d.GetOk("key_vault_credential.0.key_vault_url"); ok {
		keyVaultUrl = v.(string)
	}

	servicePrincipalName := ""
	if v, ok := d.GetOk("key_vault_credential.0.service_principal_name"); ok {
		servicePrincipalName = v.(string)
	}

	servicePrincipalSecret := ""
	if v, ok := d.GetOk("key_vault_credential.0.service_principal_secret"); ok {
		servicePrincipalSecret = v.(string)
	}

	return []any{
		map[string]any{
			"name":                     pointer.From(keyVault.CredentialName),
			"key_vault_url":            keyVaultUrl,
			"service_principal_name":   servicePrincipalName,
			"service_principal_secret": servicePrincipalSecret,
		},
	}
}

func flattenSqlVirtualMachineStorageConfigurationSettings(input *sqlvirtualmachines.StorageConfigurationSettings, storageWorkloadType string) []any {
	if input == nil {
		return []any{}
	}

	var diskType string
	if input.DiskConfigurationType != nil {
		diskType = string(*input.DiskConfigurationType)
	}

	output := map[string]any{
		"storage_workload_type":          storageWorkloadType,
		"disk_type":                      diskType,
		"system_db_on_data_disk_enabled": pointer.From(input.SqlSystemDbOnDataDisk),
		"data_settings":                  flattenSqlVirtualMachineStorageSettings(input.SqlDataSettings),
		"log_settings":                   flattenSqlVirtualMachineStorageSettings(input.SqlLogSettings),
		"temp_db_settings":               flattenSqlVirtualMachineTempDbSettings(input.SqlTempDbSettings),
	}

	if output["storage_workload_type"].(string) == "" && output["disk_type"] == "" &&
		len(output["data_settings"].([]any)) == 0 &&
		len(output["log_settings"].([]any)) == 0 &&
		len(output["temp_db_settings"].([]any)) == 0 {
		return []any{}
	}

	return []any{output}
}

func flattenSqlVirtualMachineStorageSettings(input *sqlvirtualmachines.SQLStorageSettings) []any {
	if input == nil || input.Luns == nil {
		return []any{}
	}
	attrs := make(map[string]any)

	if input.Luns != nil {
		attrs["luns"] = *input.Luns
	}

	if input.DefaultFilePath != nil {
		attrs["default_file_path"] = *input.DefaultFilePath
	}

	return []any{attrs}
}

func flattenSqlVirtualMachineTempDbSettings(input *sqlvirtualmachines.SQLTempDbSettings) []any {
	if input == nil || input.Luns == nil {
		return []any{}
	}
	attrs := make(map[string]any)

	if input.Luns != nil {
		attrs["luns"] = *input.Luns
	}

	if input.DataFileCount != nil {
		attrs["data_file_count"] = *input.DataFileCount
	}

	if input.DataFileSize != nil {
		attrs["data_file_size_mb"] = *input.DataFileSize
	}

	if input.DataGrowth != nil {
		attrs["data_file_growth_in_mb"] = *input.DataGrowth
	}

	if input.DefaultFilePath != nil {
		attrs["default_file_path"] = *input.DefaultFilePath
	}

	if input.LogFileSize != nil {
		attrs["log_file_size_mb"] = *input.LogFileSize
	}

	if input.LogGrowth != nil {
		attrs["log_file_growth_mb"] = *input.LogGrowth
	}

	return []any{attrs}
}

func flattenSqlVirtualMachineSQLInstance(input *sqlvirtualmachines.SQLInstanceSettings) []any {
	if input == nil || input.Collation == nil {
		return []any{}
	}

	collation := *input.Collation

	var maxServerMemoryMB int64 = math.MaxInt32
	if input.MaxServerMemoryMB != nil {
		maxServerMemoryMB = *input.MaxServerMemoryMB
	}

	return []any{
		map[string]any{
			"adhoc_workloads_optimization_enabled": pointer.From(input.IsOptimizeForAdHocWorkloadsEnabled),
			"collation":                            collation,
			"instant_file_initialization_enabled":  pointer.From(input.IsIfiEnabled),
			"lock_pages_in_memory_enabled":         pointer.From(input.IsLpimEnabled),
			"max_dop":                              pointer.From(input.MaxDop),
			"max_server_memory_mb":                 maxServerMemoryMB,
			"min_server_memory_mb":                 pointer.From(input.MinServerMemoryMB),
		},
	}
}

func resourceMsSqlVirtualMachineAutoBackupSettingsRefreshFunc(ctx context.Context, client *sqlvirtualmachines.SqlVirtualMachinesClient, d *pluginsdk.ResourceData) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(d.Id())
		if err != nil {
			return nil, "Error", err
		}

		resp, err := client.Get(ctx, *id, sqlvirtualmachines.GetOperationOptions{Expand: pointer.To("*")})
		if err != nil {
			return nil, "Retry", err
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				autoBackupSettings := flattenSqlVirtualMachineAutoBackup(props.AutoBackupSettings, d)

				if len(autoBackupSettings) == 0 {
					// auto backup was nil or disabled in the response
					if v, ok := d.GetOk("auto_backup"); !ok || len(v.([]any)) == 0 {
						// also disabled in the config
						return resp, "Updated", nil
					}
					return resp, "Pending", nil
				}

				if v, ok := d.GetOk("auto_backup"); !ok || len(v.([]any)) == 0 {
					// still waiting for it to be disabled
					return resp, "Pending", nil
				}

				// check each property in the auto_backup block for drift
				for prop, val := range autoBackupSettings[0].(map[string]any) {
					v := d.Get(fmt.Sprintf("auto_backup.0.%s", prop))
					switch prop {
					case "manual_schedule":
						if m := val.([]any); len(m) > 0 {
							if b, ok := d.GetOk("auto_backup.0.manual_schedule"); !ok || len(b.([]any)) == 0 {
								// manual schedule disabled in config but still showing in response
								return resp, "Pending", nil
							}
							// check each property in the manual_schedule block for drift
							for prop2, val2 := range m[0].(map[string]any) {
								v2 := d.Get(fmt.Sprintf("auto_backup.0.manual_schedule.0.%s", prop2))
								switch prop2 {
								case "full_backup_frequency":
									if !strings.EqualFold(v2.(string), val2.(string)) {
										return resp, "Pending", nil
									}
								case "days_of_week":
									daysOfWeekLocal := make(map[string]bool, 0)
									if v2 != nil {
										for _, item := range v2.(*pluginsdk.Set).List() {
											daysOfWeekLocal[item.(string)] = true
										}
									}

									daysOfWeekRemote := make([]any, 0)
									if val2 != nil {
										daysOfWeekRemote = val2.([]any)
									}

									if len(daysOfWeekRemote) != len(daysOfWeekLocal) {
										return resp, "Pending", nil
									}

									for _, item := range daysOfWeekRemote {
										if _, ok := daysOfWeekLocal[item.(string)]; !ok {
											return resp, "Pending", nil
										}
									}
								default:
									if v2 != val2 {
										return resp, "Pending", nil
									}
								}
							}
						} else if b, ok := d.GetOk("auto_backup.0.manual_schedule"); ok || len(b.([]any)) > 0 {
							// manual schedule set in config but not reflecting in response
							return resp, "Pending", nil
						}
					default:
						if v != val {
							return resp, "Pending", nil
						}
					}
				}

				return resp, "Updated", nil
			}
		}

		return resp, "Retry", nil
	}
}

func resourceMsSqlVirtualMachineAutoPatchingSettingsRefreshFunc(ctx context.Context, client *sqlvirtualmachines.SqlVirtualMachinesClient, d *pluginsdk.ResourceData) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(d.Id())
		if err != nil {
			return nil, "Error", err
		}

		resp, err := client.Get(ctx, *id, sqlvirtualmachines.GetOperationOptions{Expand: pointer.To("*")})
		if err != nil {
			return nil, "Retry", err
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				autoPatchingSettings := flattenSqlVirtualMachineAutoPatching(props.AutoPatchingSettings)

				if len(autoPatchingSettings) == 0 {
					if v, ok := d.GetOk("auto_patching"); !ok || len(v.([]any)) == 0 {
						return resp, "Updated", nil
					}
					return resp, "Pending", nil
				}

				if v, ok := d.GetOk("auto_patching"); !ok || len(v.([]any)) == 0 {
					return resp, "Pending", nil
				}

				for prop, val := range autoPatchingSettings[0].(map[string]any) {
					v := d.Get(fmt.Sprintf("auto_patching.0.%s", prop))
					if v != val {
						return resp, "Pending", nil
					}
				}

				return resp, "Updated", nil
			}
		}

		return resp, "Retry", nil
	}
}
