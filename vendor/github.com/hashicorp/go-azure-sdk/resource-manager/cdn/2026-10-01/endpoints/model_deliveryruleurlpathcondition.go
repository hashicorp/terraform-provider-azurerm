package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleCondition = DeliveryRuleURLPathCondition{}

type DeliveryRuleURLPathCondition struct {
	Parameters URLPathMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition

	Name MatchVariable `json:"name"`
}

func (s DeliveryRuleURLPathCondition) DeliveryRuleCondition() BaseDeliveryRuleConditionImpl {
	return BaseDeliveryRuleConditionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = DeliveryRuleURLPathCondition{}

func (s DeliveryRuleURLPathCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleURLPathCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleURLPathCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleURLPathCondition: %+v", err)
	}

	decoded["name"] = "UrlPath"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleURLPathCondition: %+v", err)
	}

	return encoded, nil
}
