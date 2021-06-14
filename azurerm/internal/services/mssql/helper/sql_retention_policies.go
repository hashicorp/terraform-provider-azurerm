package helper

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func LongTermRetentionPolicySchema() *pluginsdk.Schema {
	atLeastOneOf := []string{"long_term_retention_policy.0.weekly_retention", "long_term_retention_policy.0.monthly_retention",
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
					ValidateFunc: validation.IntBetween(1, 52),
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
					ValidateFunc: validation.IntBetween(7, 35),
				},
			},
		},
	}
}

func ExpandLongTermRetentionPolicy(input []interface{}) *sql.LongTermRetentionPolicyProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	longTermRetentionPolicy := input[0].(map[string]interface{})

	longTermPolicyProperties := sql.LongTermRetentionPolicyProperties{
		WeeklyRetention:  utils.String("PT0S"),
		MonthlyRetention: utils.String("PT0S"),
		YearlyRetention:  utils.String("PT0S"),
		WeekOfYear:       utils.Int32(1),
	}

	if v, ok := longTermRetentionPolicy["weekly_retention"]; ok {
		longTermPolicyProperties.WeeklyRetention = utils.String(v.(string))
	}

	if v, ok := longTermRetentionPolicy["monthly_retention"]; ok {
		longTermPolicyProperties.MonthlyRetention = utils.String(v.(string))
	}

	if v, ok := longTermRetentionPolicy["yearly_retention"]; ok {
		longTermPolicyProperties.YearlyRetention = utils.String(v.(string))
	}

	if v, ok := longTermRetentionPolicy["week_of_year"]; ok {
		longTermPolicyProperties.WeekOfYear = utils.Int32(int32(v.(int)))
	}

	return &longTermPolicyProperties
}

func FlattenLongTermRetentionPolicy(longTermRetentionPolicy *sql.BackupLongTermRetentionPolicy, d *pluginsdk.ResourceData) []interface{} {
	if longTermRetentionPolicy == nil {
		return []interface{}{}
	}

	monthlyRetention := "PT0S"
	if longTermRetentionPolicy.MonthlyRetention != nil {
		monthlyRetention = *longTermRetentionPolicy.MonthlyRetention
	}

	weeklyRetention := "PT0S"
	if longTermRetentionPolicy.WeeklyRetention != nil {
		weeklyRetention = *longTermRetentionPolicy.WeeklyRetention
	}

	weekOfYear := int32(1)
	if longTermRetentionPolicy.WeekOfYear != nil {
		weekOfYear = *longTermRetentionPolicy.WeekOfYear
	}

	yearlyRetention := "PT0S"
	if longTermRetentionPolicy.YearlyRetention != nil {
		yearlyRetention = *longTermRetentionPolicy.YearlyRetention
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

func ExpandShortTermRetentionPolicy(input []interface{}) *sql.BackupShortTermRetentionPolicyProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	shortTermRetentionPolicy := input[0].(map[string]interface{})

	shortTermPolicyProperties := sql.BackupShortTermRetentionPolicyProperties{
		RetentionDays: utils.Int32(7),
	}

	if v, ok := shortTermRetentionPolicy["retention_days"]; ok {
		shortTermPolicyProperties.RetentionDays = utils.Int32(int32(v.(int)))
	}

	return &shortTermPolicyProperties
}

func FlattenShortTermRetentionPolicy(shortTermRetentionPolicy *sql.BackupShortTermRetentionPolicy, d *pluginsdk.ResourceData) []interface{} {
	if shortTermRetentionPolicy == nil {
		return []interface{}{}
	}

	retentionDays := int32(7)
	if shortTermRetentionPolicy.RetentionDays != nil {
		retentionDays = *shortTermRetentionPolicy.RetentionDays
	}

	return []interface{}{
		map[string]interface{}{
			"retention_days": retentionDays,
		},
	}
}
