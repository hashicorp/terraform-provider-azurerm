// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectionpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceBackupProtectionPolicyVM() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBackupProtectionPolicyVMCreateUpdate,
		Read:   resourceBackupProtectionPolicyVMRead,
		Update: resourceBackupProtectionPolicyVMCreateUpdate,
		Delete: resourceBackupProtectionPolicyVMDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := protectionpolicies.ParseBackupPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceBackupProtectionPolicyVMSchema(),

		// if daily, we need daily retention
		// if weekly daily cannot be set, and we need weekly
		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
			_, hasDaily := diff.GetOk("retention_daily")
			_, hasWeekly := diff.GetOk("retention_weekly")

			frequency, _ := diff.GetOk("backup.0.frequency")
			switch frequency.(string) {
			case string(protectionpolicies.ScheduleRunTypeHourly):
				if !hasDaily {
					return fmt.Errorf("`retention_daily` must be set when backup.0.frequency is hourly")
				}

				if _, ok := diff.GetOk("backup.0.weekdays"); ok {
					return fmt.Errorf("`backup.0.weekdays` should be not set when backup.0.frequency is hourly")
				}
			case string(protectionpolicies.ScheduleRunTypeDaily):
				if !hasDaily {
					return fmt.Errorf("`retention_daily` must be set when backup.0.frequency is daily")
				}

				if _, ok := diff.GetOk("backup.0.weekdays"); ok {
					return fmt.Errorf("`backup.0.weekdays` should be not set when backup.0.frequency is daily")
				}

				if _, ok := diff.GetOk("backup.0.hour_interval"); ok {
					return fmt.Errorf("`backup.0.hour_interval` should be not set when backup.0.frequency is daily")
				}

				if _, ok := diff.GetOk("backup.0.hour_duration"); ok {
					return fmt.Errorf("`backup.0.hour_duration` should be not set when backup.0.frequency is daily")
				}
			case string(protectionpolicies.ScheduleRunTypeWeekly):
				if hasDaily {
					return fmt.Errorf("`retention_daily` must be not set when backup.0.frequency is weekly")
				}
				if !hasWeekly {
					return fmt.Errorf("`retention_weekly` must be set when backup.0.frequency is weekly")
				}

				if _, ok := diff.GetOk("backup.0.hour_interval"); ok {
					return fmt.Errorf("`backup.0.hour_interval` should be not set when backup.0.frequency is weekly")
				}

				if _, ok := diff.GetOk("backup.0.hour_duration"); ok {
					return fmt.Errorf("`backup.0.hour_duration` should be not set when backup.0.frequency is weekly")
				}
			default:
				return fmt.Errorf("Unrecognized value for backup.0.frequency")
			}
			return nil
		}),
	}
}

func resourceBackupProtectionPolicyVMCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := protectionpolicies.NewBackupPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), d.Get("name").(string))

	log.Printf("[DEBUG] Creating/updating %s", id)

	// getting this ready now because its shared between *everything*, time is... complicated for this resource
	timeOfDay := d.Get("backup.0.time").(string)
	dateOfDay, err := time.Parse(time.RFC3339, fmt.Sprintf("2018-07-30T%s:00Z", timeOfDay))
	if err != nil {
		return fmt.Errorf("generating time from %q for %s: %+v", timeOfDay, id, err)
	}
	times := append(make([]string, 0), date.Time{Time: dateOfDay}.String())

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_backup_policy_vm", id.ID())
		}
	}

	// Less than 7 daily backups is no longer supported for create/update
	if (d.IsNewResource() || d.HasChange("retention_daily.0.count")) && (d.Get("retention_daily.0.count").(int) > 1 && d.Get("retention_daily.0.count").(int) < 7) {
		return fmt.Errorf("The Azure API has recently changed behaviour so that provisioning a `count` for the `retention_daily` field can no longer be less than 7 days for new/updates to existing Backup Policies. Please ensure that `count` is greater than 7, currently %d", d.Get("retention_daily.0.count").(int))
	}

	schedulePolicy, err := expandBackupProtectionPolicyVMSchedule(d, times)
	if err != nil {
		return err
	}

	policyType := protectionpolicies.IAASVMPolicyType(d.Get("policy_type").(string))
	vmProtectionPolicyProperties := &protectionpolicies.AzureIaaSVMProtectionPolicy{
		TimeZone:         utils.String(d.Get("timezone").(string)),
		PolicyType:       pointer.To(policyType),
		SchedulePolicy:   schedulePolicy,
		InstantRPDetails: expandBackupProtectionPolicyVMResourceGroup(d),
		RetentionPolicy: &protectionpolicies.LongTermRetentionPolicy{ // SimpleRetentionPolicy only has duration property ¯\_(ツ)_/¯
			DailySchedule:   expandBackupProtectionPolicyVMRetentionDaily(d, times),
			WeeklySchedule:  expandBackupProtectionPolicyVMRetentionWeekly(d, times),
			MonthlySchedule: expandBackupProtectionPolicyVMRetentionMonthly(d, times),
			YearlySchedule:  expandBackupProtectionPolicyVMRetentionYearly(d, times),
		},
	}

	if d.HasChange("instant_restore_retention_days") {
		days := d.Get("instant_restore_retention_days").(int)
		if protectionpolicies.IAASVMPolicyTypeVOne == policyType && days > 5 {
			return fmt.Errorf("`instant_restore_retention_days` must be less than or equal to `5` when `policy_type` is `V1`")
		}

		vmProtectionPolicyProperties.InstantRpRetentionRangeInDays = pointer.To(int64(days))
	}

	policy := protectionpolicies.ProtectionPolicyResource{
		Properties: vmProtectionPolicyProperties,
	}

	if _, err = client.CreateOrUpdate(ctx, id, policy); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = resourceBackupProtectionPolicyVMWaitForUpdate(ctx, client, id, d); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceBackupProtectionPolicyVMRead(d, meta)
}

func resourceBackupProtectionPolicyVMRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protectionpolicies.ParseBackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading %s", id)

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.Set("name", id.BackupPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	if model := resp.Model; model != nil {
		if properties, ok := model.Properties.(protectionpolicies.AzureIaaSVMProtectionPolicy); ok {
			d.Set("timezone", properties.TimeZone)
			d.Set("instant_restore_retention_days", properties.InstantRpRetentionRangeInDays)

			if schedule, ok := properties.SchedulePolicy.(protectionpolicies.SimpleSchedulePolicy); ok {
				if err := d.Set("backup", flattenBackupProtectionPolicyVMSchedule(schedule)); err != nil {
					return fmt.Errorf("setting `backup`: %+v", err)
				}
			}

			if schedule, ok := properties.SchedulePolicy.(protectionpolicies.SimpleSchedulePolicyV2); ok {
				if err := d.Set("backup", flattenBackupProtectionPolicyVMScheduleV2(schedule)); err != nil {
					return fmt.Errorf("setting `backup`: %+v", err)
				}
			}

			policyType := string(protectionpolicies.IAASVMPolicyTypeVOne)
			if pointer.From(properties.PolicyType) != "" {
				policyType = string(pointer.From(properties.PolicyType))
			}
			d.Set("policy_type", policyType)

			if retention, ok := properties.RetentionPolicy.(protectionpolicies.LongTermRetentionPolicy); ok {
				if s := retention.DailySchedule; s != nil {
					if err := d.Set("retention_daily", flattenBackupProtectionPolicyVMRetentionDaily(s)); err != nil {
						return fmt.Errorf("setting `retention_daily`: %+v", err)
					}
				} else {
					d.Set("retention_daily", nil)
				}

				if s := retention.WeeklySchedule; s != nil {
					if err := d.Set("retention_weekly", flattenBackupProtectionPolicyVMRetentionWeekly(s)); err != nil {
						return fmt.Errorf("setting `retention_weekly`: %+v", err)
					}
				} else {
					d.Set("retention_weekly", nil)
				}

				if s := retention.MonthlySchedule; s != nil {
					if err := d.Set("retention_monthly", flattenBackupProtectionPolicyVMRetentionMonthly(s)); err != nil {
						return fmt.Errorf("setting `retention_monthly`: %+v", err)
					}
				} else {
					d.Set("retention_monthly", nil)
				}

				if s := retention.YearlySchedule; s != nil {
					if err := d.Set("retention_yearly", flattenBackupProtectionPolicyVMRetentionYearly(s)); err != nil {
						return fmt.Errorf("setting `retention_yearly`: %+v", err)
					}
				} else {
					d.Set("retention_yearly", nil)
				}
			}

			if instantRPDetail := properties.InstantRPDetails; instantRPDetail != nil {
				d.Set("instant_restore_resource_group", flattenBackupProtectionPolicyVMResourceGroup(*instantRPDetail))
			}
		}
	}

	return nil
}

func resourceBackupProtectionPolicyVMDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protectionpolicies.ParseBackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", id)

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = resourceBackupProtectionPolicyVMWaitForDeletion(ctx, client, *id, d); err != nil {
		return err
	}

	return nil
}

func expandBackupProtectionPolicyVMSchedule(d *pluginsdk.ResourceData, times []string) (protectionpolicies.SchedulePolicy, error) {
	if bb, ok := d.Get("backup").([]interface{}); ok && len(bb) > 0 {
		block := bb[0].(map[string]interface{})

		policyType := d.Get("policy_type").(string)
		if policyType == string(protectionpolicies.IAASVMPolicyTypeVOne) {
			schedule := protectionpolicies.SimpleSchedulePolicy{ // LongTermSchedulePolicy has no properties
				ScheduleRunTimes: &times,
			}

			if v, ok := block["frequency"].(string); ok {
				schedule.ScheduleRunFrequency = pointer.To(protectionpolicies.ScheduleRunType(v))
			}

			if v, ok := block["weekdays"].(*pluginsdk.Set); ok {
				days := make([]protectionpolicies.DayOfWeek, 0)
				for _, day := range v.List() {
					days = append(days, protectionpolicies.DayOfWeek(day.(string)))
				}
				schedule.ScheduleRunDays = &days
			}

			return schedule, nil
		} else {
			frequency := block["frequency"].(string)
			schedule := protectionpolicies.SimpleSchedulePolicyV2{
				ScheduleRunFrequency: pointer.To(protectionpolicies.ScheduleRunType(frequency)),
			}

			switch frequency {
			case string(protectionpolicies.ScheduleRunTypeHourly):
				interval, ok := block["hour_interval"].(int)
				if !ok {
					return nil, fmt.Errorf("`hour_interval` must be specified when `backup.0.frequency` is `Hourly`")
				}

				duration, ok := block["hour_duration"].(int)
				if !ok {
					return nil, fmt.Errorf("`hour_duration` must be specified when `backup.0.frequency` is `Hourly`")
				}

				if interval == 0 && duration == 0 {
					return nil, fmt.Errorf("`hour_interval` and `hour_duration` must be specified when `backup.0.frequency` is `Hourly`")
				}
				if interval == 0 {
					return nil, fmt.Errorf("`hour_interval` must be specified when `backup.0.frequency` is `Hourly`")
				}
				if duration == 0 {
					return nil, fmt.Errorf("`hour_duration` must be specified when `backup.0.frequency` is `Hourly`")
				}

				if duration%interval != 0 {
					return nil, fmt.Errorf("`hour_duration` must be multiplier of `hour_interval`")
				}

				schedule.HourlySchedule = &protectionpolicies.HourlySchedule{
					Interval:                pointer.To(int64(interval)),
					ScheduleWindowStartTime: &times[0],
					ScheduleWindowDuration:  pointer.To(int64(duration)),
				}
			case string(protectionpolicies.ScheduleRunTypeDaily):
				schedule.DailySchedule = &protectionpolicies.DailySchedule{
					ScheduleRunTimes: &times,
				}
			case string(protectionpolicies.ScheduleRunTypeWeekly):
				weekDays, ok := block["weekdays"].(*pluginsdk.Set)
				if !ok {
					return nil, fmt.Errorf("`weekdays` must be specified when `backup.0.frequency` is `Weekly`")
				}

				days := make([]protectionpolicies.DayOfWeek, 0)
				for _, day := range weekDays.List() {
					days = append(days, protectionpolicies.DayOfWeek(day.(string)))
				}

				schedule.WeeklySchedule = &protectionpolicies.WeeklySchedule{
					ScheduleRunDays:  &days,
					ScheduleRunTimes: &times,
				}
			default:
				return nil, fmt.Errorf("Unrecognized value for backup.0.frequency")
			}

			return schedule, nil
		}
	}

	return nil, nil
}

