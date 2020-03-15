package deliveryruleconditions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func RequestScheme() *schema.Resource {
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

			"match_value": {
				Type:     schema.TypeString,
				Required: true,
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
		Parameters: &cdn.RequestSchemeMatchConditionParameters{
			OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleRequestSchemeConditionParameters"),
			NegateCondition: utils.Bool(rsc["negate_condition"].(bool)),
		},
	}

	matchValues := []string{rsc["match_value"].(string)}
	requestSchemeCondition.Parameters.MatchValues = &matchValues

	if operator := rsc["operator"]; operator.(string) != "" {
		requestSchemeCondition.Parameters.Operator = utils.String(operator.(string))
	}

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

		if params.MatchValues != nil && len(*params.MatchValues) > 0 {
			res["match_value"] = (*params.MatchValues)[0]
		}
	}

	return res
}
