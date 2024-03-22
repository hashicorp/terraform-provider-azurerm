// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectionPolicyFileShare() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBackupProtectionPolicyFileShareCreateUpdate,
		Read:   resourceBackupProtectionPolicyFileShareRead,
		Update: resourceBackupProtectionPolicyFileShareCreateUpdate,
		Delete: resourceBackupProtectionPolicyFileShareDelete,

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

		Schema: resourceBackupProtectionPolicyFileShareSchema(),

		// if daily, we need daily retention
		// if weekly daily cannot be set, and we need weekly
		CustomizeDiff: func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
			_, hasDaily := diff.GetOk("retention_daily")
			_, hasWeekly := diff.GetOk("retention_weekly")

			frequencyI, _ := diff.GetOk("backup.0.frequency")
			switch strings.ToLower(frequencyI.(string)) {
			case "daily":
				if !hasDaily {
					return fmt.Errorf("`retention_daily` must be set when backup.0.frequency is daily")
				}

				if _, ok := diff.GetOk("backup.0.weekdays"); ok {
					return fmt.Errorf("`backup.0.weekdays` should be not set when backup.0.frequency is daily")
				}
			case "weekly":
				if hasDaily {
					return fmt.Errorf("`retention_daily` must be not set when backup.0.frequency is weekly")
				}
				if !hasWeekly {
					return fmt.Errorf("`retention_weekly` must be set when backup.0.frequency is weekly")
				}
			case "hourly":
				if !hasDaily {
					return fmt.Errorf("`retention_daily` must be set when backup.0.frequency is hourly")
				}
			default:
				return fmt.Errorf("Unrecognized value for backup.0.frequency")
			}
			return nil
		},
	}
}

func resourceBackupProtectionPolicyFileShareCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := protectionpolicies.NewBackupPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), d.Get("name").(string))

	log.Printf("[DEBUG] Creating/updating %s", id)

	// getting this ready now because its shared between *everything*, time is... complicated for this resource
	// if it's using hourly backup schedule, it passes null in times
	var times []string
	if strings.ToLower(d.Get("backup.0.frequency").(string)) == "daily" {
		timeOfDay := d.Get("backup.0.time").(string)
		dateOfDay, err := time.Parse(time.RFC3339, fmt.Sprintf("2018-07-30T%s:00Z", timeOfDay))
		if err != nil {
			return fmt.Errorf("generating time from %q for %s: %+v", timeOfDay, id, err)
		}
		times = append(make([]string, 0), date.Time{Time: dateOfDay}.String())
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_backup_policy_file_share", id.ID())
		}
	}

	AzureFileShareProtectionPolicyProperties := &protectionpolicies.AzureFileShareProtectionPolicy{
		TimeZone:       pointer.To(d.Get("timezone").(string)),
		WorkLoadType:   pointer.To(protectionpolicies.WorkloadTypeAzureFileShare),
		SchedulePolicy: expandBackupProtectionPolicyFileShareSchedule(d, times),
		RetentionPolicy: &protectionpolicies.LongTermRetentionPolicy{ // SimpleRetentionPolicy only has duration property ¯\_(ツ)_/¯
			DailySchedule:   expandBackupProtectionPolicyFileShareRetentionDaily(d, times),
			WeeklySchedule:  expandBackupProtectionPolicyFileShareRetentionWeekly(d, times),
			MonthlySchedule: expandBackupProtectionPolicyFileShareRetentionMonthly(d, times),
			YearlySchedule:  expandBackupProtectionPolicyFileShareRetentionYearly(d, times),
		},
	}

	policy := protectionpolicies.ProtectionPolicyResource{
		Properties: AzureFileShareProtectionPolicyProperties,
	}

	if _, err := client.CreateOrUpdate(ctx, id, policy); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := resourceBackupProtectionPolicyFileShareWaitForUpdate(ctx, client, id, d); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceBackupProtectionPolicyFileShareRead(d, meta)
}

func resourceBackupProtectionPolicyFileShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		if properties, ok := model.Properties.(protectionpolicies.AzureFileShareProtectionPolicy); ok {
			d.Set("timezone", properties.TimeZone)

			if schedule, ok := properties.SchedulePolicy.(protectionpolicies.SimpleSchedulePolicy); ok {
				if s, err := flattenBackupProtectionPolicyFileShareSchedule(schedule); err != nil {
					return fmt.Errorf("flattening `backup`: %+v", err)
				} else if err := d.Set("backup", s); err != nil {
					return fmt.Errorf("setting `backup`: %+v", err)
				}
			}

			if retention, ok := properties.RetentionPolicy.(protectionpolicies.LongTermRetentionPolicy); ok {
				if s := retention.DailySchedule; s != nil {
					if err := d.Set("retention_daily", flattenBackupProtectionPolicyFileShareRetentionDaily(s)); err != nil {
						return fmt.Errorf("setting `retention_daily`: %+v", err)
					}
				} else {
					d.Set("retention_daily", nil)
				}

				if s := retention.WeeklySchedule; s != nil {
					if err := d.Set("retention_weekly", flattenBackupProtectionPolicyFileShareRetentionWeekly(s)); err != nil {
						return fmt.Errorf("setting `retention_weekly`: %+v", err)
					}
				} else {
					d.Set("retention_weekly", nil)
				}

				if s := retention.MonthlySchedule; s != nil {
					if err := d.Set("retention_monthly", flattenBackupProtectionPolicyFileShareRetentionMonthly(s)); err != nil {
						return fmt.Errorf("setting `retention_monthly`: %+v", err)
					}
				} else {
					d.Set("retention_monthly", nil)
				}

				if s := retention.YearlySchedule; s != nil {
					if err := d.Set("retention_yearly", flattenBackupProtectionPolicyFileShareRetentionYearly(s)); err != nil {
						return fmt.Errorf("setting `retention_yearly`: %+v", err)
					}
				} else {
					d.Set("retention_yearly", nil)
				}
			}
		}
	}

	return nil
}

func resourceBackupProtectionPolicyFileShareDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

	return resourceBackupProtectionPolicyFileShareWaitForDeletion(ctx, client, *id, d)
}

func expandBackupProtectionPolicyFileShareSchedule(d *pluginsdk.ResourceData, times []string) *protectionpolicies.SimpleSchedulePolicy {
	if bb, ok := d.Get("backup").([]interface{}); ok && len(bb) > 0 {
		block := bb[0].(map[string]interface{})

		schedule := protectionpolicies.SimpleSchedulePolicy{ // LongTermSchedulePolicy has no properties
			ScheduleRunTimes: &times,
		}

		if v, ok := block["frequency"].(string); ok {
			schedule.ScheduleRunFrequency = pointer.To(protectionpolicies.ScheduleRunType(v))
		}

		if v, ok := block["hourly"].([]interface{}); ok && len(v) > 0 {
			hourlyBlock := v[0].(map[string]interface{})

			if schedule.ScheduleRunFrequency != nil && *schedule.ScheduleRunFrequency == protectionpolicies.ScheduleRunTypeHourly {
				schedule.HourlySchedule = &protectionpolicies.HourlySchedule{}
				if v, ok := hourlyBlock["interval"].(int); ok {
					schedule.HourlySchedule.Interval = pointer.To(int64(v))
				}
				if t, ok := hourlyBlock["start_time"].(string); ok {
					schedule.HourlySchedule.ScheduleWindowStartTime = pointer.To(fmt.Sprintf("2018-07-30T%s:00.000Z", t))
				}
				if v, ok := hourlyBlock["window_duration"].(int); ok {
					schedule.HourlySchedule.ScheduleWindowDuration = pointer.To(int64(v))
				}
			}
		}

		return &schedule
	}

	return nil
}

