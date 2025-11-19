// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DefaultRetentionRule struct {
	LifeCycle []LifeCycle `tfschema:"life_cycle"`
}

type RetentionRule struct {
	Name      string      `tfschema:"name"`
	Criteria  []Criteria  `tfschema:"criteria"`
	Priority  int64       `tfschema:"priority"`
	LifeCycle []LifeCycle `tfschema:"life_cycle"`
}

type LifeCycle struct {
	DataStoreType string `tfschema:"data_store_type"`
	Duration      string `tfschema:"duration"`
}

type Criteria struct {
	AbsoluteCriteria     string   `tfschema:"absolute_criteria"`
	DaysOfWeek           []string `tfschema:"days_of_week"`
	DaysOfMonth          []int64  `tfschema:"days_of_month"`
	MonthsOfYear         []string `tfschema:"months_of_year"`
	ScheduledBackupTimes []string `tfschema:"scheduled_backup_times"`
	WeeksOfMonth         []string `tfschema:"weeks_of_month"`
}

func ExpandBackupPolicyAzureBackupRuleArray(input []string, timeZone string, backupType string, datastoreType backuppolicies.DataStoreTypes, taggingCriteria *[]backuppolicies.TaggingCriteria) []backuppolicies.BasePolicyRule {
	results := make([]backuppolicies.BasePolicyRule, 0)
	results = append(results, backuppolicies.AzureBackupRule{
		Name: "BackupIntervals",
		DataStore: backuppolicies.DataStoreInfoBase{
			DataStoreType: datastoreType,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: &backuppolicies.AzureBackupParams{
			BackupType: backupType,
		},
		Trigger: backuppolicies.ScheduleBasedTriggerContext{
			Schedule: backuppolicies.BackupSchedule{
				RepeatingTimeIntervals: input,
				TimeZone:               pointer.To(timeZone),
			},
			TaggingCriteria: *taggingCriteria,
		},
	})

	return results
}

func ExpandBackupPolicyDefaultRetentionRule(input []DefaultRetentionRule) *backuppolicies.AzureRetentionRule {
	if len(input) == 0 {
		return nil
	}
	return &backuppolicies.AzureRetentionRule{
		Name:       "Default",
		IsDefault:  pointer.To(true),
		Lifecycles: expandBackupPolicyLifeCycle(input[0].LifeCycle),
	}
}

func ExpandBackupPolicyDefaultRetentionRuleArray(input string, dataStoreType backuppolicies.DataStoreTypes) backuppolicies.BasePolicyRule {
	return &backuppolicies.AzureRetentionRule{
		Name:      "Default",
		IsDefault: pointer.To(true),
		Lifecycles: []backuppolicies.SourceLifeCycle{
			{
				DeleteAfter: backuppolicies.AbsoluteDeleteOption{
					Duration: input,
				},
				SourceDataStore: backuppolicies.DataStoreInfoBase{
					DataStoreType: dataStoreType,
					ObjectType:    "DataStoreInfoBase",
				},
				TargetDataStoreCopySettings: &[]backuppolicies.TargetCopySetting{},
			},
		},
	}
}

func ExpandBackupPolicyAzureRetentionRules(input []RetentionRule) []backuppolicies.BasePolicyRule {
	results := make([]backuppolicies.BasePolicyRule, 0)
	for _, item := range input {
		lifeCycle := expandBackupPolicyLifeCycle(item.LifeCycle)

		results = append(results, backuppolicies.AzureRetentionRule{
			Name:       item.Name,
			IsDefault:  pointer.To(false),
			Lifecycles: lifeCycle,
		})
	}
	return results
}

func expandBackupPolicyLifeCycle(input []LifeCycle) []backuppolicies.SourceLifeCycle {
	results := make([]backuppolicies.SourceLifeCycle, 0)
	for _, item := range input {
		sourceLifeCycle := backuppolicies.SourceLifeCycle{
			DeleteAfter: backuppolicies.AbsoluteDeleteOption{
				Duration: item.Duration,
			},
			SourceDataStore: backuppolicies.DataStoreInfoBase{
				DataStoreType: backuppolicies.DataStoreTypes(item.DataStoreType),
				ObjectType:    "DataStoreInfoBase",
			},
			TargetDataStoreCopySettings: &[]backuppolicies.TargetCopySetting{},
		}
		results = append(results, sourceLifeCycle)
	}

	return results
}

func ExpandBackupPolicyTaggingCriteriaArray(input []RetentionRule) (*[]backuppolicies.TaggingCriteria, error) {
	results := []backuppolicies.TaggingCriteria{
		{
			Criteria:        nil,
			IsDefault:       true,
			TaggingPriority: 99,
			TagInfo: backuppolicies.RetentionTag{
				Id:      utils.String("Default_"),
				TagName: "Default",
			},
		},
	}
	for _, item := range input {
		result := backuppolicies.TaggingCriteria{
			IsDefault:       false,
			TaggingPriority: item.Priority,
			TagInfo: backuppolicies.RetentionTag{
				Id:      utils.String(item.Name + "_"),
				TagName: item.Name,
			},
		}

		criteria, err := expandBackupPolicyCriteriaArray(item.Criteria)
		if err != nil {
			return nil, err
		}
		result.Criteria = criteria
		results = append(results, result)
	}
	return &results, nil
}

func expandBackupPolicyCriteriaArray(input []Criteria) (*[]backuppolicies.BackupCriteria, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("criteria is a required field, cannot leave blank")
	}

	results := make([]backuppolicies.BackupCriteria, 0)

	for _, item := range input {
		var absoluteCriteria []backuppolicies.AbsoluteMarker
		if absoluteCriteriaRaw := item.AbsoluteCriteria; len(absoluteCriteriaRaw) > 0 {
			absoluteCriteria = []backuppolicies.AbsoluteMarker{backuppolicies.AbsoluteMarker(absoluteCriteriaRaw)}
		}

		var daysOfWeek []backuppolicies.DayOfWeek
		if len(item.DaysOfWeek) > 0 {
			daysOfWeek = make([]backuppolicies.DayOfWeek, 0)
			for _, value := range item.DaysOfWeek {
				daysOfWeek = append(daysOfWeek, backuppolicies.DayOfWeek(value))
			}
		}

		var daysOfMonth []backuppolicies.Day
		if len(item.DaysOfMonth) > 0 {
			daysOfMonth = make([]backuppolicies.Day, 0)
			for _, value := range item.DaysOfMonth {
				isLast := false
				if value == 0 {
					isLast = true
				}
				daysOfMonth = append(daysOfMonth, backuppolicies.Day{
					Date:   pointer.To(value),
					IsLast: pointer.To(isLast),
				})
			}
		} else {
			// for kubernetes clusters backup policies at least this is not used and always nil
			daysOfMonth = nil
		}

		var monthsOfYear []backuppolicies.Month
		if len(item.MonthsOfYear) > 0 {
			monthsOfYear = make([]backuppolicies.Month, 0)
			for _, value := range item.MonthsOfYear {
				monthsOfYear = append(monthsOfYear, backuppolicies.Month(value))
			}
		}

		var weeksOfMonth []backuppolicies.WeekNumber
		if len(item.WeeksOfMonth) > 0 {
			weeksOfMonth = make([]backuppolicies.WeekNumber, 0)
			for _, value := range item.WeeksOfMonth {
				weeksOfMonth = append(weeksOfMonth, backuppolicies.WeekNumber(value))
			}
		}

		results = append(results, backuppolicies.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: &absoluteCriteria,
			DaysOfMonth:      &daysOfMonth,
			DaysOfTheWeek:    &daysOfWeek,
			MonthsOfYear:     &monthsOfYear,
			ScheduleTimes:    pointer.To(item.ScheduledBackupTimes),
			WeeksOfTheMonth:  &weeksOfMonth,
		})
	}
	return &results, nil
}

