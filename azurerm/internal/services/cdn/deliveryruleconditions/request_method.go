package deliveryruleconditions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func RequestMethod() *schema.Resource {
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
						"DELETE",
						"GET",
						"HEAD",
						"OPTIONS",
						"POST",
						"PUT",
					}, false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionRequestMethod(hvc map[string]interface{}) *cdn.DeliveryRuleRequestMethodCondition {
	return &cdn.DeliveryRuleRequestMethodCondition{
		Name: cdn.NameRequestMethod,
		Parameters: &cdn.RequestMethodMatchConditionParameters{
			OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleRequestMethodConditionParameters"),
			Operator:        utils.String(hvc["operator"].(string)),
			NegateCondition: utils.Bool(hvc["negate_condition"].(bool)),
			MatchValues:     utils.ExpandStringSlice(hvc["match_values"].(*schema.Set).List()),
		},
	}
}

func FlattenArmCdnEndpointConditionRequestMethod(hvc *cdn.DeliveryRuleRequestMethodCondition) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := hvc.Parameters; params != nil {
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
