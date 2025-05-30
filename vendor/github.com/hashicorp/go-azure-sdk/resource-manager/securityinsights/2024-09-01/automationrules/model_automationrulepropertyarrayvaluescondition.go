package automationrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationRulePropertyArrayValuesCondition struct {
	ArrayConditionType *AutomationRulePropertyArrayConditionSupportedArrayConditionType `json:"arrayConditionType,omitempty"`
	ArrayType          *AutomationRulePropertyArrayConditionSupportedArrayType          `json:"arrayType,omitempty"`
	ItemConditions     *[]AutomationRuleCondition                                       `json:"itemConditions,omitempty"`
}

var _ json.Unmarshaler = &AutomationRulePropertyArrayValuesCondition{}

func (s *AutomationRulePropertyArrayValuesCondition) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ArrayConditionType *AutomationRulePropertyArrayConditionSupportedArrayConditionType `json:"arrayConditionType,omitempty"`
		ArrayType          *AutomationRulePropertyArrayConditionSupportedArrayType          `json:"arrayType,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ArrayConditionType = decoded.ArrayConditionType
	s.ArrayType = decoded.ArrayType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AutomationRulePropertyArrayValuesCondition into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["itemConditions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling ItemConditions into list []json.RawMessage: %+v", err)
		}

		output := make([]AutomationRuleCondition, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalAutomationRuleConditionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'ItemConditions' for 'AutomationRulePropertyArrayValuesCondition': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.ItemConditions = &output
	}

	return nil
}