func FlattenBackupPolicyBackupRuleArray(input *[]backuppolicies.BasePolicyRule) []string {
	if input == nil {
		return make([]string, 0)
	}
	for _, item := range *input {
		if backupRule, ok := item.(backuppolicies.AzureBackupRule); ok {
			if backupRule.Trigger != nil {
				if scheduleBasedTrigger, ok := backupRule.Trigger.(backuppolicies.ScheduleBasedTriggerContext); ok {
					return scheduleBasedTrigger.Schedule.RepeatingTimeIntervals
				}
			}
		}
	}
	return make([]string, 0)
}

func FlattenBackupPolicyBackupTimeZone(input *[]backuppolicies.BasePolicyRule) string {
	if input == nil {
		return ""
	}
	for _, item := range *input {
		if backupRule, ok := item.(backuppolicies.AzureBackupRule); ok {
			if backupRule.Trigger != nil {
				if scheduleBasedTrigger, ok := backupRule.Trigger.(backuppolicies.ScheduleBasedTriggerContext); ok {
					return pointer.From(scheduleBasedTrigger.Schedule.TimeZone)
				}
			}
		}
	}
	return ""
}

func FlattenBackupPolicyRetentionRules(input *[]backuppolicies.BasePolicyRule) []RetentionRule {
	results := make([]RetentionRule, 0)
	if input == nil {
		return results
	}

	var taggingCriterias []backuppolicies.TaggingCriteria
	for _, item := range *input {
		if backupRule, ok := item.(backuppolicies.AzureBackupRule); ok {
			if trigger, ok := backupRule.Trigger.(backuppolicies.ScheduleBasedTriggerContext); ok {
				taggingCriterias = trigger.TaggingCriteria
			}
		}
	}

	for _, item := range *input {
		if retentionRule, ok := item.(backuppolicies.AzureRetentionRule); ok {
			var name string
			var taggingPriority int64
			var taggingCriteria []Criteria
			if retentionRule.IsDefault == nil || !*retentionRule.IsDefault {
				name = retentionRule.Name
				for _, criteria := range taggingCriterias {
					if strings.EqualFold(criteria.TagInfo.TagName, name) {
						taggingPriority = criteria.TaggingPriority
						taggingCriteria = flattenBackupPolicyBackupCriteriaArray(criteria.Criteria)
						break
					}
				}

				var lifeCycle []LifeCycle
				if v := retentionRule.Lifecycles; len(v) > 0 {
					lifeCycle = flattenBackupPolicyBackupLifeCycleArray(v)
				}
				results = append(results, RetentionRule{
					Name:      name,
					Priority:  taggingPriority,
					Criteria:  taggingCriteria,
					LifeCycle: lifeCycle,
				})
			}
		}
	}
	return results
}

