package automations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationAction interface {
	AutomationAction() BaseAutomationActionImpl
}

var _ AutomationAction = BaseAutomationActionImpl{}

type BaseAutomationActionImpl struct {
	ActionType ActionType `json:"actionType"`
}

func (s BaseAutomationActionImpl) AutomationAction() BaseAutomationActionImpl {
	return s
}

var _ AutomationAction = RawAutomationActionImpl{}

// RawAutomationActionImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawAutomationActionImpl struct {
	automationAction BaseAutomationActionImpl
	Type             string
	Values           map[string]interface{}
}

func (s RawAutomationActionImpl) AutomationAction() BaseAutomationActionImpl {
	return s.automationAction
}

func UnmarshalAutomationActionImplementation(input []byte) (AutomationAction, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AutomationAction into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["actionType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "EventHub") {
		var out AutomationActionEventHub
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutomationActionEventHub: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "LogicApp") {
		var out AutomationActionLogicApp
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutomationActionLogicApp: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Workspace") {
		var out AutomationActionWorkspace
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutomationActionWorkspace: %+v", err)
		}
		return out, nil
	}

	var parent BaseAutomationActionImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAutomationActionImpl: %+v", err)
	}

	return RawAutomationActionImpl{
		automationAction: parent,
		Type:             value,
		Values:           temp,
	}, nil

}
