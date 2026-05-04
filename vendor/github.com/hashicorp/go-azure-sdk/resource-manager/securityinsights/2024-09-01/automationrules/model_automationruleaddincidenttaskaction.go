package automationrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutomationRuleAction = AutomationRuleAddIncidentTaskAction{}

type AutomationRuleAddIncidentTaskAction struct {
	ActionConfiguration *AddIncidentTaskActionProperties `json:"actionConfiguration,omitempty"`

	// Fields inherited from AutomationRuleAction

	ActionType ActionType `json:"actionType"`
	Order      int64      `json:"order"`
}

func (s AutomationRuleAddIncidentTaskAction) AutomationRuleAction() BaseAutomationRuleActionImpl {
	return BaseAutomationRuleActionImpl{
		ActionType: s.ActionType,
		Order:      s.Order,
	}
}

var _ json.Marshaler = AutomationRuleAddIncidentTaskAction{}

func (s AutomationRuleAddIncidentTaskAction) MarshalJSON() ([]byte, error) {
	type wrapper AutomationRuleAddIncidentTaskAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutomationRuleAddIncidentTaskAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutomationRuleAddIncidentTaskAction: %+v", err)
	}

	decoded["actionType"] = "AddIncidentTask"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutomationRuleAddIncidentTaskAction: %+v", err)
	}

	return encoded, nil
}
