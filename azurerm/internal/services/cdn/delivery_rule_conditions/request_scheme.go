package delivery_rule_conditions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func RuleConditionRequestScheme() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Equal",
				ValidateFunc: validation.StringInSlice([]string{
					"Equal",
				}, false),
			},

			"negate_condition": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"match_values": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP",
					"HTTPS",
				}, false),
			},
		},
	}
}

func ExpandArmCdnEndpointConditionRequestScheme(rsc map[string]interface{}) *cdn.DeliveryRuleRequestSchemeCondition {
	requestSchemeCondition := cdn.DeliveryRuleRequestSchemeCondition{
		Name: cdn.NameRequestScheme,
	}

	matchValues := rsc["match_values"].([]string)
	params := cdn.RequestSchemeMatchConditionParameters{
		MatchValues: &matchValues,
	}

	if operator, ok := rsc["operator"]; ok {
		params.Operator = utils.String(operator.(string))
	}

	if negate, ok := rsc["negate_condition"]; ok {
		params.NegateCondition = utils.Bool(negate.(bool))
	}

	requestSchemeCondition.Parameters = &params

	return &requestSchemeCondition
}

func FlattenArmCdnEndpointConditionRequestScheme(condition *cdn.DeliveryRuleRequestSchemeCondition) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := condition.Parameters; params != nil {
		if params.Operator != nil {
			res["operator"] = *params.Operator
		}

		if params.NegateCondition != nil {
			res["negate_condition"] = *params.NegateCondition
		}

		if params.MatchValues != nil {
			res["match_values"] = *params.MatchValues
		}
	}

	return res
}
