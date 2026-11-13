package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScoringProfile struct {
	FunctionAggregation *ScoringFunctionAggregation `json:"functionAggregation,omitempty"`
	Functions           *[]ScoringFunction          `json:"functions,omitempty"`
	Name                string                      `json:"name"`
	Text                *TextWeights                `json:"text,omitempty"`
}

var _ json.Unmarshaler = &ScoringProfile{}

func (s *ScoringProfile) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		FunctionAggregation *ScoringFunctionAggregation `json:"functionAggregation,omitempty"`
		Name                string                      `json:"name"`
		Text                *TextWeights                `json:"text,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.FunctionAggregation = decoded.FunctionAggregation
	s.Name = decoded.Name
	s.Text = decoded.Text

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ScoringProfile into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["functions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Functions into list []json.RawMessage: %+v", err)
		}

		output := make([]ScoringFunction, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalScoringFunctionImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Functions' for 'ScoringProfile': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Functions = &output
	}

	return nil
}
