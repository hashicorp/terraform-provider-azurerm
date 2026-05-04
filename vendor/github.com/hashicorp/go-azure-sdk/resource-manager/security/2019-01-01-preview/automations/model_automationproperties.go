package automations

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationProperties struct {
	Actions     *[]AutomationAction `json:"actions,omitempty"`
	Description *string             `json:"description,omitempty"`
	IsEnabled   *bool               `json:"isEnabled,omitempty"`
	Scopes      *[]AutomationScope  `json:"scopes,omitempty"`
	Sources     *[]AutomationSource `json:"sources,omitempty"`
}

var _ json.Unmarshaler = &AutomationProperties{}

func (s *AutomationProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Description *string             `json:"description,omitempty"`
		IsEnabled   *bool               `json:"isEnabled,omitempty"`
		Scopes      *[]AutomationScope  `json:"scopes,omitempty"`
		Sources     *[]AutomationSource `json:"sources,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Description = decoded.Description
	s.IsEnabled = decoded.IsEnabled
	s.Scopes = decoded.Scopes
	s.Sources = decoded.Sources

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AutomationProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["actions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Actions into list []json.RawMessage: %+v", err)
		}

		output := make([]AutomationAction, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalAutomationActionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Actions' for 'AutomationProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Actions = &output
	}

	return nil
}
