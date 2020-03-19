package helper

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SqlLongTermRetentionPolicy() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// WeeklyRetention - The weekly retention policy for an LTR backup in an ISO 8601 format. 1-520 weeks
				"weekly_retention": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: azure.ValidateLongTermRetentionPoliciesIsoFormat,
				},
				// MonthlyRetention - The monthly retention policy for an LTR backup in an ISO 8601 format. 4-520 weeks
				"monthly_retention": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: azure.ValidateLongTermRetentionPoliciesIsoFormat,
				},
				// YearlyRetention - The yearly retention policy for an LTR backup in an ISO 8601 format. 52-520 weeks
				"yearly_retention": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: azure.ValidateLongTermRetentionPoliciesIsoFormat,
				},
				// WeekOfYear - The week of year to take the yearly backup in an ISO 8601 format. 1-52
				"week_of_year": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 52),
				},
			},
		},
	}
}

func SqlShortTermRetentionPolicy() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// RetentionDays - The backup retention period in days. This is how many days Point-in-Time Restore will be supported.
				"retention_days": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(7, 35),
				},
			},
		},
	}
}

func ExpandSqlLongTermRetentionPolicyProperties(input []interface{}) *sql.LongTermRetentionPolicyProperties {
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

func FlattenSqlLongTermRetentionPolicy(longTermRetentionPolicy *sql.BackupLongTermRetentionPolicy) []interface{} {
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
