package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
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
			_, err := parse.BackupPolicyID(id)
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
			case string(backup.ScheduleRunTypeHourly):
				if !hasDaily {
					return fmt.Errorf("`retention_daily` must be set when backup.0.frequency is hourly")
				}

				if _, ok := diff.GetOk("backup.0.weekdays"); ok {
					return fmt.Errorf("`backup.0.weekdays` should be not set when backup.0.frequency is hourly")
				}
			case string(backup.ScheduleRunTypeDaily):
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
			case string(backup.ScheduleRunTypeWeekly):
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
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	policyName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)

	log.Printf("[DEBUG] Creating/updating Azure Backup Protection Policy %s (resource group %q)", policyName, resourceGroup)

	// getting this ready now because its shared between *everything*, time is... complicated for this resource
	timeOfDay := d.Get("backup.0.time").(string)
	dateOfDay, err := time.Parse(time.RFC3339, fmt.Sprintf("2018-07-30T%s:00Z", timeOfDay))
	if err != nil {
		return fmt.Errorf("generating time from %q for policy %q (Resource Group %q): %+v", timeOfDay, policyName, resourceGroup, err)
	}
	times := append(make([]date.Time, 0), date.Time{Time: dateOfDay})

	if d.IsNewResource() {
		existing, err2 := client.Get(ctx, vaultName, resourceGroup, policyName)
		if err2 != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Azure Backup Protection Policy %q (Resource Group %q): %+v", policyName, resourceGroup, err2)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_backup_policy_vm", *existing.ID)
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

	policyType := backup.IAASVMPolicyType(d.Get("policy_type").(string))
	vmProtectionPolicyProperties := &backup.AzureIaaSVMProtectionPolicy{
		TimeZone:             utils.String(d.Get("timezone").(string)),
		BackupManagementType: backup.ManagementTypeBasicProtectionPolicyBackupManagementTypeAzureIaasVM,
		PolicyType:           policyType,
		SchedulePolicy:       *schedulePolicy,
		RetentionPolicy: &backup.LongTermRetentionPolicy{ // SimpleRetentionPolicy only has duration property ¯\_(ツ)_/¯
			RetentionPolicyType: backup.RetentionPolicyTypeLongTermRetentionPolicy,
			DailySchedule:       expandBackupProtectionPolicyVMRetentionDaily(d, times),
			WeeklySchedule:      expandBackupProtectionPolicyVMRetentionWeekly(d, times),
			MonthlySchedule:     expandBackupProtectionPolicyVMRetentionMonthly(d, times),
			YearlySchedule:      expandBackupProtectionPolicyVMRetentionYearly(d, times),
		},
	}

	if d.HasChange("instant_restore_retention_days") {
		days := d.Get("instant_restore_retention_days").(int)
		if backup.IAASVMPolicyTypeV1 == policyType && days > 5 {
			return fmt.Errorf("`instant_restore_retention_days` must be less than or equal to `5` when `policy_type` is `V1`")
		}

		vmProtectionPolicyProperties.InstantRpRetentionRangeInDays = utils.Int32(int32(days))
	}

	policy := backup.ProtectionPolicyResource{
		Properties: vmProtectionPolicyProperties,
	}

	if _, err = client.CreateOrUpdate(ctx, vaultName, resourceGroup, policyName, policy); err != nil {
		return fmt.Errorf("creating/updating Azure Backup Protection Policy %q (Resource Group %q): %+v", policyName, resourceGroup, err)
	}

	resp, err := resourceBackupProtectionPolicyVMWaitForUpdate(ctx, client, vaultName, resourceGroup, policyName, d)
	if err != nil {
		return err
	}

	id := strings.Replace(*resp.ID, "Subscriptions", "subscriptions", 1)
	d.SetId(id)

	return resourceBackupProtectionPolicyVMRead(d, meta)
}

func resourceBackupProtectionPolicyVMRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading Azure Backup Protection Policy %q (resource group %q)", id.Name, id.ResourceGroup)

	resp, err := client.Get(ctx, id.VaultName, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Azure Backup Protection Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("recovery_vault_name", id.VaultName)

	if properties, ok := resp.Properties.AsAzureIaaSVMProtectionPolicy(); ok && properties != nil {
		d.Set("timezone", properties.TimeZone)
		d.Set("instant_restore_retention_days", properties.InstantRpRetentionRangeInDays)

		if schedule, ok := properties.SchedulePolicy.AsSimpleSchedulePolicy(); ok && schedule != nil {
			if err := d.Set("backup", flattenBackupProtectionPolicyVMSchedule(schedule)); err != nil {
				return fmt.Errorf("setting `backup`: %+v", err)
			}
		}

		if schedule, ok := properties.SchedulePolicy.AsSimpleSchedulePolicyV2(); ok && schedule != nil {
			if err := d.Set("backup", flattenBackupProtectionPolicyVMScheduleV2(schedule)); err != nil {
				return fmt.Errorf("setting `backup`: %+v", err)
			}
		}

		policyType := string(backup.IAASVMPolicyTypeV1)
		if properties.PolicyType != "" {
			policyType = string(properties.PolicyType)
		}
		d.Set("policy_type", policyType)

		if retention, ok := properties.RetentionPolicy.AsLongTermRetentionPolicy(); ok && retention != nil {
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
	}

	return nil
}

func resourceBackupProtectionPolicyVMDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Azure Backup Protected Item %q (resource group %q)", id.Name, id.ResourceGroup)

	future, err := client.Delete(ctx, id.VaultName, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	resp, err := future.Result(*client)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("issuing delete request for Azure Backup Protection Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	if _, err := resourceBackupProtectionPolicyVMWaitForDeletion(ctx, client, id.VaultName, id.ResourceGroup, id.Name, d); err != nil {
		return err
	}

	return nil
}

func expandBackupProtectionPolicyVMSchedule(d *pluginsdk.ResourceData, times []date.Time) (*backup.BasicSchedulePolicy, error) {
	if bb, ok := d.Get("backup").([]interface{}); ok && len(bb) > 0 {
		block := bb[0].(map[string]interface{})

		policyType := d.Get("policy_type").(string)
		if policyType == string(backup.IAASVMPolicyTypeV1) {
			schedule := backup.SimpleSchedulePolicy{ // LongTermSchedulePolicy has no properties
				SchedulePolicyType: backup.SchedulePolicyTypeSimpleSchedulePolicy,
				ScheduleRunTimes:   &times,
			}

			if v, ok := block["frequency"].(string); ok {
				schedule.ScheduleRunFrequency = backup.ScheduleRunType(v)
			}

			if v, ok := block["weekdays"].(*pluginsdk.Set); ok {
				days := make([]backup.DayOfWeek, 0)
				for _, day := range v.List() {
					days = append(days, backup.DayOfWeek(day.(string)))
				}
				schedule.ScheduleRunDays = &days
			}

			result, _ := schedule.AsBasicSchedulePolicy()
			return &result, nil
		} else {
			frequency := block["frequency"].(string)
			schedule := backup.SimpleSchedulePolicyV2{
				SchedulePolicyType:   backup.SchedulePolicyTypeSimpleSchedulePolicyV2,
				ScheduleRunFrequency: backup.ScheduleRunType(frequency),
			}

			switch frequency {
			case string(backup.ScheduleRunTypeHourly):
				interval, ok := block["hour_interval"].(int)
				if !ok {
					return nil, fmt.Errorf("`hour_interval` must be specified when `backup.0.frequency` is `Hourly`")
				}

				duration, ok := block["hour_duration"].(int)
				if !ok {
					return nil, fmt.Errorf("`hour_duration` must be specified when `backup.0.frequency` is `Hourly`")
				}

				if duration%interval != 0 {
					return nil, fmt.Errorf("`hour_duration` must be multiplier of `hour_interval`")
				}

				schedule.HourlySchedule = &backup.HourlySchedule{
					Interval:                utils.Int32(int32(interval)),
					ScheduleWindowStartTime: &times[0],
					ScheduleWindowDuration:  utils.Int32(int32(duration)),
				}
			case string(backup.ScheduleRunTypeDaily):
				schedule.DailySchedule = &backup.DailySchedule{
					ScheduleRunTimes: &times,
				}
			case string(backup.ScheduleRunTypeWeekly):
				weekDays, ok := block["weekdays"].(*pluginsdk.Set)
				if !ok {
					return nil, fmt.Errorf("`weekdays` must be specified when `backup.0.frequency` is `Weekly`")
				}

				days := make([]backup.DayOfWeek, 0)
				for _, day := range weekDays.List() {
					days = append(days, backup.DayOfWeek(day.(string)))
				}

				schedule.WeeklySchedule = &backup.WeeklySchedule{
					ScheduleRunDays:  &days,
					ScheduleRunTimes: &times,
				}
			default:
				return nil, fmt.Errorf("Unrecognized value for backup.0.frequency")
			}

			result, _ := schedule.AsBasicSchedulePolicy()
			return &result, nil
		}
	}

	return nil, nil
}

func expandBackupProtectionPolicyVMRetentionDaily(d *pluginsdk.ResourceData, times []date.Time) *backup.DailyRetentionSchedule {
	if rb, ok := d.Get("retention_daily").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		return &backup.DailyRetentionSchedule{
			RetentionTimes: &times,
			RetentionDuration: &backup.RetentionDuration{
				Count:        utils.Int32(int32(block["count"].(int))),
				DurationType: backup.RetentionDurationTypeDays,
			},
		}
	}

	return nil
}

