package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ScheduleActionBase = JobScheduleAction{}

type JobScheduleAction struct {
	JobDefinition JobBase `json:"jobDefinition"`

	// Fields inherited from ScheduleActionBase
}

var _ json.Marshaler = JobScheduleAction{}

func (s JobScheduleAction) MarshalJSON() ([]byte, error) {
	type wrapper JobScheduleAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JobScheduleAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JobScheduleAction: %+v", err)
	}
	decoded["actionType"] = "CreateJob"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JobScheduleAction: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &JobScheduleAction{}

func (s *JobScheduleAction) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling JobScheduleAction into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["jobDefinition"]; ok {
		impl, err := unmarshalJobBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'JobDefinition' for 'JobScheduleAction': %+v", err)
		}
		s.JobDefinition = impl
	}
	return nil
}