func expandBackupProtectionPolicyVMResourceGroup(d *pluginsdk.ResourceData) *protectionpolicies.InstantRPAdditionalDetails {
	if v, ok := d.Get("instant_restore_resource_group").([]interface{}); ok && len(v) > 0 {
		rgRaw := v[0].(map[string]interface{})
		output := &protectionpolicies.InstantRPAdditionalDetails{
			AzureBackupRGNamePrefix: utils.String(rgRaw["prefix"].(string)),
		}
		if suffix, ok := rgRaw["suffix"]; ok && suffix != "" {
			output.AzureBackupRGNameSuffix = utils.String(suffix.(string))
		}
		return output
	}

	return nil
}

func expandBackupProtectionPolicyVMRetentionDaily(d *pluginsdk.ResourceData, times []string) *protectionpolicies.DailyRetentionSchedule {
	if rb, ok := d.Get("retention_daily").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		return &protectionpolicies.DailyRetentionSchedule{
			RetentionTimes: &times,
			RetentionDuration: &protectionpolicies.RetentionDuration{
				Count:        pointer.To(int64(block["count"].(int))),
				DurationType: pointer.To(protectionpolicies.RetentionDurationTypeDays),
			},
		}
	}

	return nil
}

func expandBackupProtectionPolicyVMRetentionWeekly(d *pluginsdk.ResourceData, times []string) *protectionpolicies.WeeklyRetentionSchedule {
	if rb, ok := d.Get("retention_weekly").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		retention := protectionpolicies.WeeklyRetentionSchedule{
			RetentionTimes: &times,
			RetentionDuration: &protectionpolicies.RetentionDuration{
				Count:        pointer.To(int64(block["count"].(int))),
				DurationType: pointer.To(protectionpolicies.RetentionDurationTypeWeeks),
			},
		}

		if v, ok := block["weekdays"].(*pluginsdk.Set); ok {
			days := make([]protectionpolicies.DayOfWeek, 0)
			for _, day := range v.List() {
				days = append(days, protectionpolicies.DayOfWeek(day.(string)))
			}
			retention.DaysOfTheWeek = &days
		}

		return &retention
	}

	return nil
}

func expandBackupProtectionPolicyVMRetentionMonthly(d *pluginsdk.ResourceData, times []string) *protectionpolicies.MonthlyRetentionSchedule {
	if rb, ok := d.Get("retention_monthly").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		scheduleFormat := protectionpolicies.RetentionScheduleFormatWeekly
		var weekly *protectionpolicies.WeeklyRetentionFormat = nil
		var daily *protectionpolicies.DailyRetentionFormat = nil
		if v, ok := block["days"]; ok && v.(*pluginsdk.Set).Len() > 0 {
			scheduleFormat = protectionpolicies.RetentionScheduleFormatDaily
			daily = expandBackupProtectionPolicyVMRetentionDailyFormat(block)
		} else {
			weekly = expandBackupProtectionPolicyVMRetentionWeeklyFormat(block)
		}

		retention := protectionpolicies.MonthlyRetentionSchedule{
			RetentionScheduleFormatType: &scheduleFormat,
			RetentionScheduleDaily:      daily,
			RetentionScheduleWeekly:     weekly,
			RetentionTimes:              &times,
			RetentionDuration: &protectionpolicies.RetentionDuration{
				Count:        pointer.To(int64(block["count"].(int))),
				DurationType: pointer.To(protectionpolicies.RetentionDurationTypeMonths),
			},
		}

		return &retention
	}

	return nil
}

