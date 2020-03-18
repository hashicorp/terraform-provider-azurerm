package helper

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SQLLongTermRententionPolicy() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// WeeklyRetention - The weekly retention policy for an LTR backup in an ISO 8601 format. 1-520 weeks (P520W) P520D
				"weekly_retention": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				// MonthlyRetention - The monthly retention policy for an LTR backup in an ISO 8601 format. 4-520 weeks P520W P120M
				"monthly_retention": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty, // Validate ISO8601 at least 1 month (4 weeks/ 30 days?)
				},
				// YearlyRetention - The yearly retention policy for an LTR backup in an ISO 8601 format. 52-520 weeks P52W P5Y
				"yearly_retention": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty, // Validate ISO8601 at least 1 year (12 months/52 weeks/365 days)
				},
				// WeekOfYear - The week of year to take the yearly backup in an ISO 8601 format. 1-52
				"week_of_year": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 52),
				},
			},
		},
	}
}

func SQLShortTermRetentionPolicy() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// RetentionDays - The backup retention period in days. This is how many days Point-in-Time Restore will be supported.
				"retention_days": {
					Type:     schema.TypeInt,
					Optional: true,
					// ValidateFunc: validation.IntBetween(0, 1000),
					ValidateFunc: validation.IntInSlice([]int{
						7,
						14,
						21,
						28,
						35,
					}),
				},
			},
		},
	}
}

func ExpandSQLLongTermRetentionPolicyProperties(input []interface{}) *sql.LongTermRetentionPolicyProperties {
	LongTermRetentionPolicyProperties := sql.LongTermRetentionPolicyProperties{
		WeeklyRetention:  utils.String("P0W"),
		MonthlyRetention: utils.String("P0W"),
		YearlyRetention:  utils.String("P0W"),
		WeekOfYear:       utils.Int32(1),
	}

	if len(input) == 0 {
		return &LongTermRetentionPolicyProperties
	}

	longTermPolicies := input[0].(map[string]interface{})

	if v, ok := longTermPolicies["weekly_retention"]; ok {
		LongTermRetentionPolicyProperties.WeeklyRetention = utils.String(v.(string))
	}

	if v, ok := longTermPolicies["monthly_retention"]; ok {
		LongTermRetentionPolicyProperties.MonthlyRetention = utils.String(v.(string))
	}

	if v, ok := longTermPolicies["yearly_retention"]; ok {
		LongTermRetentionPolicyProperties.YearlyRetention = utils.String(v.(string))
	}

	if v, ok := longTermPolicies["week_of_year"]; ok {
		LongTermRetentionPolicyProperties.WeekOfYear = utils.Int32(int32(v.(int)))
	}

	return &LongTermRetentionPolicyProperties
}

func FlattenSQLLongTermRetentionPolicyProperties(longTermRetentionPolicy *sql.BackupLongTermRetentionPolicy) []interface{} {
	if longTermRetentionPolicy == nil || longTermRetentionPolicy.LongTermRetentionPolicyProperties == nil {
		return []interface{}{}
	}

	var weeklyRetention string
	if longTermRetentionPolicy.WeeklyRetention != nil {
		weeklyRetention = *longTermRetentionPolicy.WeeklyRetention
	}

	var monthlyRetention string
	if longTermRetentionPolicy.MonthlyRetention != nil {
		monthlyRetention = *longTermRetentionPolicy.MonthlyRetention
	}

	var yearlyRetention string
	if longTermRetentionPolicy.YearlyRetention != nil {
		yearlyRetention = *longTermRetentionPolicy.YearlyRetention
	}

	var weekOfYear int32
	if longTermRetentionPolicy.WeekOfYear != nil {
		weekOfYear = *longTermRetentionPolicy.WeekOfYear
	}

	return []interface{}{
		map[string]interface{}{
			"weekly_retention":  weeklyRetention,
			"monthly_retention": monthlyRetention,
			"yearly_retention":  yearlyRetention,
			"week_of_year":      weekOfYear,
		},
	}
}

func ExpandSQLShortTermRetentionPolicyProperties(input []interface{}) *sql.BackupShortTermRetentionPolicyProperties {
	ShortTermRetentionPolicyProperties := sql.BackupShortTermRetentionPolicyProperties{
		RetentionDays: utils.Int32(1),
	}

	if len(input) == 0 {
		return &ShortTermRetentionPolicyProperties
	}

	shortTermPolicies := input[0].(map[string]interface{})

	if v, ok := shortTermPolicies["retention_days"]; ok {
		ShortTermRetentionPolicyProperties.RetentionDays = utils.Int32(int32(v.(int)))
	}

	return &ShortTermRetentionPolicyProperties
}

func FlattenSQLShortTermRetentionPolicy(shortTermRetentionPolicy *sql.BackupShortTermRetentionPolicy) []interface{} {
	if shortTermRetentionPolicy == nil || shortTermRetentionPolicy.BackupShortTermRetentionPolicyProperties == nil {
		return []interface{}{}
	}

	var retentionDays int32
	if shortTermRetentionPolicy.RetentionDays != nil {
		retentionDays = *shortTermRetentionPolicy.RetentionDays
	}

	return []interface{}{
		map[string]interface{}{
			"retention_days": retentionDays,
		},
	}
}
