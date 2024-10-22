package automations

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutomationAction = AutomationActionLogicApp{}

type AutomationActionLogicApp struct {
	LogicAppResourceId *string `json:"logicAppResourceId,omitempty"`
	Uri                *string `json:"uri,omitempty"`

	// Fields inherited from AutomationAction

	ActionType ActionType `json:"actionType"`
}

func (s AutomationActionLogicApp) AutomationAction() BaseAutomationActionImpl {
	return BaseAutomationActionImpl{
		ActionType: s.ActionType,
	}
}

var _ json.Marshaler = AutomationActionLogicApp{}

func (s AutomationActionLogicApp) MarshalJSON() ([]byte, error) {
	type wrapper AutomationActionLogicApp
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutomationActionLogicApp: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutomationActionLogicApp: %+v", err)
	}

	decoded["actionType"] = "LogicApp"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutomationActionLogicApp: %+v", err)
	}

	return encoded, nil
}