func expandBackupProtectionPolicyVMRetentionYearly(d *pluginsdk.ResourceData, times []string) *protectionpolicies.YearlyRetentionSchedule {
	if rb, ok := d.Get("retention_yearly").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		scheduleFormat := protectionpolicies.RetentionScheduleFormatWeekly
		var weekly *protectionpolicies.WeeklyRetentionFormat = nil
		var daily *protectionpolicies.DailyRetentionFormat = nil
		if v, ok := block["days"]; ok && v.(*pluginsdk.Set).Len() > 0 {
			scheduleFormat = protectionpolicies.RetentionScheduleFormatDaily
			daily = expandBackupProtectionPolicyVMRetentionDailyFormat(block)
		} else {
			weekly = expandBackupProtectionPolicyVMRetentionWeeklyFormat(block)
		}

		retention := protectionpolicies.YearlyRetentionSchedule{
			RetentionScheduleFormatType: &scheduleFormat,
			RetentionScheduleDaily:      daily,
			RetentionScheduleWeekly:     weekly,
			RetentionTimes:              &times,
			RetentionDuration: &protectionpolicies.RetentionDuration{
				Count:        pointer.To(int64(block["count"].(int))),
				DurationType: pointer.To(protectionpolicies.RetentionDurationTypeYears),
			},
		}

		if v, ok := block["months"].(*pluginsdk.Set); ok {
			months := make([]protectionpolicies.MonthOfYear, 0)
			for _, month := range v.List() {
				months = append(months, protectionpolicies.MonthOfYear(month.(string)))
			}
			retention.MonthsOfYear = &months
		}

		return &retention
	}

	return nil
}

func expandBackupProtectionPolicyVMRetentionWeeklyFormat(block map[string]interface{}) *protectionpolicies.WeeklyRetentionFormat {
	weekly := protectionpolicies.WeeklyRetentionFormat{}

	if v, ok := block["weekdays"].(*pluginsdk.Set); ok {
		days := make([]protectionpolicies.DayOfWeek, 0)
		for _, day := range v.List() {
			days = append(days, protectionpolicies.DayOfWeek(day.(string)))
		}
		weekly.DaysOfTheWeek = &days
	}

	if v, ok := block["weeks"].(*pluginsdk.Set); ok {
		weeks := make([]protectionpolicies.WeekOfMonth, 0)
		for _, week := range v.List() {
			weeks = append(weeks, protectionpolicies.WeekOfMonth(week.(string)))
		}
		weekly.WeeksOfTheMonth = &weeks
	}

	return &weekly
}

func expandBackupProtectionPolicyVMRetentionDailyFormat(block map[string]interface{}) *protectionpolicies.DailyRetentionFormat {
	days := make([]protectionpolicies.Day, 0)

	if block["include_last_days"].(bool) {
		days = append(days, protectionpolicies.Day{
			Date:   pointer.To(int64(0)),
			IsLast: pointer.To(true),
		})
	}

	if v, ok := block["days"].(*pluginsdk.Set); ok {
		for _, day := range v.List() {
			days = append(days, protectionpolicies.Day{
				Date:   pointer.To(int64(day.(int))),
				IsLast: pointer.To(false),
			})
		}
	}

	daily := protectionpolicies.DailyRetentionFormat{
		DaysOfTheMonth: &days,
	}

	return &daily
}

func flattenBackupProtectionPolicyVMResourceGroup(rpDetail protectionpolicies.InstantRPAdditionalDetails) []interface{} {
	if rpDetail.AzureBackupRGNamePrefix == nil {
		return nil
	}

	block := map[string]interface{}{}

	prefix := ""
	if rpDetail.AzureBackupRGNamePrefix != nil {
		prefix = *rpDetail.AzureBackupRGNamePrefix
	}
	block["prefix"] = prefix

	suffix := ""
	if rpDetail.AzureBackupRGNameSuffix != nil {
		suffix = *rpDetail.AzureBackupRGNameSuffix
	}
	block["suffix"] = suffix

	return []interface{}{block}
}

