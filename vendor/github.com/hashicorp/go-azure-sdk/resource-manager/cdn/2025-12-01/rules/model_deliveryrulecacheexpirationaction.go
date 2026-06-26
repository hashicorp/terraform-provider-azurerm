package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleAction = DeliveryRuleCacheExpirationAction{}

type DeliveryRuleCacheExpirationAction struct {
	Parameters CacheExpirationActionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleAction

	Name DeliveryRuleActionName `json:"name"`
}

func (s DeliveryRuleCacheExpirationAction) DeliveryRuleAction() BaseDeliveryRuleActionImpl {
	return BaseDeliveryRuleActionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = DeliveryRuleCacheExpirationAction{}

func (s DeliveryRuleCacheExpirationAction) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleCacheExpirationAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleCacheExpirationAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleCacheExpirationAction: %+v", err)
	}

	decoded["name"] = "CacheExpiration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleCacheExpirationAction: %+v", err)
	}

	return encoded, nil
}
