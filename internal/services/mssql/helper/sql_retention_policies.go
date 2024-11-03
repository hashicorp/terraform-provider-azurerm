// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/backupshorttermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/longtermretentionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func LongTermRetentionPolicySchema() *pluginsdk.Schema {
	atLeastOneOf := []string{
		"long_term_retention_policy.0.weekly_retention", "long_term_retention_policy.0.monthly_retention",
		"long_term_retention_policy.0.yearly_retention", "long_term_retention_policy.0.week_of_year",
	}
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// WeeklyRetention - The weekly retention policy for an LTR backup in an ISO 8601 format.
				"weekly_retention": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validate.ISO8601Duration,
					AtLeastOneOf: atLeastOneOf,
				},

				// MonthlyRetention - The monthly retention policy for an LTR backup in an ISO 8601 format.
				"monthly_retention": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validate.ISO8601Duration,
					AtLeastOneOf: atLeastOneOf,
				},

				// YearlyRetention - The yearly retention policy for an LTR backup in an ISO 8601 format.
				"yearly_retention": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validate.ISO8601Duration,
					AtLeastOneOf: atLeastOneOf,
				},

				// WeekOfYear - The week of year to take the yearly backup in an ISO 8601 format.
				"week_of_year": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 52),
					AtLeastOneOf: atLeastOneOf,
				},
			},
		},
	}
}

func ShortTermRetentionPolicySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"retention_days": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 35),
				},
				"backup_interval_in_hours": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntInSlice([]int{12, 24}),
				},
			},
		},
	}
}

func ExpandLongTermRetentionPolicy(input []interface{}) *longtermretentionpolicies.BaseLongTermRetentionPolicyProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	policy := input[0].(map[string]interface{})

	output := longtermretentionpolicies.BaseLongTermRetentionPolicyProperties{
		WeeklyRetention:  pointer.To("PT0S"),
		MonthlyRetention: pointer.To("PT0S"),
		YearlyRetention:  pointer.To("PT0S"),
		WeekOfYear:       pointer.To(int64(1)),
	}

	if v, ok := policy["weekly_retention"].(string); ok && v != "" {
		output.WeeklyRetention = pointer.To(v)
	}

	if v, ok := policy["monthly_retention"].(string); ok && v != "" {
		output.MonthlyRetention = pointer.To(v)
	}

	if v, ok := policy["yearly_retention"].(string); ok && v != "" {
		output.YearlyRetention = pointer.To(v)
	}

	if v, ok := policy["week_of_year"].(int); ok && v != 0 {
		output.WeekOfYear = pointer.To(int64(v))
	}
	return pointer.To(output)
}

func FlattenLongTermRetentionPolicy(input *longtermretentionpolicies.LongTermRetentionPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	monthlyRetention := "PT0S"
	if input.Properties.MonthlyRetention != nil {
		monthlyRetention = pointer.From(input.Properties.MonthlyRetention)
	}

	weeklyRetention := "PT0S"
	if input.Properties.WeeklyRetention != nil {
		weeklyRetention = pointer.From(input.Properties.WeeklyRetention)
	}

	weekOfYear := int64(1)
	if input.Properties.WeekOfYear != nil && pointer.From(input.Properties.WeekOfYear) != 0 {
		weekOfYear = pointer.From(input.Properties.WeekOfYear)
	}

	yearlyRetention := "PT0S"
	if input.Properties.YearlyRetention != nil {
		yearlyRetention = *input.Properties.YearlyRetention
	}

	return []interface{}{
		map[string]interface{}{
			"monthly_retention": monthlyRetention,
			"weekly_retention":  weeklyRetention,
			"week_of_year":      weekOfYear,
			"yearly_retention":  yearlyRetention,
		},
	}
}

func ExpandShortTermRetentionPolicy(input []interface{}) *backupshorttermretentionpolicies.BackupShortTermRetentionPolicyProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	policy := input[0].(map[string]interface{})

	props := backupshorttermretentionpolicies.BackupShortTermRetentionPolicyProperties{
		RetentionDays: pointer.To(int64(7)),
	}

	if v, ok := policy["retention_days"]; ok {
		props.RetentionDays = pointer.To(int64(v.(int)))
	}

	if v, ok := policy["backup_interval_in_hours"]; ok {
		props.DiffBackupIntervalInHours = pointer.To(backupshorttermretentionpolicies.DiffBackupIntervalInHours(int64(v.(int))))
	}

	return &props
}

func FlattenShortTermRetentionPolicy(input *backupshorttermretentionpolicies.BackupShortTermRetentionPolicy) []interface{} {
	result := make([]interface{}, 0)

	if input == nil {
		return result
	}

	output := map[string]interface{}{}

	output["retention_days"] = int64(7)
	if input.Properties.RetentionDays != nil {
		output["retention_days"] = pointer.From(input.Properties.RetentionDays)
	}

	if input.Properties.DiffBackupIntervalInHours != nil {
		output["backup_interval_in_hours"] = pointer.From(input.Properties.DiffBackupIntervalInHours)
	}

	result = append(result, output)

	return result
}