func flattenBackupProtectionPolicyVMSchedule(schedule protectionpolicies.SimpleSchedulePolicy) []interface{} {
	block := map[string]interface{}{}

	block["frequency"] = string(pointer.From(schedule.ScheduleRunFrequency))

	if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
		policyTime, _ := time.Parse(time.RFC3339, (*times)[0])
		block["time"] = policyTime.Format("15:04")
	}

	if days := schedule.ScheduleRunDays; days != nil {
		weekdays := make([]interface{}, 0)
		for _, d := range *days {
			weekdays = append(weekdays, string(d))
		}
		block["weekdays"] = pluginsdk.NewSet(pluginsdk.HashString, weekdays)
	}

	return []interface{}{block}
}

func flattenBackupProtectionPolicyVMScheduleV2(schedule protectionpolicies.SimpleSchedulePolicyV2) []interface{} {
	block := map[string]interface{}{}

	frequency := pointer.From(schedule.ScheduleRunFrequency)
	block["frequency"] = string(frequency)

	switch frequency {
	case protectionpolicies.ScheduleRunTypeHourly:
		schedule := schedule.HourlySchedule
		if schedule.Interval != nil {
			block["hour_interval"] = *schedule.Interval
		}

		if schedule.ScheduleWindowDuration != nil {
			block["hour_duration"] = *schedule.ScheduleWindowDuration
		}

		if schedule.ScheduleWindowStartTime != nil {
			policyTime, _ := time.Parse(time.RFC3339, pointer.From(schedule.ScheduleWindowStartTime))
			block["time"] = policyTime.Format("15:04")
		}
	case protectionpolicies.ScheduleRunTypeDaily:
		schedule := schedule.DailySchedule
		if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
			policyTime, _ := time.Parse(time.RFC3339, (*times)[0])
			block["time"] = policyTime.Format("15:04")
		}
	case protectionpolicies.ScheduleRunTypeWeekly:
		schedule := schedule.WeeklySchedule
		if days := schedule.ScheduleRunDays; days != nil {
			weekdays := make([]interface{}, 0)
			for _, d := range *days {
				weekdays = append(weekdays, string(d))
			}
			block["weekdays"] = pluginsdk.NewSet(pluginsdk.HashString, weekdays)
		}

		if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
			policyTime, _ := time.Parse(time.RFC3339, (*times)[0])
			block["time"] = policyTime.Format("15:04")
		}
	default:
	}

	return []interface{}{block}
}

func flattenBackupProtectionPolicyVMRetentionDaily(daily *protectionpolicies.DailyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := daily.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	return []interface{}{block}
}

func flattenBackupProtectionPolicyVMRetentionWeekly(weekly *protectionpolicies.WeeklyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := weekly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if days := weekly.DaysOfTheWeek; days != nil {
		weekdays := make([]interface{}, 0)
		for _, d := range *days {
			weekdays = append(weekdays, string(d))
		}
		block["weekdays"] = pluginsdk.NewSet(pluginsdk.HashString, weekdays)
	}

	return []interface{}{block}
}

func flattenBackupProtectionPolicyVMRetentionMonthly(monthly *protectionpolicies.MonthlyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := monthly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if weekly := monthly.RetentionScheduleWeekly; weekly != nil {
		block["weekdays"], block["weeks"] = flattenBackupProtectionPolicyVMRetentionWeeklyFormat(weekly)
	}

	if daily := monthly.RetentionScheduleDaily; daily != nil {
		block["days"], block["include_last_days"] = flattenBackupProtectionPolicyVMRetentionDailyFormat(daily)
	}

	return []interface{}{block}
}

func flattenBackupProtectionPolicyVMRetentionYearly(yearly *protectionpolicies.YearlyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := yearly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if weekly := yearly.RetentionScheduleWeekly; weekly != nil {
		block["weekdays"], block["weeks"] = flattenBackupProtectionPolicyVMRetentionWeeklyFormat(weekly)
	}

	if months := yearly.MonthsOfYear; months != nil {
		slice := make([]interface{}, 0)
		for _, d := range *months {
			slice = append(slice, string(d))
		}
		block["months"] = pluginsdk.NewSet(pluginsdk.HashString, slice)
	}

	if daily := yearly.RetentionScheduleDaily; daily != nil {
		block["days"], block["include_last_days"] = flattenBackupProtectionPolicyVMRetentionDailyFormat(daily)
	}

	return []interface{}{block}
}

