package alertprocessingrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Recurrence = DailyRecurrence{}

type DailyRecurrence struct {

	// Fields inherited from Recurrence

	EndTime        *string        `json:"endTime,omitempty"`
	RecurrenceType RecurrenceType `json:"recurrenceType"`
	StartTime      *string        `json:"startTime,omitempty"`
}

func (s DailyRecurrence) Recurrence() BaseRecurrenceImpl {
	return BaseRecurrenceImpl{
		EndTime:        s.EndTime,
		RecurrenceType: s.RecurrenceType,
		StartTime:      s.StartTime,
	}
}

var _ json.Marshaler = DailyRecurrence{}

func (s DailyRecurrence) MarshalJSON() ([]byte, error) {
	type wrapper DailyRecurrence
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DailyRecurrence: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DailyRecurrence: %+v", err)
	}

	decoded["recurrenceType"] = "Daily"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DailyRecurrence: %+v", err)
	}

	return encoded, nil
}
