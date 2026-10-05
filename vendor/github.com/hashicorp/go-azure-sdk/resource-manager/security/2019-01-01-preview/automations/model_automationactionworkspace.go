package automations

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutomationAction = AutomationActionWorkspace{}

type AutomationActionWorkspace struct {
	WorkspaceResourceId *string `json:"workspaceResourceId,omitempty"`

	// Fields inherited from AutomationAction

	ActionType ActionType `json:"actionType"`
}

func (s AutomationActionWorkspace) AutomationAction() BaseAutomationActionImpl {
	return BaseAutomationActionImpl{
		ActionType: s.ActionType,
	}
}

var _ json.Marshaler = AutomationActionWorkspace{}

func (s AutomationActionWorkspace) MarshalJSON() ([]byte, error) {
	type wrapper AutomationActionWorkspace
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutomationActionWorkspace: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutomationActionWorkspace: %+v", err)
	}

	decoded["actionType"] = "Workspace"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutomationActionWorkspace: %+v", err)
	}

	return encoded, nil
}
