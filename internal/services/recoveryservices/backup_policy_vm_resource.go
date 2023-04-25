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
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/protectionpolicies"
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

	var retentionPolicy protectionpolicies.RetentionPolicy = protectionpolicies.LongTermRetentionPolicy{ // SimpleRetentionPolicy only has duration property ¯\_(ツ)_/¯
		DailySchedule:   expandBackupProtectionPolicyVMRetentionDaily(d, times),
		WeeklySchedule:  expandBackupProtectionPolicyVMRetentionWeekly(d, times),
		MonthlySchedule: expandBackupProtectionPolicyVMRetentionMonthly(d, times),
		YearlySchedule:  expandBackupProtectionPolicyVMRetentionYearly(d, times),
	}
	policyType := protectionpolicies.IAASVMPolicyType(d.Get("policy_type").(string))
	vmProtectionPolicyProperties := &protectionpolicies.AzureIaaSVMProtectionPolicy{
		TimeZone:         utils.String(d.Get("timezone").(string)),
		PolicyType:       pointer.To(policyType),
		SchedulePolicy:   pointer.To(schedulePolicy),
		InstantRPDetails: expandBackupProtectionPolicyVMResourceGroup(d),
		RetentionPolicy:  pointer.To(retentionPolicy),
	}

	if d.HasChange("instant_restore_retention_days") {
		days := d.Get("instant_restore_retention_days").(int)
		if protectionpolicies.IAASVMPolicyTypeVOne == policyType && days > 5 {
			return fmt.Errorf("`instant_restore_retention_days` must be less than or equal to `5` when `policy_type` is `V1`")
		}

		vmProtectionPolicyProperties.InstantRpRetentionRangeInDays = pointer.To(int64(days))
	}

	var policyProps protectionpolicies.ProtectionPolicy = vmProtectionPolicyProperties
	policy := protectionpolicies.ProtectionPolicyResource{
		Properties: pointer.To(policyProps),
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
		if props := model.Properties; props != nil {
			if policy, ok := (*props).(protectionpolicies.AzureIaaSVMProtectionPolicy); ok {
				d.Set("instant_restore_retention_days", policy.InstantRpRetentionRangeInDays)
				d.Set("timezone", policy.TimeZone)

				if schedulePolicyPtr := policy.SchedulePolicy; schedulePolicyPtr != nil {
					schedulePolicy := *schedulePolicyPtr
					if _, ok = schedulePolicy.(protectionpolicies.SimpleSchedulePolicy); ok {
						if err = d.Set("backup", flattenBackupProtectionPolicyVMSchedule(policy.SchedulePolicy)); err != nil {
							return fmt.Errorf("setting `backup`: %+v", err)
						}
					}

					if _, ok = schedulePolicy.(protectionpolicies.SimpleSchedulePolicyV2); ok {
						if err = d.Set("backup", flattenBackupProtectionPolicyVMScheduleV2(policy.SchedulePolicy)); err != nil {
							return fmt.Errorf("setting `backup`: %+v", err)
						}
					}
				}

				policyType := string(protectionpolicies.IAASVMPolicyTypeVOne)
				if pointer.From(policy.PolicyType) != "" {
					policyType = string(pointer.From(policy.PolicyType))
				}
				d.Set("policy_type", policyType)

				if err := d.Set("retention_daily", flattenBackupProtectionPolicyVMRetentionDaily(policy.RetentionPolicy)); err != nil {
					return fmt.Errorf("setting `retention_daily`: %+v", err)
				}
				if err := d.Set("retention_monthly", flattenBackupProtectionPolicyVMRetentionMonthly(policy.RetentionPolicy)); err != nil {
					return fmt.Errorf("setting `retention_monthly`: %+v", err)
				}
				if err := d.Set("retention_weekly", flattenBackupProtectionPolicyVMRetentionWeekly(policy.RetentionPolicy)); err != nil {
					return fmt.Errorf("setting `retention_weekly`: %+v", err)
				}
				if err := d.Set("retention_yearly", flattenBackupProtectionPolicyVMRetentionYearly(policy.RetentionPolicy)); err != nil {
					return fmt.Errorf("setting `retention_yearly`: %+v", err)
				}

				if err := d.Set("instant_restore_resource_group", flattenBackupProtectionPolicyVMResourceGroup(policy.InstantRPDetails)); err != nil {
					return fmt.Errorf("setting `instant_restore_resource_group`: %+v", err)
				}
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

		retention := protectionpolicies.MonthlyRetentionSchedule{
			RetentionScheduleFormatType: pointer.To(protectionpolicies.RetentionScheduleFormatWeekly), // this is always weekly ¯\_(ツ)_/¯
			RetentionScheduleDaily:      nil,                                                          // and this is always nil..
			RetentionScheduleWeekly:     expandBackupProtectionPolicyVMRetentionWeeklyFormat(block),
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

		retention := protectionpolicies.YearlyRetentionSchedule{
			RetentionScheduleFormatType: pointer.To(protectionpolicies.RetentionScheduleFormatWeekly), // this is always weekly ¯\_(ツ)_/¯
			RetentionScheduleDaily:      nil,                                                          // and this is always nil..
			RetentionScheduleWeekly:     expandBackupProtectionPolicyVMRetentionWeeklyFormat(block),
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

func flattenBackupProtectionPolicyVMResourceGroup(input *protectionpolicies.InstantRPAdditionalDetails) []interface{} {
	output := make([]interface{}, 0)

	if input != nil && input.AzureBackupRGNamePrefix != nil {
		prefix := ""
		if input.AzureBackupRGNamePrefix != nil {
			prefix = *input.AzureBackupRGNamePrefix
		}

		suffix := ""
		if input.AzureBackupRGNameSuffix != nil {
			suffix = *input.AzureBackupRGNameSuffix
		}

		output = append(output, map[string]interface{}{
			"prefix": prefix,
			"suffix": suffix,
		})
	}

	return output
}

func flattenBackupProtectionPolicyVMSchedule(input *protectionpolicies.SchedulePolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if v, ok := (*input).(protectionpolicies.SimpleSchedulePolicy); ok {
			timeValue := ""
			if times := v.ScheduleRunTimes; times != nil && len(*times) > 0 {
				policyTime, err := time.Parse(time.RFC3339, (*times)[0])
				if err == nil {
					timeValue = policyTime.Format("15:04")
				}
			}

			weekdays := make([]interface{}, 0)
			if days := v.ScheduleRunDays; days != nil {
				for _, d := range *days {
					weekdays = append(weekdays, string(d))
				}
			}

			output = append(output, map[string]interface{}{
				"frequency": string(pointer.From(v.ScheduleRunFrequency)),
				"time":      timeValue,
				"weekdays":  pluginsdk.NewSet(pluginsdk.HashString, weekdays),
			})
		}
	}

	return output
}

func flattenBackupProtectionPolicyVMScheduleV2(input *protectionpolicies.SchedulePolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if v, ok := (*input).(protectionpolicies.SimpleSchedulePolicyV2); ok && v.ScheduleRunFrequency != nil {
			frequency := *v.ScheduleRunFrequency

			hourDuration := 0
			hourInterval := 0
			timeValue := ""
			weekdays := make([]interface{}, 0)

			switch frequency {
			case protectionpolicies.ScheduleRunTypeHourly:
				{
					schedule := v.HourlySchedule
					if schedule.Interval != nil {
						hourInterval = int(*schedule.Interval)
					}

					if schedule.ScheduleWindowDuration != nil {
						hourDuration = int(*schedule.ScheduleWindowDuration)
					}

					if schedule.ScheduleWindowStartTime != nil {
						startTime, err := time.Parse(time.RFC3339, pointer.From(schedule.ScheduleWindowStartTime))
						if err == nil {
							timeValue = startTime.Format("15:04")
						}
					}
				}

			case protectionpolicies.ScheduleRunTypeDaily:
				{
					schedule := v.DailySchedule
					if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
						policyTime, err := time.Parse(time.RFC3339, (*times)[0])
						if err == nil {
							timeValue = policyTime.Format("15:04")
						}
					}
				}

			case protectionpolicies.ScheduleRunTypeWeekly:
				{
					schedule := v.WeeklySchedule
					if days := schedule.ScheduleRunDays; days != nil {
						for _, d := range *days {
							weekdays = append(weekdays, string(d))
						}
					}

					if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
						policyTime, err := time.Parse(time.RFC3339, (*times)[0])
						if err == nil {
							timeValue = policyTime.Format("15:04")
						}
					}
				}
			}

			output = append(output, map[string]interface{}{
				"frequency":     string(*v.ScheduleRunFrequency),
				"hour_duration": hourDuration,
				"hour_interval": hourInterval,
				"time":          timeValue,
				"weekdays":      pluginsdk.NewSet(pluginsdk.HashString, weekdays),
			})
		}
	}

	return output
}

func flattenBackupProtectionPolicyVMRetentionDaily(input *protectionpolicies.RetentionPolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if retention, ok := (*input).(protectionpolicies.LongTermRetentionPolicy); ok {
			if daily := retention.DailySchedule; daily != nil {
				count := 0
				if daily.RetentionDuration != nil && daily.RetentionDuration.Count != nil {
					count = int(*daily.RetentionDuration.Count)
				}
				output = append(output, map[string]interface{}{
					"count": count,
				})
			}
		}
	}

	return output
}

