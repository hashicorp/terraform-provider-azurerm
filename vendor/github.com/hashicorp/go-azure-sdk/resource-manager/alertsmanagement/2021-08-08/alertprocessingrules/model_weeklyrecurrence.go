package alertprocessingrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Recurrence = WeeklyRecurrence{}

type WeeklyRecurrence struct {
	DaysOfWeek []DaysOfWeek `json:"daysOfWeek"`

	// Fields inherited from Recurrence

	EndTime        *string        `json:"endTime,omitempty"`
	RecurrenceType RecurrenceType `json:"recurrenceType"`
	StartTime      *string        `json:"startTime,omitempty"`
}

func (s WeeklyRecurrence) Recurrence() BaseRecurrenceImpl {
	return BaseRecurrenceImpl{
		EndTime:        s.EndTime,
		RecurrenceType: s.RecurrenceType,
		StartTime:      s.StartTime,
	}
}

var _ json.Marshaler = WeeklyRecurrence{}

func (s WeeklyRecurrence) MarshalJSON() ([]byte, error) {
	type wrapper WeeklyRecurrence
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WeeklyRecurrence: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WeeklyRecurrence: %+v", err)
	}

	decoded["recurrenceType"] = "Weekly"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WeeklyRecurrence: %+v", err)
	}

	return encoded, nil
}
