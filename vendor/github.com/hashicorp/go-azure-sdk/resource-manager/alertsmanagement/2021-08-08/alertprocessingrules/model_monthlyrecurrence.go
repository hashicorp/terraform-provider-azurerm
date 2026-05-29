package alertprocessingrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Recurrence = MonthlyRecurrence{}

type MonthlyRecurrence struct {
	DaysOfMonth []int64 `json:"daysOfMonth"`

	// Fields inherited from Recurrence

	EndTime        *string        `json:"endTime,omitempty"`
	RecurrenceType RecurrenceType `json:"recurrenceType"`
	StartTime      *string        `json:"startTime,omitempty"`
}

func (s MonthlyRecurrence) Recurrence() BaseRecurrenceImpl {
	return BaseRecurrenceImpl{
		EndTime:        s.EndTime,
		RecurrenceType: s.RecurrenceType,
		StartTime:      s.StartTime,
	}
}

var _ json.Marshaler = MonthlyRecurrence{}

func (s MonthlyRecurrence) MarshalJSON() ([]byte, error) {
	type wrapper MonthlyRecurrence
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MonthlyRecurrence: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MonthlyRecurrence: %+v", err)
	}

	decoded["recurrenceType"] = "Monthly"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MonthlyRecurrence: %+v", err)
	}

	return encoded, nil
}
