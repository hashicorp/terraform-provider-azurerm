package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobBase = CommandJob{}

type CommandJob struct {
	CodeId               *string                   `json:"codeId,omitempty"`
	Command              string                    `json:"command"`
	Distribution         DistributionConfiguration `json:"distribution"`
	EnvironmentId        string                    `json:"environmentId"`
	EnvironmentVariables *map[string]string        `json:"environmentVariables,omitempty"`
	Inputs               *map[string]JobInput      `json:"inputs,omitempty"`
	Limits               JobLimits                 `json:"limits"`
	Outputs              *map[string]JobOutput     `json:"outputs,omitempty"`
	Parameters           *interface{}              `json:"parameters,omitempty"`
	Resources            *JobResourceConfiguration `json:"resources,omitempty"`

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

var _ json.Marshaler = CommandJob{}

func (s CommandJob) MarshalJSON() ([]byte, error) {
	type wrapper CommandJob
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CommandJob: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CommandJob: %+v", err)
	}
	decoded["jobType"] = "Command"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CommandJob: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &CommandJob{}

func (s *CommandJob) UnmarshalJSON(bytes []byte) error {
	type alias CommandJob
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into CommandJob: %+v", err)
	}

	s.CodeId = decoded.CodeId
	s.Command = decoded.Command
	s.ComponentId = decoded.ComponentId
	s.ComputeId = decoded.ComputeId
	s.Description = decoded.Description
	s.DisplayName = decoded.DisplayName
	s.EnvironmentId = decoded.EnvironmentId
	s.EnvironmentVariables = decoded.EnvironmentVariables
	s.ExperimentName = decoded.ExperimentName
	s.IsArchived = decoded.IsArchived
	s.Parameters = decoded.Parameters
	s.Properties = decoded.Properties
	s.Resources = decoded.Resources
	s.Services = decoded.Services
	s.Status = decoded.Status
	s.Tags = decoded.Tags

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CommandJob into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["distribution"]; ok {
		impl, err := unmarshalDistributionConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Distribution' for 'CommandJob': %+v", err)
		}
		s.Distribution = impl
	}

	if v, ok := temp["identity"]; ok {
		impl, err := unmarshalIdentityConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Identity' for 'CommandJob': %+v", err)
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
				return fmt.Errorf("unmarshaling key %q field 'Inputs' for 'CommandJob': %+v", key, err)
			}
			output[key] = impl
		}
		s.Inputs = &output
	}

	if v, ok := temp["limits"]; ok {
		impl, err := unmarshalJobLimitsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Limits' for 'CommandJob': %+v", err)
		}
		s.Limits = impl
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
				return fmt.Errorf("unmarshaling key %q field 'Outputs' for 'CommandJob': %+v", key, err)
			}
			output[key] = impl
		}
		s.Outputs = &output
	}
	return nil
}
