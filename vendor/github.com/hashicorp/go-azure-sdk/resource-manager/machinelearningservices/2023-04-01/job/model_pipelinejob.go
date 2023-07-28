package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobBase = PipelineJob{}

type PipelineJob struct {
	Inputs      *map[string]JobInput    `json:"inputs,omitempty"`
	Jobs        *map[string]interface{} `json:"jobs,omitempty"`
	Outputs     *map[string]JobOutput   `json:"outputs,omitempty"`
	Settings    *interface{}            `json:"settings,omitempty"`
	SourceJobId *string                 `json:"sourceJobId,omitempty"`

	// Fields inherited from JobBase
	ComponentId    *string                `json:"componentId,omitempty"`
	ComputeId      *string                `json:"computeId,omitempty"`
	Description    *string                `json:"description,omitempty"`
	DisplayName    *string                `json:"displayName,omitempty"`
	ExperimentName *string                `json:"experimentName,omitempty"`
	Identity       IdentityConfiguration  `json:"identity"`
	IsArchived     *bool                  `json:"isArchived,omitempty"`
	Properties     *map[string]string     `json:"properties,omitempty"`
	Services       *map[string]JobService `json:"services,omitempty"`
	Status         *JobStatus             `json:"status,omitempty"`
	Tags           *map[string]string     `json:"tags,omitempty"`
}

var _ json.Marshaler = PipelineJob{}

func (s PipelineJob) MarshalJSON() ([]byte, error) {
	type wrapper PipelineJob
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PipelineJob: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PipelineJob: %+v", err)
	}
	decoded["jobType"] = "Pipeline"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PipelineJob: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &PipelineJob{}

func (s *PipelineJob) UnmarshalJSON(bytes []byte) error {
	type alias PipelineJob
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into PipelineJob: %+v", err)
	}

	s.ComponentId = decoded.ComponentId
	s.ComputeId = decoded.ComputeId
	s.Description = decoded.Description
	s.DisplayName = decoded.DisplayName
	s.ExperimentName = decoded.ExperimentName
	s.IsArchived = decoded.IsArchived
	s.Jobs = decoded.Jobs
	s.Properties = decoded.Properties
	s.Services = decoded.Services
	s.Settings = decoded.Settings
	s.SourceJobId = decoded.SourceJobId
	s.Status = decoded.Status
	s.Tags = decoded.Tags

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling PipelineJob into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["identity"]; ok {
		impl, err := unmarshalIdentityConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Identity' for 'PipelineJob': %+v", err)
		}
		s.Identity = impl
	}

	if v, ok := temp["inputs"]; ok {
		var dictionaryTemp map[string]json.RawMessage
		if err := json.Unmarshal(v, &dictionaryTemp); err != nil {
			return fmt.Errorf("unmarshaling Inputs into dictionary map[string]json.RawMessage: %+v", err)
		}

		output := make(map[string]JobInput)
		for key, val := range dictionaryTemp {
			impl, err := unmarshalJobInputImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling key %q field 'Inputs' for 'PipelineJob': %+v", key, err)
			}
			output[key] = impl
		}
		s.Inputs = &output
	}

	if v, ok := temp["outputs"]; ok {
		var dictionaryTemp map[string]json.RawMessage
		if err := json.Unmarshal(v, &dictionaryTemp); err != nil {
			return fmt.Errorf("unmarshaling Outputs into dictionary map[string]json.RawMessage: %+v", err)
		}

		output := make(map[string]JobOutput)
		for key, val := range dictionaryTemp {
			impl, err := unmarshalJobOutputImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling key %q field 'Outputs' for 'PipelineJob': %+v", key, err)
			}
			output[key] = impl
		}
		s.Outputs = &output
	}
	return nil
}