func flattenBackupProtectionPolicyVMRetentionWeeklyFormat(retention *protectionpolicies.WeeklyRetentionFormat) (weekdays, weeks *pluginsdk.Set) {
	if days := retention.DaysOfTheWeek; days != nil {
		slice := make([]interface{}, 0)
		for _, d := range *days {
			slice = append(slice, string(d))
		}
		weekdays = pluginsdk.NewSet(pluginsdk.HashString, slice)
	}

	if days := retention.WeeksOfTheMonth; days != nil {
		slice := make([]interface{}, 0)
		for _, d := range *days {
			slice = append(slice, string(d))
		}
		weeks = pluginsdk.NewSet(pluginsdk.HashString, slice)
	}

	return weekdays, weeks
}

func flattenBackupProtectionPolicyVMRetentionDailyFormat(retention *protectionpolicies.DailyRetentionFormat) (days []interface{}, includeLastDay bool) {
	if dotm := retention.DaysOfTheMonth; dotm != nil {
		for _, d := range *dotm {
			if d.Date != nil && *d.Date != 0 {
				days = append(days, *d.Date)
			}
			if d.IsLast != nil && *d.IsLast {
				includeLastDay = true
			}
		}
	}

	return days, includeLastDay
}

func resourceBackupProtectionPolicyVMWaitForUpdate(ctx context.Context, client *protectionpolicies.ProtectionPoliciesClient, id protectionpolicies.BackupPolicyId, d *pluginsdk.ResourceData) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"NotFound"},
		Target:     []string{"Found"},
		Refresh:    resourceBackupProtectionPolicyVMRefreshFunc(ctx, client, id),
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	_, err := state.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s to provision: %+v", id, err)
	}

	return nil
}

func resourceBackupProtectionPolicyVMWaitForDeletion(ctx context.Context, client *protectionpolicies.ProtectionPoliciesClient, id protectionpolicies.BackupPolicyId, d *pluginsdk.ResourceData) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"Found"},
		Target:     []string{"NotFound"},
		Refresh:    resourceBackupProtectionPolicyVMRefreshFunc(ctx, client, id),
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
	}

	_, err := state.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s to provision: %+v", id, err)
	}

	return nil
}

func resourceBackupProtectionPolicyVMRefreshFunc(ctx context.Context, client *protectionpolicies.ProtectionPoliciesClient, id protectionpolicies.BackupPolicyId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "NotFound", nil
			}

			return resp, "Error", fmt.Errorf("making Read request on %s: %+v", id, err)
		}

		return resp, "Found", nil
	}
}

func resourceBackupProtectionPolicyVMSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z][-_!a-zA-Z0-9]{2,149}$"),
				"Backup Policy name must be 3 - 150 characters long, start with a letter, contain only letters and numbers.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"instant_restore_retention_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(1, 30),
		},

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"timezone": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "UTC",
		},

		"instant_restore_resource_group": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"prefix": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"suffix": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"backup": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"frequency": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(protectionpolicies.ScheduleRunTypeHourly),
							string(protectionpolicies.ScheduleRunTypeDaily),
							string(protectionpolicies.ScheduleRunTypeWeekly),
						}, false),
					},

					"time": { // applies to all backup schedules & retention times (they all must be the same)
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
							"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
						),
					},

					"weekdays": { // only for weekly
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
						},
					},

					"hour_interval": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						ValidateFunc: validation.IntInSlice([]int{
							4,
							6,
							8,
							12,
						}),
					},

					"hour_duration": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(4, 24),
					},
				},
			},
		},

		"policy_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(protectionpolicies.IAASVMPolicyTypeVOne),
			ValidateFunc: validation.StringInSlice([]string{
				string(protectionpolicies.IAASVMPolicyTypeVOne),
				string(protectionpolicies.IAASVMPolicyTypeVTwo),
			}, false),
		},

		"retention_daily": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 9999), // Azure no longer supports less than 7 daily backups. This should be updated in 3.0 provider

					},
				},
			},
		},

		"retention_weekly": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 9999),
					},

					"weekdays": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Set:      set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
						},
					},
				},
			},
		},

		"retention_monthly": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 9999),
					},

					"weeks": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(protectionpolicies.WeekOfMonthFirst),
								string(protectionpolicies.WeekOfMonthSecond),
								string(protectionpolicies.WeekOfMonthThird),
								string(protectionpolicies.WeekOfMonthFourth),
								string(protectionpolicies.WeekOfMonthLast),
							}, false),
						},
						ConflictsWith: []string{
							"retention_monthly.0.days",
							"retention_monthly.0.include_last_days",
						},
						AtLeastOneOf: []string{
							"retention_monthly.0.weekdays",
							"retention_monthly.0.weeks",
							"retention_monthly.0.days",
							"retention_monthly.0.include_last_days",
						},
						RequiredWith: []string{
							"retention_monthly.0.weekdays",
						},
					},

					"weekdays": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
						},
						RequiredWith: []string{
							"retention_monthly.0.weeks",
						},
						ConflictsWith: []string{
							"retention_monthly.0.days",
							"retention_monthly.0.include_last_days",
						},
						AtLeastOneOf: []string{
							"retention_monthly.0.weekdays",
							"retention_monthly.0.weeks",
							"retention_monthly.0.days",
							"retention_monthly.0.include_last_days",
						},
					},

					"days": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(1, 31), // days in months
						},
						ConflictsWith: []string{
							"retention_monthly.0.weeks",
							"retention_monthly.0.weekdays",
						},
						AtLeastOneOf: []string{
							"retention_monthly.0.weekdays",
							"retention_monthly.0.weeks",
							"retention_monthly.0.days",
							"retention_monthly.0.include_last_days",
						},
					},

					"include_last_days": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
						ConflictsWith: []string{
							"retention_monthly.0.weeks",
							"retention_monthly.0.weekdays",
						},
						AtLeastOneOf: []string{
							"retention_monthly.0.weekdays",
							"retention_monthly.0.weeks",
							"retention_monthly.0.days",
							"retention_monthly.0.include_last_days",
						},
					},
				},
			},
		},

		"retention_yearly": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 9999),
					},

					"months": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Set:      set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsMonth(true),
						},
					},

					"weeks": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(protectionpolicies.WeekOfMonthFirst),
								string(protectionpolicies.WeekOfMonthSecond),
								string(protectionpolicies.WeekOfMonthThird),
								string(protectionpolicies.WeekOfMonthFourth),
								string(protectionpolicies.WeekOfMonthLast),
							}, false),
						},
						RequiredWith: []string{
							"retention_yearly.0.weekdays",
						},
						ConflictsWith: []string{
							"retention_yearly.0.days",
							"retention_yearly.0.include_last_days",
						},
						AtLeastOneOf: []string{
							"retention_yearly.0.weeks",
							"retention_yearly.0.weekdays",
							"retention_yearly.0.days",
							"retention_yearly.0.include_last_days",
						},
					},

					"weekdays": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      set.HashStringIgnoreCase,
						Elem: &pluginsdk.Schema{
							Type:             pluginsdk.TypeString,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
						},
						RequiredWith: []string{
							"retention_yearly.0.weeks",
						},
						ConflictsWith: []string{
							"retention_yearly.0.days",
							"retention_yearly.0.include_last_days",
						},
						AtLeastOneOf: []string{
							"retention_yearly.0.weeks",
							"retention_yearly.0.weekdays",
							"retention_yearly.0.days",
							"retention_yearly.0.include_last_days",
						},
					},

					"days": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(1, 31), // days in months
						},
						ConflictsWith: []string{
							"retention_yearly.0.weeks",
							"retention_yearly.0.weekdays",
						},
						AtLeastOneOf: []string{
							"retention_yearly.0.weeks",
							"retention_yearly.0.weekdays",
							"retention_yearly.0.days",
							"retention_yearly.0.include_last_days",
						},
					},

					"include_last_days": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
						ConflictsWith: []string{
							"retention_yearly.0.weeks",
							"retention_yearly.0.weekdays",
						},
						AtLeastOneOf: []string{
							"retention_yearly.0.weeks",
							"retention_yearly.0.weekdays",
							"retention_yearly.0.days",
							"retention_yearly.0.include_last_days",
						},
					},
				},
			},
		},
	}
}
