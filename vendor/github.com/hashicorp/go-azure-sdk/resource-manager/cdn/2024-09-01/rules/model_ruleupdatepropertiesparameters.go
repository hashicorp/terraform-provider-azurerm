package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RuleUpdatePropertiesParameters struct {
	Actions                 *[]DeliveryRuleAction    `json:"actions,omitempty"`
	Conditions              *[]DeliveryRuleCondition `json:"conditions,omitempty"`
	MatchProcessingBehavior *MatchProcessingBehavior `json:"matchProcessingBehavior,omitempty"`
	Order                   *int64                   `json:"order,omitempty"`
	RuleSetName             *string                  `json:"ruleSetName,omitempty"`
}

var _ json.Unmarshaler = &RuleUpdatePropertiesParameters{}

func (s *RuleUpdatePropertiesParameters) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		MatchProcessingBehavior *MatchProcessingBehavior `json:"matchProcessingBehavior,omitempty"`
		Order                   *int64                   `json:"order,omitempty"`
		RuleSetName             *string                  `json:"ruleSetName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.MatchProcessingBehavior = decoded.MatchProcessingBehavior
	s.Order = decoded.Order
	s.RuleSetName = decoded.RuleSetName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling RuleUpdatePropertiesParameters into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["actions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Actions into list []json.RawMessage: %+v", err)
		}

		output := make([]DeliveryRuleAction, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalDeliveryRuleActionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Actions' for 'RuleUpdatePropertiesParameters': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Actions = &output
	}

	if v, ok := temp["conditions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Conditions into list []json.RawMessage: %+v", err)
		}

		output := make([]DeliveryRuleCondition, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalDeliveryRuleConditionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Conditions' for 'RuleUpdatePropertiesParameters': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Conditions = &output
	}

	return nil
}
