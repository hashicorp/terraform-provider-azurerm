package deliveryruleconditions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func RequestHeader() *schema.Resource {
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
					string(cdn.RequestHeaderOperatorAny),
					string(cdn.RequestHeaderOperatorBeginsWith),
					string(cdn.RequestHeaderOperatorContains),
					string(cdn.RequestHeaderOperatorEndsWith),
					string(cdn.RequestHeaderOperatorEqual),
					string(cdn.RequestHeaderOperatorGreaterThan),
					string(cdn.RequestHeaderOperatorGreaterThanOrEqual),
					string(cdn.RequestHeaderOperatorLessThan),
					string(cdn.RequestHeaderOperatorLessThanOrEqual),
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

func ExpandArmCdnEndpointConditionRequestHeader(input []interface{}) []cdn.BasicDeliveryRuleCondition {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		requestHeaderCondition := cdn.DeliveryRuleRequestHeaderCondition{
			Name: cdn.NameRequestHeader,
			Parameters: &cdn.RequestHeaderMatchConditionParameters{
				OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleRequestHeaderConditionParameters"),
				Selector:        utils.String(item["selector"].(string)),
				Operator:        cdn.RequestHeaderOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].(*schema.Set).List()),
			},
		}

		if rawTransforms := item["transforms"].([]interface{}); len(rawTransforms) != 0 {
			transforms := make([]cdn.Transform, 0)
			for _, t := range rawTransforms {
				transforms = append(transforms, cdn.Transform(t.(string)))
			}
			requestHeaderCondition.Parameters.Transforms = &transforms
		}

		output = append(output, requestHeaderCondition)
	}

	return output
}

func FlattenArmCdnEndpointConditionRequestHeader(cc *cdn.DeliveryRuleRequestHeaderCondition) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := cc.Parameters; params != nil {
		if params.Selector != nil {
			res["selector"] = *params.Selector
		}

		res["operator"] = string(params.Operator)

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
