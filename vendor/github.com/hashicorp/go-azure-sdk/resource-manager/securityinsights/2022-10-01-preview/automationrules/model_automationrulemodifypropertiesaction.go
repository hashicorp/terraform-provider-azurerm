package automationrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutomationRuleAction = AutomationRuleModifyPropertiesAction{}

type AutomationRuleModifyPropertiesAction struct {
	ActionConfiguration *IncidentPropertiesAction `json:"actionConfiguration,omitempty"`

	// Fields inherited from AutomationRuleAction

	ActionType ActionType `json:"actionType"`
	Order      int64      `json:"order"`
}

func (s AutomationRuleModifyPropertiesAction) AutomationRuleAction() BaseAutomationRuleActionImpl {
	return BaseAutomationRuleActionImpl{
		ActionType: s.ActionType,
		Order:      s.Order,
	}
}

var _ json.Marshaler = AutomationRuleModifyPropertiesAction{}

func (s AutomationRuleModifyPropertiesAction) MarshalJSON() ([]byte, error) {
	type wrapper AutomationRuleModifyPropertiesAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutomationRuleModifyPropertiesAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutomationRuleModifyPropertiesAction: %+v", err)
	}

	decoded["actionType"] = "ModifyProperties"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutomationRuleModifyPropertiesAction: %+v", err)
	}

	return encoded, nil
}
