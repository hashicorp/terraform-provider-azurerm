package automations

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutomationAction = AutomationActionEventHub{}

type AutomationActionEventHub struct {
	ConnectionString   *string `json:"connectionString,omitempty"`
	EventHubResourceId *string `json:"eventHubResourceId,omitempty"`
	SasPolicyName      *string `json:"sasPolicyName,omitempty"`

	// Fields inherited from AutomationAction

	ActionType ActionType `json:"actionType"`
}

func (s AutomationActionEventHub) AutomationAction() BaseAutomationActionImpl {
	return BaseAutomationActionImpl{
		ActionType: s.ActionType,
	}
}

var _ json.Marshaler = AutomationActionEventHub{}

func (s AutomationActionEventHub) MarshalJSON() ([]byte, error) {
	type wrapper AutomationActionEventHub
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutomationActionEventHub: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutomationActionEventHub: %+v", err)
	}

	decoded["actionType"] = "EventHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutomationActionEventHub: %+v", err)
	}

	return encoded, nil
}
