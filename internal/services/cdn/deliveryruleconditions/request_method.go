// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package deliveryruleconditions

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func RequestMethod() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"operator": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "Equal",
				ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRequestMethodOperator(),
					false),
			},

			"negate_condition": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"match_values": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRequestMethodMatchValue(),
						false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionRequestMethod(input []interface{}) []rules.DeliveryRuleCondition {
	output := make([]rules.DeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		output = append(output, rules.DeliveryRuleRequestMethodCondition{
			Name: rules.MatchVariableRequestMethod,
			Parameters: rules.RequestMethodMatchConditionParameters{
				TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestMethodConditionParameters,
				Operator:        rules.RequestMethodOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     expandRequestMethodMatchValue(item["match_values"].(*pluginsdk.Set).List()),
			},
		})
	}

	return output
}

func FlattenArmCdnEndpointConditionRequestMethod(input rules.DeliveryRuleCondition) (*map[string]interface{}, error) {
	condition, ok := AsDeliveryRuleRequestMethodCondition(input)
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request method condition")
	}

	operator := ""
	negateCondition := false
	matchValues := make([]interface{}, 0)

	if params := condition; params != nil {
		operator = string(params.Operator)

		if params.NegateCondition != nil {
			negateCondition = *params.NegateCondition
		}

		if params.MatchValues != nil {
			matchValues = flattenRequestMethodMatchValue(params.MatchValues)
		}
	}

	return &map[string]interface{}{
		"match_values":     pluginsdk.NewSet(pluginsdk.HashString, matchValues),
		"negate_condition": negateCondition,
		"operator":         operator,
	}, nil
}