func expandBackupProtectionPolicyFileShareRetentionDaily(d *pluginsdk.ResourceData, times []string) *protectionpolicies.DailyRetentionSchedule {
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

func expandBackupProtectionPolicyFileShareRetentionWeekly(d *pluginsdk.ResourceData, times []string) *protectionpolicies.WeeklyRetentionSchedule {
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

func expandBackupProtectionPolicyFileShareRetentionMonthly(d *pluginsdk.ResourceData, times []string) *protectionpolicies.MonthlyRetentionSchedule {
	if rb, ok := d.Get("retention_monthly").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		scheduleFormat := protectionpolicies.RetentionScheduleFormatWeekly
		var weekly *protectionpolicies.WeeklyRetentionFormat = nil
		var daily *protectionpolicies.DailyRetentionFormat = nil
		if v, ok := block["days"]; ok && v.(*pluginsdk.Set).Len() > 0 || block["include_last_days"].(bool) {
			scheduleFormat = protectionpolicies.RetentionScheduleFormatDaily
			daily = expandBackupProtectionPolicyFileShareRetentionDailyFormat(block)
		} else {
			weekly = expandBackupProtectionPolicyFileShareRetentionWeeklyFormat(block)
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

func expandBackupProtectionPolicyFileShareRetentionYearly(d *pluginsdk.ResourceData, times []string) *protectionpolicies.YearlyRetentionSchedule {
	if rb, ok := d.Get("retention_yearly").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		scheduleFormat := protectionpolicies.RetentionScheduleFormatWeekly
		var weekly *protectionpolicies.WeeklyRetentionFormat = nil
		var daily *protectionpolicies.DailyRetentionFormat = nil
		if v, ok := block["days"]; ok && v.(*pluginsdk.Set).Len() > 0 || block["include_last_days"].(bool) {
			scheduleFormat = protectionpolicies.RetentionScheduleFormatDaily
			daily = expandBackupProtectionPolicyFileShareRetentionDailyFormat(block)
		} else {
			weekly = expandBackupProtectionPolicyFileShareRetentionWeeklyFormat(block)
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

func expandBackupProtectionPolicyFileShareRetentionWeeklyFormat(block map[string]interface{}) *protectionpolicies.WeeklyRetentionFormat {
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

func expandBackupProtectionPolicyFileShareRetentionDailyFormat(block map[string]interface{}) *protectionpolicies.DailyRetentionFormat {
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

func flattenBackupProtectionPolicyFileShareSchedule(schedule protectionpolicies.SimpleSchedulePolicy) ([]interface{}, error) {
	block := map[string]interface{}{}

	block["frequency"] = string(pointer.From(schedule.ScheduleRunFrequency))

	if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
		policyTime, _ := time.Parse(time.RFC3339, (*times)[0])
		block["time"] = policyTime.Format("15:04")
	}

	if hourly := schedule.HourlySchedule; hourly != nil {
		hourlyBlock := make(map[string]interface{}, 0)
		if hourly.ScheduleWindowStartTime != nil {
			startTime, err := time.Parse(time.RFC3339, *hourly.ScheduleWindowStartTime)
			if err != nil {
				return nil, fmt.Errorf("error parsing schedule window start time: %s", err)
			}
			hourlyBlock["start_time"] = startTime.Format("15:04")
		}
		hourlyBlock["interval"] = pointer.From(hourly.Interval)
		hourlyBlock["window_duration"] = pointer.From(hourly.ScheduleWindowDuration)
		block["hourly"] = []interface{}{hourlyBlock}
	}

	return []interface{}{block}, nil
}

func flattenBackupProtectionPolicyFileShareRetentionDaily(daily *protectionpolicies.DailyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := daily.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	return []interface{}{block}
}

func flattenBackupProtectionPolicyFileShareRetentionWeekly(weekly *protectionpolicies.WeeklyRetentionSchedule) []interface{} {
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

func flattenBackupProtectionPolicyFileShareRetentionMonthly(monthly *protectionpolicies.MonthlyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := monthly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if weekly := monthly.RetentionScheduleWeekly; weekly != nil {
		block["weekdays"], block["weeks"] = flattenBackupProtectionPolicyFileShareRetentionWeeklyFormat(weekly)
	}

	if daily := monthly.RetentionScheduleDaily; daily != nil {
		block["days"], block["include_last_days"] = flattenBackupProtectionPolicyFileRetentionDailyFormat(daily)
	}

	return []interface{}{block}
}

func flattenBackupProtectionPolicyFileShareRetentionYearly(yearly *protectionpolicies.YearlyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := yearly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if weekly := yearly.RetentionScheduleWeekly; weekly != nil {
		block["weekdays"], block["weeks"] = flattenBackupProtectionPolicyFileShareRetentionWeeklyFormat(weekly)
	}

	if daily := yearly.RetentionScheduleDaily; daily != nil {
		block["days"], block["include_last_days"] = flattenBackupProtectionPolicyFileRetentionDailyFormat(daily)
	}

	if months := yearly.MonthsOfYear; months != nil {
		slice := make([]interface{}, 0)
		for _, d := range *months {
			slice = append(slice, string(d))
		}
		block["months"] = pluginsdk.NewSet(pluginsdk.HashString, slice)
	}

	return []interface{}{block}
}

func flattenBackupProtectionPolicyFileShareRetentionWeeklyFormat(retention *protectionpolicies.WeeklyRetentionFormat) (weekdays, weeks *pluginsdk.Set) {
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

func flattenBackupProtectionPolicyFileRetentionDailyFormat(retention *protectionpolicies.DailyRetentionFormat) (days []interface{}, includeLastDay bool) {
	if dotm := retention.DaysOfTheMonth; dotm != nil {
		for _, d := range *dotm {
			// for the last date, the service will return a record with date == 0.
			if d.Date != nil && (*d.Date != 0) {
				days = append(days, *d.Date)
			}
			if d.IsLast != nil && *d.IsLast {
				includeLastDay = true
			}
		}
	}
	return days, includeLastDay
}

func resourceBackupProtectionPolicyFileShareWaitForUpdate(ctx context.Context, client *protectionpolicies.ProtectionPoliciesClient, id protectionpolicies.BackupPolicyId, d *pluginsdk.ResourceData) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"NotFound"},
		Target:     []string{"Found"},
		Refresh:    resourceBackupProtectionPolicyFileShareRefreshFunc(ctx, client, id),
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	_, err := state.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for update %s: %+v", id, err)
	}

	return nil
}

func resourceBackupProtectionPolicyFileShareWaitForDeletion(ctx context.Context, client *protectionpolicies.ProtectionPoliciesClient, id protectionpolicies.BackupPolicyId, d *pluginsdk.ResourceData) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"Found"},
		Target:     []string{"NotFound"},
		Refresh:    resourceBackupProtectionPolicyFileShareRefreshFunc(ctx, client, id),
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
	}

	_, err := state.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for delete to finish for %s: %+v", id, err)
	}

	return nil
}

func resourceBackupProtectionPolicyFileShareRefreshFunc(ctx context.Context, client *protectionpolicies.ProtectionPoliciesClient, id protectionpolicies.BackupPolicyId) pluginsdk.StateRefreshFunc {
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

func resourceBackupProtectionPolicyFileShareSchema() map[string]*pluginsdk.Schema {
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

		"backup": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"frequency": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc: validation.StringInSlice([]string{
							string(protectionpolicies.ScheduleRunTypeDaily),
							string(protectionpolicies.ScheduleRunTypeHourly),
						}, false),
					},

					"time": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
							"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
						),
						AtLeastOneOf: []string{
							"backup.0.time",
							"backup.0.hourly",
						},
						ConflictsWith: []string{
							"backup.0.hourly",
						},
					},

					"hourly": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						AtLeastOneOf: []string{
							"backup.0.time",
							"backup.0.hourly",
						},
						ConflictsWith: []string{
							"backup.0.time",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"interval": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntInSlice([]int{4, 6, 8, 12}),
								},

								"start_time": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
										"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
									),
								},

								"window_duration": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(4, 24),
								},
							},
						},
					},
				},
			},
		},

		"retention_daily": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 200),
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
						ValidateFunc: validation.IntBetween(1, 200),
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
						ValidateFunc: validation.IntBetween(1, 120),
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
						ValidateFunc: validation.IntBetween(1, 10),
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
