// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_policy_vm_workload

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupProtectionPolicyVMWorkloadModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.RecoveryServices.ProtectionPoliciesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := protectionpolicies.NewBackupPolicyID(subscriptionId, model.ResourceGroupName, model.RecoveryVaultName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_backup_policy_vm_workload", id.ID())
				}
			}

			protectionPolicy, err := expandBackupProtectionPolicyVMWorkloadProtectionPolicies(model.ProtectionPolicies, model.WorkloadType)
			if err != nil {
				return err
			}

			properties := &protectionpolicies.ProtectionPolicyResource{
				Properties: &protectionpolicies.AzureVMWorkloadProtectionPolicy{
					Settings:            expandBackupProtectionPolicyVMWorkloadSettings(model.Settings),
					SubProtectionPolicy: protectionPolicy,
					WorkLoadType:        pointer.ToEnum[protectionpolicies.WorkloadType](model.WorkloadType),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties, protectionpolicies.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandBackupProtectionPolicyVMWorkloadSettings(input []Settings) *protectionpolicies.Settings {
	if len(input) == 0 {
		return &protectionpolicies.Settings{}
	}

	settings := input[0]
	result := &protectionpolicies.Settings{
		IsCompression: pointer.To(settings.CompressionEnabled),
	}

	if settings.TimeZone != "" {
		result.TimeZone = pointer.To(settings.TimeZone)
	}

	return result
}

func expandBackupProtectionPolicyVMWorkloadProtectionPolicies(input []ProtectionPolicy, workloadType string) (*[]protectionpolicies.SubProtectionPolicy, error) {
	if len(input) == 0 {
		return nil, nil
	}

	results := make([]protectionpolicies.SubProtectionPolicy, 0)

	for _, item := range input {
		if workloadType == string(protectionpolicies.WorkloadTypeSQLDataBase) && item.PolicyType == string(protectionpolicies.PolicyTypeIncremental) {
			return nil, fmt.Errorf("the Incremental backup isn't supported when `workload_type` is `SQLDataBase`")
		}

		backupBlock := item.Backup[0]

		// getting this ready now because its shared between *everything*, time is... complicated for this resource
		timeOfDay := backupBlock.Time
		times := make([]string, 0)
		if timeOfDay != "" {
			dateOfDay, err := time.Parse(time.RFC3339, fmt.Sprintf("2018-07-30T%s:00Z", timeOfDay))
			if err != nil {
				return nil, fmt.Errorf("generating time from %q for policy): %+v", timeOfDay, err)
			}
			times = append(times, date.Time{Time: dateOfDay}.String())
		}

		switch backupBlock.Frequency {
		case string(protectionpolicies.ScheduleRunTypeDaily):
			if len(item.RetentionDaily) == 0 {
				return nil, fmt.Errorf("`retention_daily` must be set when `backup.0.frequency` is `Daily`")
			}

			if weekdays := backupBlock.Weekdays; len(weekdays) > 0 {
				return nil, fmt.Errorf("`backup.0.weekdays` should be not set when `backup.0.frequency` is `Daily`")
			}
		case string(protectionpolicies.ScheduleRunTypeWeekly):
			if len(item.RetentionDaily) > 0 {
				return nil, fmt.Errorf("`retention_daily` must be not set when `backup.0.frequency` is `Weekly`")
			}

			if item.PolicyType != string(protectionpolicies.PolicyTypeLog) && len(backupBlock.Weekdays) == 0 {
				return nil, fmt.Errorf("`backup.weekdays` must be set when `policy_type` is not `Log` and `backup.frequency` is `Weekly`")
			}

			if item.PolicyType == string(protectionpolicies.PolicyTypeFull) && len(item.RetentionWeekly) == 0 {
				return nil, fmt.Errorf("`retention_weekly` must be set when `policy_type` is `Full` and `backup.frequency` is `Weekly`")
			}
		}

		result := protectionpolicies.SubProtectionPolicy{
			PolicyType:     pointer.ToEnum[protectionpolicies.PolicyType](item.PolicyType),
			SchedulePolicy: expandBackupProtectionPolicyVMWorkloadSchedulePolicy(item, times),
		}

		if v, err := expandBackupProtectionPolicyVMWorkloadRetentionPolicy(item, times); err != nil {
			return nil, err
		} else {
			result.RetentionPolicy = v
		}

		results = append(results, result)
	}

	return &results, nil
}

func expandBackupProtectionPolicyVMWorkloadSchedulePolicy(input ProtectionPolicy, times []string) protectionpolicies.SchedulePolicy {
	if input.PolicyType == string(protectionpolicies.PolicyTypeLog) {
		schedule := protectionpolicies.LogSchedulePolicy{}

		if v := input.Backup[0].FrequencyInMinutes; v != 0 {
			schedule.ScheduleFrequencyInMins = pointer.To(v)
		}

		return schedule
	} else {
		schedule := protectionpolicies.SimpleSchedulePolicy{}

		backupBlock := input.Backup[0]
		if backupBlock.Frequency != "" {
			schedule.ScheduleRunFrequency = pointer.ToEnum[protectionpolicies.ScheduleRunType](backupBlock.Frequency)
		}

		if len(times) > 0 {
			schedule.ScheduleRunTimes = &times
		}

		if v := backupBlock.Weekdays; len(v) > 0 {
			days := make([]protectionpolicies.DayOfWeek, 0)
			for _, day := range v {
				days = append(days, protectionpolicies.DayOfWeek(day))
			}
			schedule.ScheduleRunDays = &days
		}

		return schedule
	}
}

func expandBackupProtectionPolicyVMWorkloadRetentionPolicy(input ProtectionPolicy, times []string) (protectionpolicies.RetentionPolicy, error) {
	if input.PolicyType == string(protectionpolicies.PolicyTypeFull) {
		retentionPolicy := protectionpolicies.LongTermRetentionPolicy{}

		if len(input.RetentionDaily) > 0 {
			retentionDaily := input.RetentionDaily[0]

			retentionPolicy.DailySchedule = &protectionpolicies.DailyRetentionSchedule{
				RetentionTimes: &times,
				RetentionDuration: &protectionpolicies.RetentionDuration{
					Count:        pointer.To(retentionDaily.Count),
					DurationType: pointer.To(protectionpolicies.RetentionDurationTypeDays),
				},
			}
		}

		if len(input.RetentionWeekly) > 0 {
			retentionWeekly := input.RetentionWeekly[0]

			retentionPolicy.WeeklySchedule = &protectionpolicies.WeeklyRetentionSchedule{
				RetentionTimes: &times,
				RetentionDuration: &protectionpolicies.RetentionDuration{
					Count:        pointer.To(retentionWeekly.Count),
					DurationType: pointer.To(protectionpolicies.RetentionDurationTypeWeeks),
				},
			}

			if v := retentionWeekly.Weekdays; len(v) > 0 {
				days := make([]protectionpolicies.DayOfWeek, 0)
				for _, day := range v {
					days = append(days, protectionpolicies.DayOfWeek(day))
				}
				retentionPolicy.WeeklySchedule.DaysOfTheWeek = &days
			}
		}

		if len(input.RetentionMonthly) > 0 {
			retentionMonthly := input.RetentionMonthly[0]

			if input.Backup[0].Frequency == string(protectionpolicies.ScheduleRunTypeWeekly) && retentionMonthly.FormatType != string(protectionpolicies.RetentionScheduleFormatWeekly) {
				return nil, fmt.Errorf("`retention_monthly.format_type` must be `Weekly` when `policy_type` is `Full` and `frequency` is `Weekly`")
			}

			if retentionMonthly.FormatType == string(protectionpolicies.RetentionScheduleFormatDaily) && (len(retentionMonthly.Monthdays) == 0) {
				return nil, fmt.Errorf("`retention_monthly.monthdays` must be set when `retention_monthly.format_type` is `Daily`")
			}

			if retentionMonthly.FormatType == string(protectionpolicies.RetentionScheduleFormatWeekly) && (len(retentionMonthly.Weeks) == 0 || len(retentionMonthly.Weekdays) == 0) {
				return nil, fmt.Errorf("`retention_monthly.weeks` and `retention_monthly.weekdays` must be set when `retention_monthly.format_type` is `Weekly`")
			}

			retentionPolicy.MonthlySchedule = &protectionpolicies.MonthlyRetentionSchedule{
				RetentionScheduleFormatType: pointer.ToEnum[protectionpolicies.RetentionScheduleFormat](retentionMonthly.FormatType),
				RetentionScheduleDaily:      expandBackupProtectionPolicyVMWorkloadRetentionDailyFormat(retentionMonthly.Monthdays),
				RetentionScheduleWeekly:     expandBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(retentionMonthly.Weekdays, retentionMonthly.Weeks),
				RetentionTimes:              &times,
				RetentionDuration: &protectionpolicies.RetentionDuration{
					Count:        pointer.To(retentionMonthly.Count),
					DurationType: pointer.To(protectionpolicies.RetentionDurationTypeMonths),
				},
			}
		}

		if len(input.RetentionYearly) > 0 {
			retentionYearly := input.RetentionYearly[0]

			if input.Backup[0].Frequency == string(protectionpolicies.ScheduleRunTypeWeekly) && retentionYearly.FormatType != string(protectionpolicies.RetentionScheduleFormatWeekly) {
				return nil, fmt.Errorf("`retention_yearly.format_type` must be `Weekly` when `policy_type` is `Full` and `frequency` is `Weekly`")
			}

			if retentionYearly.FormatType == string(protectionpolicies.RetentionScheduleFormatDaily) && (len(retentionYearly.Monthdays) == 0) {
				return nil, fmt.Errorf("`retention_yearly.monthdays` must be set when `retention_yearly.format_type` is `Daily`")
			}

			if retentionYearly.FormatType == string(protectionpolicies.RetentionScheduleFormatWeekly) && (len(retentionYearly.Weeks) == 0 || len(retentionYearly.Weekdays) == 0) {
				return nil, fmt.Errorf("`retention_yearly.weeks` and `retention_yearly.weekdays` must be set when `retention_yearly.format_type` is `Weekly`")
			}

			retentionPolicy.YearlySchedule = &protectionpolicies.YearlyRetentionSchedule{
				RetentionScheduleFormatType: pointer.ToEnum[protectionpolicies.RetentionScheduleFormat](retentionYearly.FormatType),
				RetentionScheduleDaily:      expandBackupProtectionPolicyVMWorkloadRetentionDailyFormat(retentionYearly.Monthdays),
				RetentionScheduleWeekly:     expandBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(retentionYearly.Weekdays, retentionYearly.Weeks),
				RetentionTimes:              &times,
				RetentionDuration: &protectionpolicies.RetentionDuration{
					Count:        pointer.To(retentionYearly.Count),
					DurationType: pointer.To(protectionpolicies.RetentionDurationTypeYears),
				},
			}

			if v := retentionYearly.Months; v != nil {
				months := make([]protectionpolicies.MonthOfYear, 0)
				for _, month := range v {
					months = append(months, protectionpolicies.MonthOfYear(month))
				}
				retentionPolicy.YearlySchedule.MonthsOfYear = &months
			}
		}

		return retentionPolicy, nil
	} else {
		retentionPolicy := protectionpolicies.SimpleRetentionPolicy{}

		if len(input.SimpleRetention) > 0 {
			simpleRetention := input.SimpleRetention[0]

			retentionPolicy.RetentionDuration = &protectionpolicies.RetentionDuration{
				Count:        pointer.To(simpleRetention.Count),
				DurationType: pointer.To(protectionpolicies.RetentionDurationTypeDays),
			}
		}

		return retentionPolicy, nil
	}
}

func expandBackupProtectionPolicyVMWorkloadRetentionDailyFormat(input []int64) *protectionpolicies.DailyRetentionFormat {
	if len(input) == 0 {
		return nil
	}

	daily := protectionpolicies.DailyRetentionFormat{}

	days := make([]protectionpolicies.Day, 0)
	for _, item := range input {
		day := protectionpolicies.Day{
			Date: pointer.To(item),
		}

		if item == 0 {
			day.IsLast = pointer.To(true)
		} else {
			day.IsLast = pointer.To(false)
		}

		days = append(days, day)
	}
	daily.DaysOfTheMonth = &days

	return &daily
}

func expandBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(weekdays, weeks []string) *protectionpolicies.WeeklyRetentionFormat {
	if len(weekdays) == 0 && len(weeks) == 0 {
		return nil
	}

	weekly := protectionpolicies.WeeklyRetentionFormat{}

	if len(weekdays) > 0 {
		weekdaysBlock := make([]protectionpolicies.DayOfWeek, 0)
		for _, day := range weekdays {
			weekdaysBlock = append(weekdaysBlock, protectionpolicies.DayOfWeek(day))
		}
		weekly.DaysOfTheWeek = &weekdaysBlock
	}

	if len(weeks) > 0 {
		weeksBlock := make([]protectionpolicies.WeekOfMonth, 0)
		for _, week := range weeks {
			weeksBlock = append(weeksBlock, protectionpolicies.WeekOfMonth(week))
		}
		weekly.WeeksOfTheMonth = &weeksBlock
	}

	return &weekly
}
