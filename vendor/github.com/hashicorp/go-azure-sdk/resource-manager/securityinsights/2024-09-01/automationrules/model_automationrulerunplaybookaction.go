package automationrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutomationRuleAction = AutomationRuleRunPlaybookAction{}

type AutomationRuleRunPlaybookAction struct {
	ActionConfiguration *PlaybookActionProperties `json:"actionConfiguration,omitempty"`

	// Fields inherited from AutomationRuleAction

	ActionType ActionType `json:"actionType"`
	Order      int64      `json:"order"`
}

func (s AutomationRuleRunPlaybookAction) AutomationRuleAction() BaseAutomationRuleActionImpl {
	return BaseAutomationRuleActionImpl{
		ActionType: s.ActionType,
		Order:      s.Order,
	}
}

var _ json.Marshaler = AutomationRuleRunPlaybookAction{}

func (s AutomationRuleRunPlaybookAction) MarshalJSON() ([]byte, error) {
	type wrapper AutomationRuleRunPlaybookAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutomationRuleRunPlaybookAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutomationRuleRunPlaybookAction: %+v", err)
	}

	decoded["actionType"] = "RunPlaybook"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutomationRuleRunPlaybookAction: %+v", err)
	}

	return encoded, nil
}
