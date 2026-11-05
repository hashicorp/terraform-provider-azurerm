// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automanage_configuration

func flattenAntiMalwareConfig(configMap map[string]any) []AntimalwareConfiguration {
	if val, ok := configMap["Antimalware/Enable"]; !ok || (val == nil) {
		return []AntimalwareConfiguration{}
	}

	antimalware := make([]AntimalwareConfiguration, 1)
	antimalware[0] = AntimalwareConfiguration{}

	if val, ok := configMap["Antimalware/EnableRealTimeProtection"]; ok {
		antimalware[0].RealTimeProtectionEnabled = val.(bool)
	}

	if val, ok := configMap["Antimalware/RunScheduledScan"]; ok {
		antimalware[0].ScheduledScanEnabled = val.(bool)
	}

	if val, ok := configMap["Antimalware/ScanType"]; ok {
		antimalware[0].ScanType = val.(string)
	}

	if val, ok := configMap["Antimalware/ScanDay"]; ok {
		antimalware[0].ScanDay = int64(val.(float64))
	}

	if val, ok := configMap["Antimalware/ScanTimeInMinutes"]; ok {
		antimalware[0].ScanTimeInMinutes = int64(val.(float64))
	}

	exclusions := AntimalwareExclusions{}
	exclusionsChanged := false
	if val, ok := configMap["Antimalware/Exclusions/Extensions"]; ok {
		exclusions.Extensions = val.(string)
		exclusionsChanged = true
	}

	if val, ok := configMap["Antimalware/Exclusions/Paths"]; ok {
		exclusions.Paths = val.(string)
		exclusionsChanged = true
	}

	if val, ok := configMap["Antimalware/Exclusions/Processes"]; ok {
		exclusions.Processes = val.(string)
		exclusionsChanged = true
	}

	if exclusionsChanged {
		antimalware[0].Exclusions = make([]AntimalwareExclusions, 1)
		antimalware[0].Exclusions[0] = exclusions
	}

	return antimalware
}

func flattenAzureSecurityBaselineConfig(configMap map[string]any) []AzureSecurityBaselineConfiguration {
	if val, ok := configMap["AzureSecurityBaseline/Enable"]; !ok || (val == nil) {
		return []AzureSecurityBaselineConfiguration{}
	}

	azureSecurityBaseline := make([]AzureSecurityBaselineConfiguration, 1)
	azureSecurityBaseline[0] = AzureSecurityBaselineConfiguration{}

	if val, ok := configMap["AzureSecurityBaseline/AssignmentType"]; ok {
		azureSecurityBaseline[0].AssignmentType = val.(string)
	}

	return azureSecurityBaseline
}

