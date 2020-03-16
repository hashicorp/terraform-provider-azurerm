package deliveryruleconditions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func PostArg() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"selector": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"operator": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.PostArgsOperatorAny),
					string(cdn.PostArgsOperatorBeginsWith),
					string(cdn.PostArgsOperatorContains),
					string(cdn.PostArgsOperatorEndsWith),
					string(cdn.PostArgsOperatorEqual),
					string(cdn.PostArgsOperatorGreaterThan),
					string(cdn.PostArgsOperatorGreaterThanOrEqual),
					string(cdn.PostArgsOperatorLessThan),
					string(cdn.PostArgsOperatorLessThanOrEqual),
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
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},

			"transforms": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(cdn.Lowercase),
						string(cdn.Uppercase),
					}, false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionPostArg(pac map[string]interface{}) *cdn.DeliveryRulePostArgsCondition {
	cookiesCondition := cdn.DeliveryRulePostArgsCondition{
		Name: cdn.NameCookies,
		Parameters: &cdn.PostArgsMatchConditionParameters{
			OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRulePostArgsConditionParameters"),
			Selector:        utils.String(pac["selector"].(string)),
			Operator:        cdn.PostArgsOperator(pac["operator"].(string)),
			NegateCondition: utils.Bool(pac["negate_condition"].(bool)),
			MatchValues:     utils.ExpandStringSlice(pac["match_values"].(*schema.Set).List()),
		},
	}

	if rawTransforms := pac["transforms"].([]interface{}); len(rawTransforms) != 0 {
		transforms := make([]cdn.Transform, 0)
		for _, t := range rawTransforms {
			transforms = append(transforms, cdn.Transform(t.(string)))
		}
		cookiesCondition.Parameters.Transforms = &transforms
	}

	return &cookiesCondition
}

func FlattenArmCdnEndpointConditionPostArg(pac *cdn.DeliveryRulePostArgsCondition) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := pac.Parameters; params != nil {
		if params.Selector != nil {
			res["selector"] = *params.Selector
		}

		res["operator"] = params.Operator

		if params.NegateCondition != nil {
			res["negate_condition"] = *params.NegateCondition
		}

		if params.MatchValues != nil {
			res["match_values"] = schema.NewSet(schema.HashString, utils.FlattenStringSlice(params.MatchValues))
		}

		if params.Transforms != nil {
			transforms := make([]string, 0)
			for _, transform := range *params.Transforms {
				transforms = append(transforms, string(transform))
			}
			res["transforms"] = &transforms
		}
	}

	return res
}
