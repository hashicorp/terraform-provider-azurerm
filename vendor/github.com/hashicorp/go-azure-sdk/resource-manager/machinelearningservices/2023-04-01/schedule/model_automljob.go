package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobBase = AutoMLJob{}

type AutoMLJob struct {
	EnvironmentId        *string                   `json:"environmentId,omitempty"`
	EnvironmentVariables *map[string]string        `json:"environmentVariables,omitempty"`
	Outputs              *map[string]JobOutput     `json:"outputs,omitempty"`
	Resources            *JobResourceConfiguration `json:"resources,omitempty"`
	TaskDetails          AutoMLVertical            `json:"taskDetails"`

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

var _ json.Marshaler = AutoMLJob{}

func (s AutoMLJob) MarshalJSON() ([]byte, error) {
	type wrapper AutoMLJob
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutoMLJob: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutoMLJob: %+v", err)
	}
	decoded["jobType"] = "AutoML"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutoMLJob: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AutoMLJob{}

func (s *AutoMLJob) UnmarshalJSON(bytes []byte) error {
	type alias AutoMLJob
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into AutoMLJob: %+v", err)
	}

	s.ComponentId = decoded.ComponentId
	s.ComputeId = decoded.ComputeId
	s.Description = decoded.Description
	s.DisplayName = decoded.DisplayName
	s.EnvironmentId = decoded.EnvironmentId
	s.EnvironmentVariables = decoded.EnvironmentVariables
	s.ExperimentName = decoded.ExperimentName
	s.IsArchived = decoded.IsArchived
	s.Properties = decoded.Properties
	s.Resources = decoded.Resources
	s.Services = decoded.Services
	s.Status = decoded.Status
	s.Tags = decoded.Tags

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AutoMLJob into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["identity"]; ok {
		impl, err := unmarshalIdentityConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Identity' for 'AutoMLJob': %+v", err)
		}
		s.Identity = impl
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
				return fmt.Errorf("unmarshaling key %q field 'Outputs' for 'AutoMLJob': %+v", key, err)
			}
			output[key] = impl
		}
		s.Outputs = &output
	}

	if v, ok := temp["taskDetails"]; ok {
		impl, err := unmarshalAutoMLVerticalImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TaskDetails' for 'AutoMLJob': %+v", err)
		}
		s.TaskDetails = impl
	}
	return nil
}
