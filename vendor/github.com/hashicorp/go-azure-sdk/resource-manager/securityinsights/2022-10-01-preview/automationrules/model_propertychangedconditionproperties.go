package automationrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutomationRuleCondition = PropertyChangedConditionProperties{}

type PropertyChangedConditionProperties struct {
	ConditionProperties *AutomationRulePropertyValuesChangedCondition `json:"conditionProperties,omitempty"`

	// Fields inherited from AutomationRuleCondition

	ConditionType ConditionType `json:"conditionType"`
}

func (s PropertyChangedConditionProperties) AutomationRuleCondition() BaseAutomationRuleConditionImpl {
	return BaseAutomationRuleConditionImpl{
		ConditionType: s.ConditionType,
	}
}

var _ json.Marshaler = PropertyChangedConditionProperties{}

func (s PropertyChangedConditionProperties) MarshalJSON() ([]byte, error) {
	type wrapper PropertyChangedConditionProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PropertyChangedConditionProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PropertyChangedConditionProperties: %+v", err)
	}

	decoded["conditionType"] = "PropertyChanged"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PropertyChangedConditionProperties: %+v", err)
	}

	return encoded, nil
}
