package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleAction = DeliveryRuleCacheKeyQueryStringAction{}

type DeliveryRuleCacheKeyQueryStringAction struct {
	Parameters CacheKeyQueryStringActionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleAction

	Name DeliveryRuleActionName `json:"name"`
}

func (s DeliveryRuleCacheKeyQueryStringAction) DeliveryRuleAction() BaseDeliveryRuleActionImpl {
	return BaseDeliveryRuleActionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = DeliveryRuleCacheKeyQueryStringAction{}

func (s DeliveryRuleCacheKeyQueryStringAction) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleCacheKeyQueryStringAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleCacheKeyQueryStringAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleCacheKeyQueryStringAction: %+v", err)
	}

	decoded["name"] = "CacheKeyQueryString"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleCacheKeyQueryStringAction: %+v", err)
	}

	return encoded, nil
}
