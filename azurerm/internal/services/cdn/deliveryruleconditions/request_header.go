package deliveryruleconditions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func RequestHeader() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"selector": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"operator": {
				Type:     pluginsdk.TypeString,
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
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"match_values": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},

			"transforms": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
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
				MatchValues:     utils.ExpandStringSlice(item["match_values"].(*pluginsdk.Set).List()),
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

func FlattenArmCdnEndpointConditionRequestHeader(input cdn.BasicDeliveryRuleCondition) (*map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestHeaderCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request header condition")
	}

	selector := ""
	operator := ""
	matchValues := make([]interface{}, 0)
	negateCondition := false
	transforms := make([]string, 0)
	if params := condition.Parameters; params != nil {
		if params.Selector != nil {
			selector = *params.Selector
		}

		operator = string(params.Operator)

		if params.NegateCondition != nil {
			negateCondition = *params.NegateCondition
		}

		if params.MatchValues != nil {
			matchValues = utils.FlattenStringSlice(params.MatchValues)
		}

		if params.Transforms != nil {
			for _, transform := range *params.Transforms {
				transforms = append(transforms, string(transform))
			}
		}
	}

	return &map[string]interface{}{
		"selector":         selector,
		"operator":         operator,
		"match_values":     pluginsdk.NewSet(pluginsdk.HashString, matchValues),
		"negate_condition": negateCondition,
		"transforms":       transforms,
	}, nil
}
