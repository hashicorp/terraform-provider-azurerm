package monitor

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/Azure/azure-sdk-for-go/services/preview/alertsmanagement/mgmt/2019-05-05/alertsmanagement"
)

var weekDays = []string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

var weekDayMap = map[string]int{
	"Sunday":    0,
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
}

func schemaActionRuleAlertContextCondtion() *schema.Schema {
	return schemaActionRuleCondtion(
		validation.StringInSlice([]string{
			string(alertsmanagement.Equals),
			string(alertsmanagement.NotEquals),
			string(alertsmanagement.Contains),
			string(alertsmanagement.DoesNotContain),
		}, false),
		validation.StringIsNotEmpty,
	)
}

func schemaActionRuleAlertRuleIDCondtion() *schema.Schema {
	return schemaActionRuleCondtion(
		validation.StringInSlice([]string{
			string(alertsmanagement.Equals),
			string(alertsmanagement.NotEquals),
			string(alertsmanagement.Contains),
			string(alertsmanagement.DoesNotContain),
		}, false),
		validation.StringIsNotEmpty,
	)
}

func schemaActionRuleDescriptionCondtion() *schema.Schema {
	return schemaActionRuleCondtion(
		validation.StringInSlice([]string{
			string(alertsmanagement.Equals),
			string(alertsmanagement.NotEquals),
			string(alertsmanagement.Contains),
			string(alertsmanagement.DoesNotContain),
		}, false),
		validation.StringIsNotEmpty,
	)
}

func schemaActionRuleMonitorCondtion() *schema.Schema {
	return schemaActionRuleCondtion(
		validation.StringInSlice([]string{
			string(alertsmanagement.Equals),
			string(alertsmanagement.NotEquals),
		}, false),
		validation.StringInSlice([]string{
			string(alertsmanagement.Fired),
			string(alertsmanagement.Resolved),
		}, false),
	)
}

func schemaActionRuleMonitorServiceCondtion() *schema.Schema {
	return schemaActionRuleCondtion(
		validation.StringInSlice([]string{
			string(alertsmanagement.Equals),
			string(alertsmanagement.NotEquals),
		}, false),
		// the supported type list is not consistent with the swagger and sdk
		// https://github.com/Azure/azure-rest-api-specs/issues/9076
		// directly use string constant
		validation.StringInSlice([]string{
			"ActivityLog Administrative",
			"ActivityLog Autoscale",
			"ActivityLog Policy",
			"ActivityLog Recommendation",
			"ActivityLog Security",
			"Application Insights",
			"Azure Backup",
			"Data Box Edge",
			"Data Box Gateway",
			"Health Platform",
			"Log Analytics",
			"Platform",
			"Resource Health",
		}, false),
	)
}

func schemaActionRuleSeverityCondtion() *schema.Schema {
	return schemaActionRuleCondtion(
		validation.StringInSlice([]string{
			string(alertsmanagement.Equals),
			string(alertsmanagement.NotEquals),
		}, false),
		validation.StringInSlice([]string{
			string(alertsmanagement.Sev0),
			string(alertsmanagement.Sev1),
			string(alertsmanagement.Sev2),
			string(alertsmanagement.Sev3),
			string(alertsmanagement.Sev4),
		}, false),
	)
}

func schemaActionRuleTargetResourceTypeCondtion() *schema.Schema {
	return schemaActionRuleCondtion(
		validation.StringInSlice([]string{
			string(alertsmanagement.Equals),
			string(alertsmanagement.NotEquals),
		}, false),
		validation.StringIsNotEmpty,
	)
}

func schemaActionRuleCondtion(operatorValidateFunc, valuesValidateFunc schema.SchemaValidateFunc) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"operator": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: operatorValidateFunc,
				},

				"values": {
					Type:     schema.TypeSet,
					Required: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: valuesValidateFunc,
					},
				},
			},
		},
	}
}

func expandArmActionRuleCondition(input []interface{}) *alertsmanagement.Condition {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &alertsmanagement.Condition{
		Operator: alertsmanagement.Operator(v["operator"].(string)),
		Values:   utils.ExpandStringSlice(v["values"].(*schema.Set).List()),
	}
}

func flattenArmActionRuleCondition(input *alertsmanagement.Condition) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var operator string
	if input.Operator != "" {
		operator = string(input.Operator)
	}
	return []interface{}{
		map[string]interface{}{
			"operator": operator,
			"values":   utils.FlattenStringSlice(input.Values),
		},
	}
}
