package automationrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutomationRuleCondition = BooleanConditionProperties{}

type BooleanConditionProperties struct {
	ConditionProperties *AutomationRuleBooleanCondition `json:"conditionProperties,omitempty"`

	// Fields inherited from AutomationRuleCondition

	ConditionType ConditionType `json:"conditionType"`
}

func (s BooleanConditionProperties) AutomationRuleCondition() BaseAutomationRuleConditionImpl {
	return BaseAutomationRuleConditionImpl{
		ConditionType: s.ConditionType,
	}
}

var _ json.Marshaler = BooleanConditionProperties{}

func (s BooleanConditionProperties) MarshalJSON() ([]byte, error) {
	type wrapper BooleanConditionProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BooleanConditionProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BooleanConditionProperties: %+v", err)
	}

	decoded["conditionType"] = "Boolean"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BooleanConditionProperties: %+v", err)
	}

	return encoded, nil
}
