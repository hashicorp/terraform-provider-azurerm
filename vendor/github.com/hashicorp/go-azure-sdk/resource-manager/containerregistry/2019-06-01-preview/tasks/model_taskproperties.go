package tasks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaskProperties struct {
	AgentConfiguration *AgentProperties    `json:"agentConfiguration,omitempty"`
	AgentPoolName      *string             `json:"agentPoolName,omitempty"`
	CreationDate       *string             `json:"creationDate,omitempty"`
	Credentials        *Credentials        `json:"credentials,omitempty"`
	IsSystemTask       *bool               `json:"isSystemTask,omitempty"`
	LogTemplate        *string             `json:"logTemplate,omitempty"`
	Platform           *PlatformProperties `json:"platform,omitempty"`
	ProvisioningState  *ProvisioningState  `json:"provisioningState,omitempty"`
	Status             *TaskStatus         `json:"status,omitempty"`
	Step               TaskStepProperties  `json:"step"`
	Timeout            *int64              `json:"timeout,omitempty"`
	Trigger            *TriggerProperties  `json:"trigger,omitempty"`
}

func (o *TaskProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *TaskProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}

var _ json.Unmarshaler = &TaskProperties{}

func (s *TaskProperties) UnmarshalJSON(bytes []byte) error {
	type alias TaskProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into TaskProperties: %+v", err)
	}

	s.AgentConfiguration = decoded.AgentConfiguration
	s.AgentPoolName = decoded.AgentPoolName
	s.CreationDate = decoded.CreationDate
	s.Credentials = decoded.Credentials
	s.IsSystemTask = decoded.IsSystemTask
	s.LogTemplate = decoded.LogTemplate
	s.Platform = decoded.Platform
	s.ProvisioningState = decoded.ProvisioningState
	s.Status = decoded.Status
	s.Timeout = decoded.Timeout
	s.Trigger = decoded.Trigger

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling TaskProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["step"]; ok {
		impl, err := unmarshalTaskStepPropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Step' for 'TaskProperties': %+v", err)
		}
		s.Step = impl
	}
	return nil
}
