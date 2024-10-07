package alertprocessingrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Recurrence interface {
	Recurrence() BaseRecurrenceImpl
}

var _ Recurrence = BaseRecurrenceImpl{}

type BaseRecurrenceImpl struct {
	EndTime        *string        `json:"endTime,omitempty"`
	RecurrenceType RecurrenceType `json:"recurrenceType"`
	StartTime      *string        `json:"startTime,omitempty"`
}

func (s BaseRecurrenceImpl) Recurrence() BaseRecurrenceImpl {
	return s
}

var _ Recurrence = RawRecurrenceImpl{}

// RawRecurrenceImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawRecurrenceImpl struct {
	recurrence BaseRecurrenceImpl
	Type       string
	Values     map[string]interface{}
}

func (s RawRecurrenceImpl) Recurrence() BaseRecurrenceImpl {
	return s.recurrence
}

func UnmarshalRecurrenceImplementation(input []byte) (Recurrence, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Recurrence into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["recurrenceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Daily") {
		var out DailyRecurrence
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DailyRecurrence: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Monthly") {
		var out MonthlyRecurrence
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MonthlyRecurrence: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Weekly") {
		var out WeeklyRecurrence
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WeeklyRecurrence: %+v", err)
		}
		return out, nil
	}

	var parent BaseRecurrenceImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseRecurrenceImpl: %+v", err)
	}

	return RawRecurrenceImpl{
		recurrence: parent,
		Type:       value,
		Values:     temp,
	}, nil

}
