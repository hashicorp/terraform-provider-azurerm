package automationrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationRuleAction interface {
}

func unmarshalAutomationRuleActionImplementation(input []byte) (AutomationRuleAction, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AutomationRuleAction into map[string]interface: %+v", err)
	}

	value, ok := temp["actionType"].(string)
	if !ok {
		return nil, nil
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

	type RawAutomationRuleActionImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawAutomationRuleActionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
