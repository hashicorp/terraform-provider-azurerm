package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Schedule struct {
	Action            ScheduleActionBase          `json:"action"`
	Description       *string                     `json:"description,omitempty"`
	DisplayName       *string                     `json:"displayName,omitempty"`
	IsEnabled         *bool                       `json:"isEnabled,omitempty"`
	Properties        *map[string]string          `json:"properties,omitempty"`
	ProvisioningState *ScheduleProvisioningStatus `json:"provisioningState,omitempty"`
	Tags              *map[string]string          `json:"tags,omitempty"`
	Trigger           TriggerBase                 `json:"trigger"`
}

var _ json.Unmarshaler = &Schedule{}

func (s *Schedule) UnmarshalJSON(bytes []byte) error {
	type alias Schedule
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into Schedule: %+v", err)
	}

	s.Description = decoded.Description
	s.DisplayName = decoded.DisplayName
	s.IsEnabled = decoded.IsEnabled
	s.Properties = decoded.Properties
	s.ProvisioningState = decoded.ProvisioningState
	s.Tags = decoded.Tags

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling Schedule into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["action"]; ok {
		impl, err := unmarshalScheduleActionBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Action' for 'Schedule': %+v", err)
		}
		s.Action = impl
	}

	if v, ok := temp["trigger"]; ok {
		impl, err := unmarshalTriggerBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Trigger' for 'Schedule': %+v", err)
		}
		s.Trigger = impl
	}
	return nil
}
