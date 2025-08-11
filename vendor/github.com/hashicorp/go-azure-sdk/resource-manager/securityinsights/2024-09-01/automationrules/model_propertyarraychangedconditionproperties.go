package automationrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutomationRuleCondition = PropertyArrayChangedConditionProperties{}

type PropertyArrayChangedConditionProperties struct {
	ConditionProperties *AutomationRulePropertyArrayChangedValuesCondition `json:"conditionProperties,omitempty"`

	// Fields inherited from AutomationRuleCondition

	ConditionType ConditionType `json:"conditionType"`
}

func (s PropertyArrayChangedConditionProperties) AutomationRuleCondition() BaseAutomationRuleConditionImpl {
	return BaseAutomationRuleConditionImpl{
		ConditionType: s.ConditionType,
	}
}

var _ json.Marshaler = PropertyArrayChangedConditionProperties{}

func (s PropertyArrayChangedConditionProperties) MarshalJSON() ([]byte, error) {
	type wrapper PropertyArrayChangedConditionProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PropertyArrayChangedConditionProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PropertyArrayChangedConditionProperties: %+v", err)
	}

	decoded["conditionType"] = "PropertyArrayChanged"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PropertyArrayChangedConditionProperties: %+v", err)
	}

	return encoded, nil
}