func expandBackupProtectionPolicyVMRetentionWeekly(d *pluginsdk.ResourceData, times []date.Time) *backup.WeeklyRetentionSchedule {
	if rb, ok := d.Get("retention_weekly").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		retention := backup.WeeklyRetentionSchedule{
			RetentionTimes: &times,
			RetentionDuration: &backup.RetentionDuration{
				Count:        utils.Int32(int32(block["count"].(int))),
				DurationType: backup.RetentionDurationTypeWeeks,
			},
		}

		if v, ok := block["weekdays"].(*pluginsdk.Set); ok {
			days := make([]backup.DayOfWeek, 0)
			for _, day := range v.List() {
				days = append(days, backup.DayOfWeek(day.(string)))
			}
			retention.DaysOfTheWeek = &days
		}

		return &retention
	}

	return nil
}

func expandBackupProtectionPolicyVMRetentionMonthly(d *pluginsdk.ResourceData, times []date.Time) *backup.MonthlyRetentionSchedule {
	if rb, ok := d.Get("retention_monthly").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		retention := backup.MonthlyRetentionSchedule{
			RetentionScheduleFormatType: backup.RetentionScheduleFormatWeekly, // this is always weekly ¯\_(ツ)_/¯
			RetentionScheduleDaily:      nil,                                  // and this is always nil..
			RetentionScheduleWeekly:     expandBackupProtectionPolicyVMRetentionWeeklyFormat(block),
			RetentionTimes:              &times,
			RetentionDuration: &backup.RetentionDuration{
				Count:        utils.Int32(int32(block["count"].(int))),
				DurationType: backup.RetentionDurationTypeMonths,
			},
		}

		return &retention
	}

	return nil
}

func expandBackupProtectionPolicyVMRetentionYearly(d *pluginsdk.ResourceData, times []date.Time) *backup.YearlyRetentionSchedule {
	if rb, ok := d.Get("retention_yearly").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		retention := backup.YearlyRetentionSchedule{
			RetentionScheduleFormatType: backup.RetentionScheduleFormatWeekly, // this is always weekly ¯\_(ツ)_/¯
			RetentionScheduleDaily:      nil,                                  // and this is always nil..
			RetentionScheduleWeekly:     expandBackupProtectionPolicyVMRetentionWeeklyFormat(block),
			RetentionTimes:              &times,
			RetentionDuration: &backup.RetentionDuration{
				Count:        utils.Int32(int32(block["count"].(int))),
				DurationType: backup.RetentionDurationTypeYears,
			},
		}

		if v, ok := block["months"].(*pluginsdk.Set); ok {
			months := make([]backup.MonthOfYear, 0)
			for _, month := range v.List() {
				months = append(months, backup.MonthOfYear(month.(string)))
			}
			retention.MonthsOfYear = &months
		}

		return &retention
	}

	return nil
}

func expandBackupProtectionPolicyVMRetentionWeeklyFormat(block map[string]interface{}) *backup.WeeklyRetentionFormat {
	weekly := backup.WeeklyRetentionFormat{}

	if v, ok := block["weekdays"].(*pluginsdk.Set); ok {
		days := make([]backup.DayOfWeek, 0)
		for _, day := range v.List() {
			days = append(days, backup.DayOfWeek(day.(string)))
		}
		weekly.DaysOfTheWeek = &days
	}

	if v, ok := block["weeks"].(*pluginsdk.Set); ok {
		weeks := make([]backup.WeekOfMonth, 0)
		for _, week := range v.List() {
			weeks = append(weeks, backup.WeekOfMonth(week.(string)))
		}
		weekly.WeeksOfTheMonth = &weeks
	}

	return &weekly
}

func flattenBackupProtectionPolicyVMSchedule(schedule *backup.SimpleSchedulePolicy) []interface{} {
	block := map[string]interface{}{}

	block["frequency"] = string(schedule.ScheduleRunFrequency)

	if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
		block["time"] = (*times)[0].Format("15:04")
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

func flattenBackupProtectionPolicyVMScheduleV2(schedule *backup.SimpleSchedulePolicyV2) []interface{} {
	block := map[string]interface{}{}

	frequency := schedule.ScheduleRunFrequency
	block["frequency"] = string(frequency)

	switch frequency {
	case backup.ScheduleRunTypeHourly:
		schedule := schedule.HourlySchedule
		if schedule.Interval != nil {
			block["hour_interval"] = *schedule.Interval
		}

		if schedule.ScheduleWindowDuration != nil {
			block["hour_duration"] = *schedule.ScheduleWindowDuration
		}

		if schedule.ScheduleWindowStartTime != nil {
			block["time"] = schedule.ScheduleWindowStartTime.Format("15:04")
		}
	case backup.ScheduleRunTypeDaily:
		schedule := schedule.DailySchedule
		if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
			block["time"] = (*times)[0].Format("15:04")
		}
	case backup.ScheduleRunTypeWeekly:
		schedule := schedule.WeeklySchedule
		if days := schedule.ScheduleRunDays; days != nil {
			weekdays := make([]interface{}, 0)
			for _, d := range *days {
				weekdays = append(weekdays, string(d))
			}
			block["weekdays"] = pluginsdk.NewSet(pluginsdk.HashString, weekdays)
		}

		if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
			block["time"] = (*times)[0].Format("15:04")
		}
	default:
	}

	return []interface{}{block}
}

