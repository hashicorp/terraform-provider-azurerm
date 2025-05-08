package automationrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutomationRuleCondition = PropertyArrayConditionProperties{}

type PropertyArrayConditionProperties struct {
	ConditionProperties *AutomationRulePropertyArrayValuesCondition `json:"conditionProperties,omitempty"`

	// Fields inherited from AutomationRuleCondition

	ConditionType ConditionType `json:"conditionType"`
}

func (s PropertyArrayConditionProperties) AutomationRuleCondition() BaseAutomationRuleConditionImpl {
	return BaseAutomationRuleConditionImpl{
		ConditionType: s.ConditionType,
	}
}

var _ json.Marshaler = PropertyArrayConditionProperties{}

func (s PropertyArrayConditionProperties) MarshalJSON() ([]byte, error) {
	type wrapper PropertyArrayConditionProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PropertyArrayConditionProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PropertyArrayConditionProperties: %+v", err)
	}

	decoded["conditionType"] = "PropertyArray"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PropertyArrayConditionProperties: %+v", err)
	}

	return encoded, nil
}