func flattenBackupProtectionPolicyVMRetentionWeekly(input *protectionpolicies.RetentionPolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if retention, ok := (*input).(protectionpolicies.LongTermRetentionPolicy); ok {
			if weekly := retention.WeeklySchedule; weekly != nil {
				count := 0
				if weekly.RetentionDuration != nil && weekly.RetentionDuration.Count != nil {
					count = int(*weekly.RetentionDuration.Count)
				}

				weekdays := make([]interface{}, 0)
				if dotw := weekly.DaysOfTheWeek; dotw != nil {
					for _, d := range *dotw {
						weekdays = append(weekdays, string(d))
					}
				}

				output = append(output, map[string]interface{}{
					"count":    count,
					"weekdays": pluginsdk.NewSet(pluginsdk.HashString, weekdays),
				})
			}
		}
	}

	return output
}

func flattenBackupProtectionPolicyVMRetentionMonthly(input *protectionpolicies.RetentionPolicy) []interface{} {
	output := make([]interface{}, 0)
	if input != nil {
		if retention, ok := (*input).(protectionpolicies.LongTermRetentionPolicy); ok {
			if monthly := retention.MonthlySchedule; monthly != nil {
				count := 0
				if monthly.RetentionDuration != nil && monthly.RetentionDuration.Count != nil {
					count = int(*monthly.RetentionDuration.Count)
				}

				weekdays := make([]interface{}, 0)
				weeks := make([]interface{}, 0)
				if weekly := monthly.RetentionScheduleWeekly; weekly != nil {
					if dotw := weekly.DaysOfTheWeek; dotw != nil {
						for _, d := range *dotw {
							weekdays = append(weekdays, string(d))
						}
					}

					if wotm := weekly.WeeksOfTheMonth; wotm != nil {
						for _, d := range *wotm {
							weeks = append(weeks, string(d))
						}
					}
				}

				output = append(output, map[string]interface{}{
					"count":    count,
					"weekdays": pluginsdk.NewSet(pluginsdk.HashString, weekdays),
					"weeks":    pluginsdk.NewSet(pluginsdk.HashString, weeks),
				})
			}
		}
	}
	return output
}

