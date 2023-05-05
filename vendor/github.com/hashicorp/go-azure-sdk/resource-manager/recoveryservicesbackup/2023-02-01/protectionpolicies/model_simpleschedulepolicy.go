package protectionpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SchedulePolicy = SimpleSchedulePolicy{}

type SimpleSchedulePolicy struct {
	HourlySchedule          *HourlySchedule  `json:"hourlySchedule,omitempty"`
	ScheduleRunDays         *[]DayOfWeek     `json:"scheduleRunDays,omitempty"`
	ScheduleRunFrequency    *ScheduleRunType `json:"scheduleRunFrequency,omitempty"`
	ScheduleRunTimes        *[]string        `json:"scheduleRunTimes,omitempty"`
	ScheduleWeeklyFrequency *int64           `json:"scheduleWeeklyFrequency,omitempty"`

	// Fields inherited from SchedulePolicy
}

var _ json.Marshaler = SimpleSchedulePolicy{}

func (s SimpleSchedulePolicy) MarshalJSON() ([]byte, error) {
	type wrapper SimpleSchedulePolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SimpleSchedulePolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SimpleSchedulePolicy: %+v", err)
	}
	decoded["schedulePolicyType"] = "SimpleSchedulePolicy"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SimpleSchedulePolicy: %+v", err)
	}

	return encoded, nil
}
