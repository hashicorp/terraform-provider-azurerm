package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleCondition = DeliveryRuleCookiesCondition{}

type DeliveryRuleCookiesCondition struct {
	Parameters CookiesMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition

	Name MatchVariable `json:"name"`
}

func (s DeliveryRuleCookiesCondition) DeliveryRuleCondition() BaseDeliveryRuleConditionImpl {
	return BaseDeliveryRuleConditionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = DeliveryRuleCookiesCondition{}

func (s DeliveryRuleCookiesCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleCookiesCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleCookiesCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleCookiesCondition: %+v", err)
	}

	decoded["name"] = "Cookies"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleCookiesCondition: %+v", err)
	}

	return encoded, nil
}
