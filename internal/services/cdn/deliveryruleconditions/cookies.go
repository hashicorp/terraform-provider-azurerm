// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package deliveryruleconditions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func Cookies() *pluginsdk.Resource {
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
					string(cdn.OperatorAny),
					string(cdn.OperatorBeginsWith),
					string(cdn.OperatorContains),
					string(cdn.OperatorEndsWith),
					string(cdn.OperatorEqual),
					string(cdn.OperatorGreaterThan),
					string(cdn.OperatorGreaterThanOrEqual),
					string(cdn.OperatorLessThan),
					string(cdn.OperatorLessThanOrEqual),
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
						string(cdn.TransformLowercase),
						string(cdn.TransformUppercase),
					}, false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionCookies(input []interface{}) []cdn.BasicDeliveryRuleCondition {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	for _, v := range input {
		item := v.(map[string]interface{})
		cookiesCondition := cdn.DeliveryRuleCookiesCondition{
			Name: cdn.NameCookies,
			Parameters: &cdn.CookiesMatchConditionParameters{
				OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleCookiesConditionParameters"),
				Selector:        utils.String(item["selector"].(string)),
				Operator:        cdn.CookiesOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].(*pluginsdk.Set).List()),
			},
		}

		if rawTransforms := item["transforms"].([]interface{}); len(rawTransforms) != 0 {
			transforms := make([]cdn.Transform, 0)
			for _, t := range rawTransforms {
				transforms = append(transforms, cdn.Transform(t.(string)))
			}
			cookiesCondition.Parameters.Transforms = &transforms
		}

		output = append(output, cookiesCondition)
	}

	return output
}

func FlattenArmCdnEndpointConditionCookies(input cdn.BasicDeliveryRuleCondition) (*map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleCookiesCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule cookie condition")
	}

	selector := ""
	operator := ""
	negateCondition := false
	matchValues := make([]interface{}, 0)
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
		"match_values":     pluginsdk.NewSet(pluginsdk.HashString, matchValues),
		"negate_condition": negateCondition,
		"operator":         operator,
		"selector":         selector,
		"transforms":       transforms,
	}, nil
}
