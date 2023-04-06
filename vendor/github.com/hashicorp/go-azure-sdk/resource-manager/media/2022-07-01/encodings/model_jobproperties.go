package encodings

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobProperties struct {
	CorrelationData *map[string]string `json:"correlationData,omitempty"`
	Created         *string            `json:"created,omitempty"`
	Description     *string            `json:"description,omitempty"`
	EndTime         *string            `json:"endTime,omitempty"`
	Input           JobInput           `json:"input"`
	LastModified    *string            `json:"lastModified,omitempty"`
	Outputs         []JobOutput        `json:"outputs"`
	Priority        *Priority          `json:"priority,omitempty"`
	StartTime       *string            `json:"startTime,omitempty"`
	State           *JobState          `json:"state,omitempty"`
}

func (o *JobProperties) GetCreatedAsTime() (*time.Time, error) {
	if o.Created == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Created, "2006-01-02T15:04:05Z07:00")
}

func (o *JobProperties) SetCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Created = &formatted
}

func (o *JobProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *JobProperties) GetLastModifiedAsTime() (*time.Time, error) {
	if o.LastModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModified, "2006-01-02T15:04:05Z07:00")
}

func (o *JobProperties) SetLastModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModified = &formatted
}

func (o *JobProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}

var _ json.Unmarshaler = &JobProperties{}

func (s *JobProperties) UnmarshalJSON(bytes []byte) error {
	type alias JobProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into JobProperties: %+v", err)
	}

	s.CorrelationData = decoded.CorrelationData
	s.Created = decoded.Created
	s.Description = decoded.Description
	s.EndTime = decoded.EndTime
	s.LastModified = decoded.LastModified
	s.Priority = decoded.Priority
	s.StartTime = decoded.StartTime
	s.State = decoded.State

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling JobProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["input"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Input' for 'JobProperties': %+v", err)
		}
		s.Input = impl
	}

	if v, ok := temp["outputs"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Outputs into list []json.RawMessage: %+v", err)
		}

		output := make([]JobOutput, 0)
		for i, val := range listTemp {
			impl, err := unmarshalJobOutputImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Outputs' for 'JobProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Outputs = output
	}
	return nil
}
