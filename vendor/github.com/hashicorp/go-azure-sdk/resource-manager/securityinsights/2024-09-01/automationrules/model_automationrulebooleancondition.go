package automationrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationRuleBooleanCondition struct {
	InnerConditions *[]AutomationRuleCondition                       `json:"innerConditions,omitempty"`
	Operator        *AutomationRuleBooleanConditionSupportedOperator `json:"operator,omitempty"`
}

var _ json.Unmarshaler = &AutomationRuleBooleanCondition{}

func (s *AutomationRuleBooleanCondition) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Operator *AutomationRuleBooleanConditionSupportedOperator `json:"operator,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Operator = decoded.Operator

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AutomationRuleBooleanCondition into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["innerConditions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling InnerConditions into list []json.RawMessage: %+v", err)
		}

		output := make([]AutomationRuleCondition, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalAutomationRuleConditionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'InnerConditions' for 'AutomationRuleBooleanCondition': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.InnerConditions = &output
	}

	return nil
}
