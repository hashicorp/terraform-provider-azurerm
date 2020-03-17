package helper

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func SqlLongTermRententionPolicy() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// TODO Update validation
				// WeeklyRetention - The weekly retention policy for an LTR backup in an ISO 8601 format.
				"weekly_retention": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				// MonthlyRetention - The monthly retention policy for an LTR backup in an ISO 8601 format.
				"monthly_retention": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				// YearlyRetention - The yearly retention policy for an LTR backup in an ISO 8601 format.
				"yearly_retention": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				// WeekOfYear - The week of year to take the yearly backup in an ISO 8601 format.
				"week_of_year": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 1000),
				},
			},
		},
	}
}