// This function is currently unused, but can be used by either of data_protection_backup_policy_mysql_flexible_server or azurerm_data_protection_backup_policy_kubernetes_cluster at least
// func FlattenBackupPolicyDefaultRetentionRule(input *[]backuppolicies.BasePolicyRule) []DefaultRetentionRule {
// 	results := make([]DefaultRetentionRule, 0)
// 	if input == nil {
// 		return results
// 	}

// 	for _, item := range *input {
// 		if retentionRule, ok := item.(backuppolicies.AzureRetentionRule); ok {
// 			if pointer.From(retentionRule.IsDefault) {
// 				var lifeCycle []LifeCycle
// 				if v := retentionRule.Lifecycles; len(v) > 0 {
// 					lifeCycle = flattenBackupPolicyBackupLifeCycleArray(v)
// 				}

// 				results = append(results, DefaultRetentionRule{
// 					LifeCycle: lifeCycle,
// 				})
// 			}
// 		}
// 	}
// 	return results
// }

func FlattenBackupPolicyDefaultRetentionRuleDuration(input *[]backuppolicies.BasePolicyRule, dsType backuppolicies.DataStoreTypes) string {
	if input == nil {
		return ""
	}

	for _, item := range *input {
		if retentionRule, ok := item.(backuppolicies.AzureRetentionRule); ok && retentionRule.IsDefault != nil && *retentionRule.IsDefault {
			if len(retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := (retentionRule.Lifecycles)[0].DeleteAfter.(backuppolicies.AbsoluteDeleteOption); ok {
					if (retentionRule.Lifecycles)[0].SourceDataStore.DataStoreType == dsType {
						return deleteOption.Duration
					}
				}
			}
		}
	}
	return ""
}

func flattenBackupPolicyBackupCriteriaArray(input *[]backuppolicies.BackupCriteria) []Criteria {
	results := make([]Criteria, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if criteria, ok := item.(backuppolicies.ScheduleBasedBackupCriteria); ok {
			var absoluteCriteria string
			if criteria.AbsoluteCriteria != nil && len(*criteria.AbsoluteCriteria) > 0 {
				absoluteCriteria = string((*criteria.AbsoluteCriteria)[0])
			}
			var daysOfWeek []string
			if criteria.DaysOfTheWeek != nil {
				daysOfWeek = make([]string, 0)
				for _, item := range *criteria.DaysOfTheWeek {
					daysOfWeek = append(daysOfWeek, (string)(item))
				}
			}
			var daysOfMonth []int64
			if criteria.DaysOfMonth != nil {
				daysOfMonth = make([]int64, 0)
				for _, item := range *criteria.DaysOfMonth {
					daysOfMonth = append(daysOfMonth, *item.Date)
				}
			}
			var monthsOfYear []string
			if criteria.MonthsOfYear != nil {
				monthsOfYear = make([]string, 0)
				for _, item := range *criteria.MonthsOfYear {
					monthsOfYear = append(monthsOfYear, (string)(item))
				}
			}
			var weeksOfMonth []string
			if criteria.WeeksOfTheMonth != nil {
				weeksOfMonth = make([]string, 0)
				for _, item := range *criteria.WeeksOfTheMonth {
					weeksOfMonth = append(weeksOfMonth, (string)(item))
				}
			}
			var scheduleTimes []string
			if criteria.ScheduleTimes != nil {
				scheduleTimes = make([]string, 0)
				scheduleTimes = append(scheduleTimes, *criteria.ScheduleTimes...)
			}

			results = append(results, Criteria{
				AbsoluteCriteria:     absoluteCriteria,
				DaysOfWeek:           daysOfWeek,
				DaysOfMonth:          daysOfMonth,
				MonthsOfYear:         monthsOfYear,
				WeeksOfMonth:         weeksOfMonth,
				ScheduledBackupTimes: scheduleTimes,
			})
		}
	}
	return results
}

func flattenBackupPolicyBackupLifeCycleArray(input []backuppolicies.SourceLifeCycle) []LifeCycle {
	results := make([]LifeCycle, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		var duration string
		var dataStoreType string
		if deleteOption, ok := item.DeleteAfter.(backuppolicies.AbsoluteDeleteOption); ok {
			duration = deleteOption.Duration
		}
		dataStoreType = string(item.SourceDataStore.DataStoreType)

		results = append(results, LifeCycle{
			Duration:      duration,
			DataStoreType: dataStoreType,
		})
	}
	return results
}