func flattenBackupConfig(configMap map[string]any) []BackupConfiguration {
	if val, ok := configMap["Backup/Enable"]; !ok || (val == nil) {
		return []BackupConfiguration{}
	}

	backup := make([]BackupConfiguration, 1)
	backup[0] = BackupConfiguration{}

	if val, ok := configMap["Backup/PolicyName"]; ok {
		backup[0].PolicyName = val.(string)
	}

	if val, ok := configMap["Backup/TimeZone"]; ok {
		backup[0].TimeZone = val.(string)
	}

	if val, ok := configMap["Backup/InstantRpRetentionRangeInDays"]; ok {
		backup[0].InstantRpRetentionRangeInDays = int64(val.(float64))
	}

	schedulePolicy := SchedulePolicyConfiguration{}
	schedulePolicyChanged := false
	if val, ok := configMap["Backup/SchedulePolicy/ScheduleRunFrequency"]; ok {
		schedulePolicy.ScheduleRunFrequency = val.(string)
		schedulePolicyChanged = true
	}

	if val, ok := configMap["Backup/SchedulePolicy/ScheduleRunTimes"]; ok {
		schedulePolicy.ScheduleRunTimes = flattenToListOfString(val)
		schedulePolicyChanged = true
	}

	if val, ok := configMap["Backup/SchedulePolicy/ScheduleRunDays"]; ok {
		schedulePolicy.ScheduleRunDays = flattenToListOfString(val)
		schedulePolicyChanged = true
	}

	if val, ok := configMap["Backup/SchedulePolicy/SchedulePolicyType"]; ok {
		schedulePolicy.SchedulePolicyType = val.(string)
		schedulePolicyChanged = true
	}

	if schedulePolicyChanged {
		backup[0].SchedulePolicy = make([]SchedulePolicyConfiguration, 1)
		backup[0].SchedulePolicy[0] = schedulePolicy
	}

	retentionPolicy := RetentionPolicyConfiguration{}
	retentionPolicyChanged := false
	if val, ok := configMap["Backup/RetentionPolicy/RetentionPolicyType"]; ok {
		retentionPolicy.RetentionPolicyType = val.(string)
		retentionPolicyChanged = true
	}

	dailySchedule := ScheduleConfiguration{}
	dailyScheduleChanged := false
	if val, ok := configMap["Backup/RetentionPolicy/DailySchedule/RetentionTimes"]; ok {
		dailySchedule.RetentionTimes = flattenToListOfString(val)
		dailyScheduleChanged = true
	}

	retentionDuration := RetentionDurationConfiguration{}
	retentionDurationChanged := false
	if val, ok := configMap["Backup/RetentionPolicy/DailySchedule/RetentionDuration/Count"]; ok {
		retentionDuration.Count = int64(val.(float64))
		retentionDurationChanged = true
	}

	if val, ok := configMap["Backup/RetentionPolicy/DailySchedule/RetentionDuration/DurationType"]; ok {
		retentionDuration.DurationType = val.(string)
		retentionDurationChanged = true
	}

	if retentionDurationChanged {
		dailySchedule.RetentionDuration = make([]RetentionDurationConfiguration, 1)
		dailySchedule.RetentionDuration[0] = retentionDuration
	}

	if dailyScheduleChanged || retentionDurationChanged {
		retentionPolicy.DailySchedule = make([]ScheduleConfiguration, 1)
		retentionPolicy.DailySchedule[0] = dailySchedule
	}

	weeklySchedule := ScheduleConfiguration{}
	weeklyScheduleChanged := false
	if val, ok := configMap["Backup/RetentionPolicy/WeeklySchedule/RetentionTimes"]; ok {
		weeklySchedule.RetentionTimes = flattenToListOfString(val)
		weeklyScheduleChanged = true
	}

	weeklyRetentionDuration := RetentionDurationConfiguration{}
	weeklyRetentionDurationChanged := false
	if val, ok := configMap["Backup/RetentionPolicy/WeeklySchedule/RetentionDuration/Count"]; ok {
		weeklyRetentionDuration.Count = int64(val.(float64))
		weeklyRetentionDurationChanged = true
	}

	if val, ok := configMap["Backup/RetentionPolicy/WeeklySchedule/RetentionDuration/DurationType"]; ok {
		weeklyRetentionDuration.DurationType = val.(string)
		weeklyRetentionDurationChanged = true
	}

	if weeklyRetentionDurationChanged {
		weeklySchedule.RetentionDuration = make([]RetentionDurationConfiguration, 1)
		weeklySchedule.RetentionDuration[0] = weeklyRetentionDuration
	}

	if weeklyScheduleChanged || weeklyRetentionDurationChanged {
		retentionPolicy.WeeklySchedule = make([]ScheduleConfiguration, 1)
		retentionPolicy.WeeklySchedule[0] = weeklySchedule
	}

	if retentionPolicyChanged || dailyScheduleChanged || retentionDurationChanged || weeklyScheduleChanged || weeklyRetentionDurationChanged {
		backup[0].RetentionPolicy = make([]RetentionPolicyConfiguration, 1)
		backup[0].RetentionPolicy[0] = retentionPolicy
	}

	return backup
}

func flattenToListOfString(val any) []string {
	lis := val.([]any)
	strs := make([]string, len(lis))
	for i, v := range lis {
		strs[i] = v.(string)
	}
	return strs
}
