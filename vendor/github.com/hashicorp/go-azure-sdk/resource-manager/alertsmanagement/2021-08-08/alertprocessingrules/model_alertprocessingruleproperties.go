package alertprocessingrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertProcessingRuleProperties struct {
	Actions     []Action     `json:"actions"`
	Conditions  *[]Condition `json:"conditions,omitempty"`
	Description *string      `json:"description,omitempty"`
	Enabled     *bool        `json:"enabled,omitempty"`
	Schedule    *Schedule    `json:"schedule,omitempty"`
	Scopes      []string     `json:"scopes"`
}

var _ json.Unmarshaler = &AlertProcessingRuleProperties{}

func (s *AlertProcessingRuleProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Conditions  *[]Condition `json:"conditions,omitempty"`
		Description *string      `json:"description,omitempty"`
		Enabled     *bool        `json:"enabled,omitempty"`
		Schedule    *Schedule    `json:"schedule,omitempty"`
		Scopes      []string     `json:"scopes"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Conditions = decoded.Conditions
	s.Description = decoded.Description
	s.Enabled = decoded.Enabled
	s.Schedule = decoded.Schedule
	s.Scopes = decoded.Scopes

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AlertProcessingRuleProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["actions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Actions into list []json.RawMessage: %+v", err)
		}

		output := make([]Action, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalActionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Actions' for 'AlertProcessingRuleProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Actions = output
	}

	return nil
}
