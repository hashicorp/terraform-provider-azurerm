package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleCondition = DeliveryRuleURLFileExtensionCondition{}

type DeliveryRuleURLFileExtensionCondition struct {
	Parameters URLFileExtensionMatchConditionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleCondition

	Name MatchVariable `json:"name"`
}

func (s DeliveryRuleURLFileExtensionCondition) DeliveryRuleCondition() BaseDeliveryRuleConditionImpl {
	return BaseDeliveryRuleConditionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = DeliveryRuleURLFileExtensionCondition{}

func (s DeliveryRuleURLFileExtensionCondition) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleURLFileExtensionCondition
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleURLFileExtensionCondition: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleURLFileExtensionCondition: %+v", err)
	}

	decoded["name"] = "UrlFileExtension"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleURLFileExtensionCondition: %+v", err)
	}

	return encoded, nil
}
