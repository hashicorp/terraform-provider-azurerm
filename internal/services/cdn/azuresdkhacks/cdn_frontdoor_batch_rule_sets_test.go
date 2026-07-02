// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"encoding/json"
	"testing"

	batchRules "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
)

func TestBatchRuleProperties_MarshalJSON_IncludesNullCacheConfiguration(t *testing.T) {
	actions := []batchRules.DeliveryRuleAction{
		batchRules.DeliveryRuleRouteConfigurationOverrideAction{
			Name: batchRules.DeliveryRuleActionNameRouteConfigurationOverride,
			Parameters: batchRules.RouteConfigurationOverrideActionParameters{
				TypeName: batchRules.DeliveryRuleActionParametersTypeDeliveryRuleRouteConfigurationOverrideActionParameters,
			},
		},
	}

	input := BatchRuleProperties{
		Actions: &actions,
	}

	encoded, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("expected no error but got %q", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid json but got %q", err)
	}

	rawActions, ok := decoded["actions"].([]interface{})
	if !ok || len(rawActions) != 1 {
		t.Fatalf("expected one action but got %#v", decoded["actions"])
	}

	action, ok := rawActions[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected action object but got %#v", rawActions[0])
	}

	params, ok := action["parameters"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected parameters object but got %#v", action["parameters"])
	}

	cacheConfiguration, exists := params["cacheConfiguration"]
	if !exists {
		t.Fatalf("expected cacheConfiguration key in marshaled payload, got %#v", params)
	}

	if cacheConfiguration != nil {
		t.Fatalf("expected cacheConfiguration to be null but got %#v", cacheConfiguration)
	}
}

func TestBatchRuleProperties_MarshalJSON_IncludesNullOriginGroupOverride(t *testing.T) {
	actions := []batchRules.DeliveryRuleAction{
		batchRules.DeliveryRuleRouteConfigurationOverrideAction{
			Name: batchRules.DeliveryRuleActionNameRouteConfigurationOverride,
			Parameters: batchRules.RouteConfigurationOverrideActionParameters{
				TypeName: batchRules.DeliveryRuleActionParametersTypeDeliveryRuleRouteConfigurationOverrideActionParameters,
				CacheConfiguration: &batchRules.CacheConfiguration{
					CacheBehavior: &[]batchRules.RuleCacheBehavior{batchRules.RuleCacheBehaviorOverrideAlways}[0],
				},
			},
		},
	}

	input := BatchRuleProperties{
		Actions: &actions,
	}

	encoded, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("expected no error but got %q", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("expected valid json but got %q", err)
	}

	rawActions, ok := decoded["actions"].([]interface{})
	if !ok || len(rawActions) != 1 {
		t.Fatalf("expected one action but got %#v", decoded["actions"])
	}

	action, ok := rawActions[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected action object but got %#v", rawActions[0])
	}

	params, ok := action["parameters"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected parameters object but got %#v", action["parameters"])
	}

	originGroupOverride, exists := params["originGroupOverride"]
	if !exists {
		t.Fatalf("expected originGroupOverride key in marshaled payload, got %#v", params)
	}

	if originGroupOverride != nil {
		t.Fatalf("expected originGroupOverride to be null but got %#v", originGroupOverride)
	}
}
