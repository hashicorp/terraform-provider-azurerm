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

			"match_values": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"HTTP",
						"HTTPS",
					}, false),
				},
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
			MatchValues:     utils.ExpandStringSlice(rsc["match_values"].(*schema.Set).List()),
		},
	}

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

		if params.MatchValues != nil {
			res["match_values"] = schema.NewSet(schema.HashString, utils.FlattenStringSlice(params.MatchValues))
		}
	}

	return res
}