func flattenBackupProtectionPolicyVMRetentionYearly(input *protectionpolicies.RetentionPolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if retention, ok := (*input).(protectionpolicies.LongTermRetentionPolicy); ok {
			if yearly := retention.YearlySchedule; yearly != nil {
				count := 0
				if yearly.RetentionDuration != nil && yearly.RetentionDuration.Count != nil {
					count = int(*yearly.RetentionDuration.Count)
				}
				weekdays := make([]interface{}, 0)
				weeks := make([]interface{}, 0)
				if weekly := yearly.RetentionScheduleWeekly; weekly != nil {
					if dotw := weekly.DaysOfTheWeek; dotw != nil {
						for _, d := range *dotw {
							weekdays = append(weekdays, string(d))
						}
					}

					if wotm := weekly.WeeksOfTheMonth; wotm != nil {
						for _, d := range *wotm {
							weeks = append(weeks, string(d))
						}
					}
				}

				months := make([]interface{}, 0)
				if yearly.MonthsOfYear != nil {
					for _, v := range *yearly.MonthsOfYear {
						months = append(months, string(v))
					}
				}

				output = append(output, map[string]interface{}{
					"count":    count,
					"months":   pluginsdk.NewSet(pluginsdk.HashString, months),
					"weekdays": pluginsdk.NewSet(pluginsdk.HashString, weekdays),
					"weeks":    pluginsdk.NewSet(pluginsdk.HashString, weeks),
				})
			}
		}
	}

	return output
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
						Required: true,
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
						Required: true,
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
	}
}
