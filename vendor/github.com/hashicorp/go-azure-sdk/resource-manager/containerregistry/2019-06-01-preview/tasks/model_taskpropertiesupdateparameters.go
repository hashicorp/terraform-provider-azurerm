package tasks

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaskPropertiesUpdateParameters struct {
	AgentConfiguration *AgentProperties          `json:"agentConfiguration,omitempty"`
	AgentPoolName      *string                   `json:"agentPoolName,omitempty"`
	Credentials        *Credentials              `json:"credentials,omitempty"`
	LogTemplate        *string                   `json:"logTemplate,omitempty"`
	Platform           *PlatformUpdateParameters `json:"platform,omitempty"`
	Status             *TaskStatus               `json:"status,omitempty"`
	Step               TaskStepUpdateParameters  `json:"step"`
	Timeout            *int64                    `json:"timeout,omitempty"`
	Trigger            *TriggerUpdateParameters  `json:"trigger,omitempty"`
}

var _ json.Unmarshaler = &TaskPropertiesUpdateParameters{}

func (s *TaskPropertiesUpdateParameters) UnmarshalJSON(bytes []byte) error {
	type alias TaskPropertiesUpdateParameters
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into TaskPropertiesUpdateParameters: %+v", err)
	}

	s.AgentConfiguration = decoded.AgentConfiguration
	s.AgentPoolName = decoded.AgentPoolName
	s.Credentials = decoded.Credentials
	s.LogTemplate = decoded.LogTemplate
	s.Platform = decoded.Platform
	s.Status = decoded.Status
	s.Timeout = decoded.Timeout
	s.Trigger = decoded.Trigger

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling TaskPropertiesUpdateParameters into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["step"]; ok {
		impl, err := unmarshalTaskStepUpdateParametersImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Step' for 'TaskPropertiesUpdateParameters': %+v", err)
		}
		s.Step = impl
	}
	return nil
}
