package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Pipeline struct {
	Activities    *[]Activity                        `json:"activities,omitempty"`
	Annotations   *[]interface{}                     `json:"annotations,omitempty"`
	Concurrency   *int64                             `json:"concurrency,omitempty"`
	Description   *string                            `json:"description,omitempty"`
	Folder        *PipelineFolder                    `json:"folder,omitempty"`
	Parameters    *map[string]ParameterSpecification `json:"parameters,omitempty"`
	Policy        *PipelinePolicy                    `json:"policy,omitempty"`
	RunDimensions *map[string]string                 `json:"runDimensions,omitempty"`
	Variables     *map[string]VariableSpecification  `json:"variables,omitempty"`
}

var _ json.Unmarshaler = &Pipeline{}

func (s *Pipeline) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Annotations   *[]interface{}                     `json:"annotations,omitempty"`
		Concurrency   *int64                             `json:"concurrency,omitempty"`
		Description   *string                            `json:"description,omitempty"`
		Folder        *PipelineFolder                    `json:"folder,omitempty"`
		Parameters    *map[string]ParameterSpecification `json:"parameters,omitempty"`
		Policy        *PipelinePolicy                    `json:"policy,omitempty"`
		RunDimensions *map[string]string                 `json:"runDimensions,omitempty"`
		Variables     *map[string]VariableSpecification  `json:"variables,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Annotations = decoded.Annotations
	s.Concurrency = decoded.Concurrency
	s.Description = decoded.Description
	s.Folder = decoded.Folder
	s.Parameters = decoded.Parameters
	s.Policy = decoded.Policy
	s.RunDimensions = decoded.RunDimensions
	s.Variables = decoded.Variables

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling Pipeline into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["activities"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Activities into list []json.RawMessage: %+v", err)
		}

		output := make([]Activity, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalActivityImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Activities' for 'Pipeline': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Activities = &output
	}

	return nil
}
