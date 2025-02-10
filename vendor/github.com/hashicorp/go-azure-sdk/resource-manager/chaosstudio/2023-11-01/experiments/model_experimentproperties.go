package experiments

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExperimentProperties struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Selectors         []Selector         `json:"selectors"`
	Steps             []Step             `json:"steps"`
}

var _ json.Unmarshaler = &ExperimentProperties{}

func (s *ExperimentProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
		Steps             []Step             `json:"steps"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ProvisioningState = decoded.ProvisioningState
	s.Steps = decoded.Steps

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ExperimentProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["selectors"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Selectors into list []json.RawMessage: %+v", err)
		}

		output := make([]Selector, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalSelectorImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Selectors' for 'ExperimentProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Selectors = output
	}

	return nil
}
