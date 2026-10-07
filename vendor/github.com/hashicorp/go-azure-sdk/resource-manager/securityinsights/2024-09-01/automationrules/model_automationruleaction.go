package automationrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationRuleAction interface {
	AutomationRuleAction() BaseAutomationRuleActionImpl
}

var _ AutomationRuleAction = BaseAutomationRuleActionImpl{}

type BaseAutomationRuleActionImpl struct {
	ActionType ActionType `json:"actionType"`
	Order      int64      `json:"order"`
}

func (s BaseAutomationRuleActionImpl) AutomationRuleAction() BaseAutomationRuleActionImpl {
	return s
}

var _ AutomationRuleAction = RawAutomationRuleActionImpl{}

// RawAutomationRuleActionImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawAutomationRuleActionImpl struct {
	automationRuleAction BaseAutomationRuleActionImpl
	Type                 string
	Values               map[string]interface{}
}

func (s RawAutomationRuleActionImpl) AutomationRuleAction() BaseAutomationRuleActionImpl {
	return s.automationRuleAction
}

func (s RawAutomationRuleActionImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalAutomationRuleActionImplementation(input []byte) (AutomationRuleAction, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AutomationRuleAction into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["actionType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AddIncidentTask") {
		var out AutomationRuleAddIncidentTaskAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutomationRuleAddIncidentTaskAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ModifyProperties") {
		var out AutomationRuleModifyPropertiesAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutomationRuleModifyPropertiesAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RunPlaybook") {
		var out AutomationRuleRunPlaybookAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutomationRuleRunPlaybookAction: %+v", err)
		}
		return out, nil
	}

	var parent BaseAutomationRuleActionImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAutomationRuleActionImpl: %+v", err)
	}

	return RawAutomationRuleActionImpl{
		automationRuleAction: parent,
		Type:                 value,
		Values:               temp,
	}, nil

}
