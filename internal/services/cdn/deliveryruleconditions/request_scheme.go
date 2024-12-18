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

func RequestScheme() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"operator": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(rules.OperatorEqual),
				ValidateFunc: validation.StringInSlice([]string{
					string(rules.OperatorEqual),
				}, false),
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
					ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRequestSchemeMatchValue(),
						false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionRequestScheme(input []interface{}) []rules.DeliveryRuleCondition {
	output := make([]rules.DeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		requestSchemeCondition := rules.DeliveryRuleRequestSchemeCondition{
			Name: rules.MatchVariableRequestScheme,
			Parameters: rules.RequestSchemeMatchConditionParameters{
				TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestSchemeConditionParameters,
				Operator:        rules.Operator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     expandRequestSchemeMatchValue(item["match_values"].(*pluginsdk.Set).List()),
			},
		}

		output = append(output, requestSchemeCondition)
	}

	return output
}

func FlattenArmCdnEndpointConditionRequestScheme(input rules.DeliveryRuleCondition) (*map[string]interface{}, error) {
	condition, ok := AsDeliveryRuleRequestSchemeCondition(input)
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request scheme condition")
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
			matchValues = flattenRequestSchemeMatchValue(params.MatchValues)
		}
	}

	return &map[string]interface{}{
		"operator":         operator,
		"match_values":     pluginsdk.NewSet(pluginsdk.HashString, matchValues),
		"negate_condition": negateCondition,
	}, nil
}
