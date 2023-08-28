package automationrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationRuleCondition interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawAutomationRuleConditionImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalAutomationRuleConditionImplementation(input []byte) (AutomationRuleCondition, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AutomationRuleCondition into map[string]interface: %+v", err)
	}

	value, ok := temp["conditionType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Boolean") {
		var out BooleanConditionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BooleanConditionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PropertyArrayChanged") {
		var out PropertyArrayChangedConditionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PropertyArrayChangedConditionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PropertyArray") {
		var out PropertyArrayConditionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PropertyArrayConditionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PropertyChanged") {
		var out PropertyChangedConditionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PropertyChangedConditionProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Property") {
		var out PropertyConditionProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PropertyConditionProperties: %+v", err)
		}
		return out, nil
	}

	out := RawAutomationRuleConditionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
