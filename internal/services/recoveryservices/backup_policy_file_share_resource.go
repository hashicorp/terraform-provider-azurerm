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
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/protectionpolicies"
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
			return tf.ImportAsExistsError("azurerm_backup_policy_file_share", id.ID())
		}
	}

	var retentionPolicy protectionpolicies.RetentionPolicy = protectionpolicies.LongTermRetentionPolicy{
		DailySchedule:   expandBackupProtectionPolicyFileShareRetentionDaily(d, times),
		WeeklySchedule:  expandBackupProtectionPolicyFileShareRetentionWeekly(d, times),
		MonthlySchedule: expandBackupProtectionPolicyFileShareRetentionMonthly(d, times),
		YearlySchedule:  expandBackupProtectionPolicyFileShareRetentionYearly(d, times),
	}
	var scheduledPolicy protectionpolicies.SchedulePolicy = expandBackupProtectionPolicyFileShareSchedule(d, times)
	var properties protectionpolicies.ProtectionPolicy = protectionpolicies.AzureFileShareProtectionPolicy{
		TimeZone:        pointer.To(d.Get("timezone").(string)),
		WorkLoadType:    pointer.To(protectionpolicies.WorkloadTypeAzureFileShare),
		SchedulePolicy:  pointer.To(scheduledPolicy),
		RetentionPolicy: pointer.To(retentionPolicy),
	}
	policy := protectionpolicies.ProtectionPolicyResource{
		Properties: pointer.To(properties),
	}

	if _, err = client.CreateOrUpdate(ctx, id, policy); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = resourceBackupProtectionPolicyFileShareWaitForUpdate(ctx, client, id, d); err != nil {
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
		if props := model.Properties; props != nil {
			if properties, ok := (*props).(protectionpolicies.AzureFileShareProtectionPolicy); ok {
				d.Set("timezone", properties.TimeZone)

				if err := d.Set("backup", flattenBackupProtectionPolicyFileShareSchedule(properties.SchedulePolicy)); err != nil {
					return fmt.Errorf("setting `backup`: %+v", err)
				}

				if err := d.Set("retention_daily", flattenBackupProtectionPolicyFileShareRetentionDaily(properties.RetentionPolicy)); err != nil {
					return fmt.Errorf("setting `retention_daily`: %+v", err)
				}
				if err := d.Set("retention_monthly", flattenBackupProtectionPolicyFileShareRetentionMonthly(properties.RetentionPolicy)); err != nil {
					return fmt.Errorf("setting `retention_monthly`: %+v", err)
				}
				if err := d.Set("retention_weekly", flattenBackupProtectionPolicyFileShareRetentionWeekly(properties.RetentionPolicy)); err != nil {
					return fmt.Errorf("setting `retention_weekly`: %+v", err)
				}
				if err := d.Set("retention_yearly", flattenBackupProtectionPolicyFileShareRetentionYearly(properties.RetentionPolicy)); err != nil {
					return fmt.Errorf("setting `retention_yearly`: %+v", err)
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

		retention := protectionpolicies.MonthlyRetentionSchedule{
			RetentionScheduleFormatType: pointer.To(protectionpolicies.RetentionScheduleFormatWeekly),
			RetentionScheduleDaily:      nil, // and this is always nil..
			RetentionScheduleWeekly:     expandBackupProtectionPolicyFileShareRetentionWeeklyFormat(block),
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

		retention := protectionpolicies.YearlyRetentionSchedule{
			RetentionScheduleFormatType: pointer.To(protectionpolicies.RetentionScheduleFormatWeekly), // this is always weekly ¯\_(ツ)_/¯
			RetentionScheduleDaily:      nil,                                                          // and this is always nil..
			RetentionScheduleWeekly:     expandBackupProtectionPolicyFileShareRetentionWeeklyFormat(block),
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

func flattenBackupProtectionPolicyFileShareSchedule(input *protectionpolicies.SchedulePolicy) []interface{} {
	output := make([]interface{}, 0)
	if input != nil {
		if v, ok := (*input).(protectionpolicies.SimpleSchedulePolicy); ok {
			frequency := ""
			if v.ScheduleRunFrequency != nil {
				frequency = string(pointer.From(v.ScheduleRunFrequency))
			}

			timeValue := ""
			if times := v.ScheduleRunTimes; times != nil && len(*times) > 0 {
				policyTime, err := time.Parse(time.RFC3339, (*times)[0])
				if err == nil {
					timeValue = policyTime.Format("15:04")
				}
			}

			output = append(output, map[string]interface{}{
				"frequency": frequency,
				"time":      timeValue,
			})
		}
	}

	return output
}

func flattenBackupProtectionPolicyFileShareRetentionDaily(input *protectionpolicies.RetentionPolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if retention, ok := (*input).(protectionpolicies.LongTermRetentionPolicy); ok {
			if dailySchedule := retention.DailySchedule; dailySchedule != nil {
				count := 0
				if duration := dailySchedule.RetentionDuration; duration != nil {
					if v := duration.Count; v != nil {
						count = int(*v)
					}
				}
				output = append(output, map[string]interface{}{
					"count": count,
				})
			}
		}
	}
	return output
}

func flattenBackupProtectionPolicyFileShareRetentionWeekly(input *protectionpolicies.RetentionPolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if retention, ok := (*input).(protectionpolicies.LongTermRetentionPolicy); ok {
			if weeklySchedule := retention.WeeklySchedule; weeklySchedule != nil {
				count := 0
				if duration := weeklySchedule.RetentionDuration; duration != nil {
					if duration.Count != nil {
						count = int(*duration.Count)
					}
				}

				weekdays := make([]interface{}, 0)
				if weeklySchedule.DaysOfTheWeek != nil {
					for _, d := range *weeklySchedule.DaysOfTheWeek {
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

func flattenBackupProtectionPolicyFileShareRetentionMonthly(input *protectionpolicies.RetentionPolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if retention, ok := (*input).(protectionpolicies.LongTermRetentionPolicy); ok {
			if monthlySchedule := retention.MonthlySchedule; monthlySchedule != nil {
				count := 0
				if duration := monthlySchedule.RetentionDuration; duration != nil {
					if duration.Count != nil {
						count = int(*duration.Count)
					}
				}

				weekdays := make([]interface{}, 0)
				weeks := make([]interface{}, 0)
				if weekly := monthlySchedule.RetentionScheduleWeekly; weekly != nil {
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

func flattenBackupProtectionPolicyFileShareRetentionYearly(input *protectionpolicies.RetentionPolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if retention, ok := (*input).(protectionpolicies.LongTermRetentionPolicy); ok {
			if yearlySchedule := retention.YearlySchedule; yearlySchedule != nil {
				count := 0
				if yearlySchedule.RetentionDuration != nil && yearlySchedule.RetentionDuration.Count != nil {
					count = int(*yearlySchedule.RetentionDuration.Count)
				}

				weekdays := make([]interface{}, 0)
				weeks := make([]interface{}, 0)
				months := make([]interface{}, 0)
				if weekly := yearlySchedule.RetentionScheduleWeekly; weekly != nil {
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

					if moty := yearlySchedule.MonthsOfYear; moty != nil {
						for _, d := range *moty {
							months = append(months, string(d))
						}
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
						}, false),
					},

					"time": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
							"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
						),
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
