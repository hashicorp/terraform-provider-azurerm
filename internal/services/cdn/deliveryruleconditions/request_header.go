// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package deliveryruleconditions

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
				ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRequestHeaderOperator(),
					false),
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
						string(rules.TransformLowercase),
						string(rules.TransformUppercase),
					}, false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionRequestHeader(input []interface{}) []rules.DeliveryRuleCondition {
	output := make([]rules.DeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		requestHeaderCondition := rules.DeliveryRuleRequestHeaderCondition{
			Name: rules.MatchVariableRequestHeader,
			Parameters: rules.RequestHeaderMatchConditionParameters{
				TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestHeaderConditionParameters,
				Selector:        pointer.To(item["selector"].(string)),
				Operator:        rules.RequestHeaderOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].(*pluginsdk.Set).List()),
				Transforms:      expandTransforms(item["transforms"].([]interface{})),
			},
		}

		output = append(output, requestHeaderCondition)
	}

	return output
}

func FlattenArmCdnEndpointConditionRequestHeader(input rules.DeliveryRuleCondition) (*map[string]interface{}, error) {
	condition, ok := AsDeliveryRuleRequestHeaderCondition(input)
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request header condition")
	}

	selector := ""
	operator := ""
	matchValues := make([]interface{}, 0)
	negateCondition := false
	transforms := make([]string, 0)

	if params := condition; params != nil {
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
			transforms = flattenTransforms(params.Transforms)
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