func flattenBackupProtectionPolicyVMRetentionDaily(daily *backup.DailyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := daily.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	return []interface{}{block}
}

func flattenBackupProtectionPolicyVMRetentionWeekly(weekly *backup.WeeklyRetentionSchedule) []interface{} {
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

func flattenBackupProtectionPolicyVMRetentionMonthly(monthly *backup.MonthlyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := monthly.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	if weekly := monthly.RetentionScheduleWeekly; weekly != nil {
		block["weekdays"], block["weeks"] = flattenBackupProtectionPolicyVMRetentionWeeklyFormat(weekly)
	}

	return []interface{}{block}
}

func flattenBackupProtectionPolicyVMRetentionYearly(yearly *backup.YearlyRetentionSchedule) []interface{} {
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

	return []interface{}{block}
}

func flattenBackupProtectionPolicyVMRetentionWeeklyFormat(retention *backup.WeeklyRetentionFormat) (weekdays, weeks *pluginsdk.Set) {
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

func resourceBackupProtectionPolicyVMWaitForUpdate(ctx context.Context, client *backup.ProtectionPoliciesClient, vaultName, resourceGroup, policyName string, d *pluginsdk.ResourceData) (backup.ProtectionPolicyResource, error) {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"NotFound"},
		Target:     []string{"Found"},
		Refresh:    resourceBackupProtectionPolicyVMRefreshFunc(ctx, client, vaultName, resourceGroup, policyName),
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	resp, err := state.WaitForStateContext(ctx)
	if err != nil {
		return resp.(backup.ProtectionPolicyResource), fmt.Errorf("waiting for the Azure Backup Protection Policy %q to be true (Resource Group %q) to provision: %+v", policyName, resourceGroup, err)
	}

	return resp.(backup.ProtectionPolicyResource), nil
}

func resourceBackupProtectionPolicyVMWaitForDeletion(ctx context.Context, client *backup.ProtectionPoliciesClient, vaultName, resourceGroup, policyName string, d *pluginsdk.ResourceData) (backup.ProtectionPolicyResource, error) {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"Found"},
		Target:     []string{"NotFound"},
		Refresh:    resourceBackupProtectionPolicyVMRefreshFunc(ctx, client, vaultName, resourceGroup, policyName),
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
	}

	resp, err := state.WaitForStateContext(ctx)
	if err != nil {
		return resp.(backup.ProtectionPolicyResource), fmt.Errorf("waiting for the Azure Backup Protection Policy %q to be false (Resource Group %q) to provision: %+v", policyName, resourceGroup, err)
	}

	return resp.(backup.ProtectionPolicyResource), nil
}

func resourceBackupProtectionPolicyVMRefreshFunc(ctx context.Context, client *backup.ProtectionPoliciesClient, vaultName, resourceGroup, policyName string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, vaultName, resourceGroup, policyName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return resp, "NotFound", nil
			}

			return resp, "Error", fmt.Errorf("making Read request on Azure Backup Protection Policy %q (Resource Group %q): %+v", policyName, resourceGroup, err)
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

		"resource_group_name": azure.SchemaResourceGroupName(),

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
							string(backup.ScheduleRunTypeHourly),
							string(backup.ScheduleRunTypeDaily),
							string(backup.ScheduleRunTypeWeekly),
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
			Default:  string(backup.IAASVMPolicyTypeV1),
			ValidateFunc: validation.StringInSlice([]string{
				string(backup.IAASVMPolicyTypeV1),
				string(backup.IAASVMPolicyTypeV2),
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
								string(backup.WeekOfMonthFirst),
								string(backup.WeekOfMonthSecond),
								string(backup.WeekOfMonthThird),
								string(backup.WeekOfMonthFourth),
								string(backup.WeekOfMonthLast),
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
								string(backup.WeekOfMonthFirst),
								string(backup.WeekOfMonthSecond),
								string(backup.WeekOfMonthThird),
								string(backup.WeekOfMonthFourth),
								string(backup.WeekOfMonthLast),
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
